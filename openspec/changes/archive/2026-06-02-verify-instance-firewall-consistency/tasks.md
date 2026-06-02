## 1. Adapter Interface Changes

- [x] 1.1 Add `ListInstanceRules(ctx, instanceID string) ([]Rule, error)` to the `FirewallAdapter` interface in `internal/adapter/adapter.go`

## 2. Tencent Cloud Adapter

- [x] 2.1 Implement `ListInstanceRules` in `internal/adapter/tencent/adapter.go` — calls Tencent Cloud Lighthouse `DescribeFirewallRules` API and maps response to `[]adapter.Rule`
- [x] 2.2 Handle nil-pointer safety for optional response fields (CidrBlock, Port, FirewallRuleDescription)

## 3. Syncer Verification Logic

- [x] 3.1 Implement `ensureInstanceSync(ctx context.Context, templateRules []adapter.Rule)` method on `Syncer` that iterates `cfg.InstanceIDs`, calls `ListInstanceRules` for each, compares against template rules, and calls `ApplyTemplate` if drift is detected
- [x] 3.2 Implement rule set comparison function: sort both slices by canonical key, compare element-wise by Protocol/Port/CidrBlock/Action/Description

## 4. Integration

- [x] 4.1 Call `ensureInstanceSync` at the end of `sync()` after the create/update/apply logic
- [x] 4.2 Pass the already-fetched template rules to `ensureInstanceSync` to avoid redundant API calls

## 5. Testing

- [x] 5.1 Write unit tests for `ensureInstanceSync` with mock `FirewallAdapter` covering: all instances match, one instance diverges, multiple instances one diverges, no instances configured, `ListInstanceRules` failure
- [x] 5.2 Write unit tests for the Tencent adapter's `ListInstanceRules` (mocked SDK response)
- [x] 5.3 Verify existing tests still pass (`make test`)
