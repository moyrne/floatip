## 1. Adapter Interface & Pagination

- [x] 1.1 Add `GetRuleByDescription(ctx, templateID, description) (*Rule, error)` and `GetInstanceRuleByDescription(ctx, instanceID, description) (*Rule, error)` to `FirewallAdapter` interface in `internal/adapter/adapter.go`
- [x] 1.2 Add pagination loop to `TencentAdapter.ListRules` using `Offset`/`Limit` parameters
- [x] 1.3 Add pagination loop to `TencentAdapter.ListInstanceRules` using `Offset`/`Limit` parameters
- [x] 1.4 Implement `GetRuleByDescription` on `TencentAdapter` — reuse paginated `ListRules` internally, filter by description
- [x] 1.5 Implement `GetInstanceRuleByDescription` on `TencentAdapter` — reuse paginated `ListInstanceRules` internally, filter by description

## 2. Syncer Refactor

- [x] 2.1 Replace manual `slices.IndexFunc` filtering in `sync` with `s.adapter.GetRuleByDescription`
- [x] 2.2 Extract a `syncRule` helper that returns `*adapter.Rule` to flatten the three-way branch
- [x] 2.3 Remove redundant `s.applyTemplate()` calls from inside the create/update branches
- [x] 2.4 Delegate all `applyTemplate` logic to `ensureInstanceSync` — call it once after `syncRule`

## 3. Test Updates

- [x] 3.1 Add `GetRuleByDescription` and `GetInstanceRuleByDescription` to the mock adapter
- [x] 3.2 Update sync tests for refactored control flow with proper instance rule isolation
- [x] 3.3 Update `ensureInstanceSync` tests to pass `*adapter.Rule` instead of `[]adapter.Rule`
- [x] 3.4 Remove obsolete `rulesEqual` tests
- [x] 3.5 Run `make test` and verify all tests pass
