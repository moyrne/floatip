## ADDED Requirements

### Requirement: Adapter supports listing instance firewall rules

The `FirewallAdapter` interface SHALL provide a `ListInstanceRules(ctx, instanceID string) ([]Rule, error)` method to query firewall rules currently applied to a Lighthouse instance.

#### Scenario: List instance rules returns applied rules

- **WHEN** `ListInstanceRules` is called with a valid instance ID
- **THEN** the adapter SHALL call the Tencent Cloud Lighthouse `DescribeFirewallRules` API
- **THEN** the adapter SHALL return the list of firewall rules currently applied to that instance

#### Scenario: Instance with no firewall rules

- **WHEN** `ListInstanceRules` is called for an instance that has no firewall rules
- **THEN** the adapter SHALL return an empty slice (not nil)

#### Scenario: ListInstanceRules error propagation

- **WHEN** the underlying API call fails (invalid instance ID, network error, auth failure)
- **THEN** the adapter SHALL return a descriptive error wrapping the root cause

### Requirement: Syncer verifies instance firewall consistency

After `sync()` completes, the syncer SHALL verify that all configured instances have firewall rules matching the template, and re-apply the template if drift is detected.

#### Scenario: All instances match template

- **WHEN** `ensureInstanceSync` is called after sync
- **WHEN** all instances' firewall rules match the template rules
- **THEN** the syncer SHALL NOT call `ApplyTemplate`
- **THEN** the syncer SHALL log that instances are in sync

#### Scenario: Instance firewall diverges from template

- **WHEN** `ensureInstanceSync` is called after sync
- **WHEN** at least one instance's firewall rules differ from the template rules
- **THEN** the syncer SHALL call `ApplyTemplate` with the template ID and all configured instance IDs
- **THEN** the syncer SHALL log that drift was detected and template was re-applied

#### Scenario: No instances configured

- **WHEN** `ensureInstanceSync` is called
- **WHEN** `cfg.InstanceIDs` is empty
- **THEN** the syncer SHALL skip the verification and return immediately

#### Scenario: ListInstanceRules fails for one instance

- **WHEN** `ensureInstanceSync` is called
- **WHEN** `ListInstanceRules` fails for one or more instances
- **THEN** the syncer SHALL log the error for each failed instance
- **THEN** the syncer SHALL continue checking remaining instances
- **THEN** the syncer SHALL call `ApplyTemplate` if any other instance diverged
