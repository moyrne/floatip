## Why

The current `ListRules` method lacks pagination handling, so it may miss rules when the API returns paginated results. The syncer's `sync` method duplicates logic with `ensureInstanceSync`, has deeply nested if-else branches, and calls `applyTemplate` redundantly.

## What Changes

1. **Add `GetRuleByDescription` to `FirewallAdapter`** — a new method that retrieves a single rule by its description field, with proper pagination support internally. Applies to both template rules and instance rules.
2. **Optimize `Syncer.sync`** — remove duplicated logic between `sync` and `ensureInstanceSync`, flatten the if-else chain, and eliminate redundant `applyTemplate` calls.

## Capabilities

### New Capabilities
*(none — all changes modify existing capabilities)*

### Modified Capabilities
- `tencent-adapter`: Add `GetRuleByDescription(ctx, templateID, description) (*Rule, error)` and `GetInstanceRuleByDescription(ctx, instanceID, description) (*Rule, error)` to the `FirewallAdapter` interface with pagination-aware implementations.
- `firewall-sync`: Optimize `sync` and `ensureInstanceSync` to use the new adapter methods, remove duplicated apply logic, and flatten control flow.

## Impact

- `internal/adapter/adapter.go` — `FirewallAdapter` interface gains two new methods
- `internal/adapter/tencent/adapter.go` — implement the new methods with pagination
- `internal/syncer/syncer.go` — refactor `sync` and `ensureInstanceSync`
- Tests in `internal/adapter/tencent/adapter_test.go` and `internal/syncer/syncer_test.go` — update to cover new methods and refactored logic
