## Context

The `FirewallAdapter` interface exposes `ListRules` and `ListInstanceRules` which return all rules. Callers must filter by description client-side — this logic lives in `Syncer.sync`. The Tencent API may paginate results, but the current implementation does not handle pagination, meaning `ListRules` can silently return partial results.

The `Syncer.sync` method contains a three-way branch (no change / update / create) with deep if-else nesting, and after each mutation it calls `applyTemplate` directly, then `ensureInstanceSync` also checks drift and may call `applyTemplate` again — duplicating both logic and API calls.

## Goals / Non-Goals

**Goals:**
- Add `GetRuleByDescription` and `GetInstanceRuleByDescription` to the adapter interface with pagination-aware implementations
- Flatten the `sync` method's control flow and eliminate duplication with `ensureInstanceSync`
- Ensure all Tencent API pagination scenarios are handled

**Non-Goals:**
- No changes to the public IP detection logic or config model
- No changes to the `applyTemplate` mechanism itself
- No changes to the periodic scheduling in `Run`

## Decisions

1. **Interface methods vs. adapter-only helpers** — Adding `GetRuleByDescription` as interface methods rather than package-level helpers ensures every adapter implementation provides an efficient lookup, and allows future implementations (e.g., mock, AWS) to use native filtering APIs.

2. **Pagination strategy** — Use an `Offset`-based or `NextToken`-based loop depending on the Tencent API. The Tencent Lighthouse `DescribeFirewallTemplateRules` API supports `Offset` and `Limit` parameters. The implementation will loop until `len(results) < Limit`.

3. **Unify `sync` and `ensureInstanceSync`** — Instead of `sync` calling `applyTemplate` on mutation and then calling `ensureInstanceSync` which may call `applyTemplate` again, `sync` will only detect IP, find/create/update the rule, then delegate to `ensureInstanceSync` which handles all `applyTemplate` calls. This removes the redundant `applyTemplate` calls in `sync`.

4. **Flatten the three-way branch** — Extract the create/update logic into a helper that returns the resulting `[]Rule` and a `dirty` flag. The top-level `sync` then only needs one `if dirty { ensureInstanceSync(...) }`.

## Risks / Trade-offs

- **[API rate limits]** — Adding pagination loops increases API calls. → Mitigation: The expected rule count is low (<100), so pagination adds at most 1-2 extra calls per sync.
- **[Backward compatibility]** — Adding methods to the interface breaks external implementations. → Mitigation: No external implementations exist; only `TencentAdapter` implements it.
