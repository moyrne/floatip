## ADDED Requirements

### Requirement: Sync detected IP to firewall template

The system SHALL synchronize the detected public IP to the configured Tencent Cloud Lighthouse firewall template.

#### Scenario: IP not present in any rule

- **WHEN** the detected IP does not exist in any firewall rule's source field
- **THEN** the system SHALL create a new rule that allows traffic from the detected IP

#### Scenario: IP exists but source differs

- **WHEN** a rule exists in the template with a `Source` field that differs from the detected IP
- **THEN** the system SHALL update that rule's `Source` to the current public IP

#### Scenario: IP already matches

- **WHEN** the detected IP already matches an existing rule's `Source` field
- **THEN** the system SHALL take no action (skip)

### Requirement: Configurable rule properties

The system SHALL allow configuration of the rule's protocol, port, and action.

#### Scenario: Default rule properties

- **WHEN** no custom rule properties are configured
- **THEN** the system SHALL use defaults: protocol `tcp`, port `22`, action `accept`

#### Scenario: Custom rule properties

- **WHEN** custom protocol, port, or action are configured
- **THEN** the system SHALL use the configured values when creating/updating rules

### Requirement: Periodic sync loop

The system SHALL run the detect-then-sync loop on a configurable interval.

#### Scenario: Sync cycle

- **WHEN** the sync interval elapses
- **THEN** the system SHALL detect the public IP
- **THEN** the system SHALL sync it to the firewall template
- **THEN** the system SHALL wait for the next interval

#### Scenario: Detection failure during sync

- **WHEN** IP detection fails during a sync cycle
- **THEN** the system SHALL log the error
- **THEN** the system SHALL skip the sync for this cycle
- **THEN** the system SHALL wait for the next interval
