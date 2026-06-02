## MODIFIED Requirements

### Requirement: Sync detected IP to firewall template

The system SHALL synchronize the detected public IP to the configured Tencent Cloud Lighthouse firewall template.

#### Scenario: IP not present in any rule

- **WHEN** no rule with the configured `rule_description` exists in the template
- **THEN** the system SHALL create a new rule using the configured protocol, port, action, and the detected IP
- **THEN** the system SHALL trigger instance sync after creation

#### Scenario: IP exists but CIDR differs

- **WHEN** a rule with the configured `rule_description` has a `CidrBlock` different from the detected IP
- **THEN** the system SHALL update that rule's `CidrBlock` to the current public IP
- **THEN** the system SHALL trigger instance sync after update

#### Scenario: IP already matches

- **WHEN** a rule with the configured `rule_description` already has a matching `CidrBlock`
- **THEN** the system SHALL still verify instance firewall consistency
- **THEN** the system SHALL skip template rule mutations

### Requirement: Instance firewall consistency check

The system SHALL verify that all managed instances have firewall rules matching the template after each sync.

#### Scenario: All instances in sync

- **WHEN** all configured instances have firewall rules matching the template
- **THEN** the system SHALL NOT call `ApplyTemplate`
- **THEN** the system SHALL log that instances are in sync

#### Scenario: One or more instances out of sync

- **WHEN** any instance has firewall rules that differ from the template
- **THEN** the system SHALL call `ApplyTemplate` once
- **THEN** the system SHALL log the drift detection

## ADDED Requirements

### Requirement: Lookup managed rule via adapter

The sync logic SHALL use `adapter.GetRuleByDescription` to find the managed rule instead of filtering `ListRules` results manually.

#### Scenario: GetRuleByDescription used

- **WHEN** the sync cycle starts
- **THEN** the system SHALL call `adapter.GetRuleByDescription` with the configured `rule_description`
- **THEN** the system SHALL use the returned rule (or `nil`) to decide create vs. update
