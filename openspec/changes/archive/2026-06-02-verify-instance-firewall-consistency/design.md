## Context

The current `sync()` flow: detect IP → list template rules → create/update rule if needed → `ApplyTemplate`. `ApplyTemplate` is only called when a rule actually changes. If the API fails silently on `ApplyTemplate`, or if an instance's firewall is modified outside floatip (manual console edits, another tool), drift goes undetected until the next IP change forces a re-apply.

This design adds a post-sync audit that checks each instance's actual firewall rules against the template and re-applies if they diverge.

## Goals / Non-Goals

**Goals:**
- Detect when a Lighthouse instance's applied firewall rules differ from the template rules
- Re-apply the template automatically when drift is detected
- Log clear messages for both match and drift cases

**Non-Goals:**
- Not repairing individual rules on instances — the fix is always a full template re-apply
- Not changing the existing sync logic or its decision tree
- Not adding real-time monitoring or alerting

## Decisions

### New Adapter Method: `ListInstanceRules`

Add `ListInstanceRules(ctx, instanceID string) ([]Rule, error)` to the `FirewallAdapter` interface. The Tencent adapter implements it via `DescribeFirewallRules` API. This method reads the rules *currently applied to an instance*, as opposed to `ListRules` which reads the template definition.

**Alternatives considered:**
- Calling `ApplyTemplate` unconditionally every cycle: wasteful, slower, and risks rate limits
- Comparing via `DescribeFirewallTemplateTemplateInstanceVersions` API: not available in the Lighthouse SDK; the instance query is the only reliable source of truth

### Comparison Strategy

After `sync()` completes, iterate over `cfg.InstanceIDs`. For each instance:
1. Call `ListInstanceRules` to get current rules
2. Call `ListRules` to get template rules (already available from sync, or re-fetch)
3. If the rule sets differ (by Protocol / Port / CidrBlock / Action / Description), call `ApplyTemplate`

Since rule ordering is not guaranteed, comparison uses set equality — sort both slices by a canonical key then compare element-wise.

### Idempotent Re-apply

`ApplyTemplate` is idempotent (Tencent Cloud applies the template regardless of current state). Calling it on a drifted instance only brings it back in sync without side effects on already-synced instances.

## Risks / Trade-offs

- [Extra API call per cycle] → `ListRules` is already called in `sync()`. Re-use its result to avoid a duplicate call. Each `ListInstanceRules` call is one API request per instance. For typical deployments (1-5 instances), this is negligible.
- [Rate limiting on DescribeFirewallRules] → This is a read API with high quota; unlikely to hit limits at 5-min intervals.
- [Consistency window] → Between the ListInstanceRules check and the ApplyTemplate call, the instance could change again. This is acceptable — next cycle will catch it.
