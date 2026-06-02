## Why

After syncing the public IP to the firewall template, `ApplyTemplate` is only called when a rule actually changes. Drift can occur — the API call may fail silently, instances may be manually modified in the Tencent Cloud console, or template application may not propagate to all instances. This change adds a verification step to detect and correct such drift automatically.

## What Changes

- Add a new `ListInstanceRules` method to the `FirewallAdapter` interface to query firewall rules currently applied to a Lighthouse instance
- Implement `ListInstanceRules` in the Tencent adapter using the `DescribeFirewallRules` API
- Add a new syncer function `ensureInstanceSync` that compares each instance's actual firewall rules against the template rules and calls `ApplyTemplate` if they diverge
- Call the new function after each `sync()` cycle

## Capabilities

### New Capabilities
- `instance-firewall-audit`: After template sync, verify that all target instances have firewall rules matching the template, and re-apply the template if drift is detected

### Modified Capabilities

*(None — no existing specs to modify)*

## Impact

- `internal/adapter/adapter.go`: New `ListInstanceRules(ctx, instanceID string) ([]Rule, error)` on the interface
- `internal/adapter/tencent/adapter.go`: Implementation using Tencent Cloud `DescribeFirewallRules` API
- `internal/syncer/syncer.go`: New `ensureInstanceSync` function; called at end of `sync()`
- No new dependencies
