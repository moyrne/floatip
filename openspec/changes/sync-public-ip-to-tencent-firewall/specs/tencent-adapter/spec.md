## ADDED Requirements

### Requirement: Define FirewallAdapter interface

The system SHALL define a `FirewallAdapter` interface in `internal/adapter/` that abstracts firewall rule operations.

#### Scenario: Interface methods

- **WHEN** the adapter package is defined
- **THEN** it SHALL expose a `FirewallAdapter` interface with at least:
  - `ListRules(ctx, templateID) ([]Rule, error)` — list all rules in a template
  - `CreateRule(ctx, templateID, rule) error` — create a new rule
  - `UpdateRule(ctx, templateID, ruleID, rule) error` — update an existing rule by ID

### Requirement: Tencent Cloud Lighthouse adapter

The system SHALL provide a Tencent Cloud Lighthouse implementation of `FirewallAdapter` in `internal/adapter/tencent/`.

#### Scenario: List firewall rules

- **WHEN** `ListRules` is called with a valid template ID
- **THEN** the adapter SHALL call the Tencent Cloud Lighthouse `DescribeFirewallRules` API
- **THEN** the adapter SHALL return the parsed rules list

#### Scenario: Create firewall rule

- **WHEN** `CreateRule` is called
- **THEN** the adapter SHALL call the Tencent Cloud Lighthouse `CreateFirewallRules` API
- **THEN** the adapter SHALL return the created rule

#### Scenario: Update firewall rule

- **WHEN** `UpdateRule` is called
- **THEN** the adapter SHALL call the Tencent Cloud Lighthouse `ModifyFirewallRule` API
- **THEN** the adapter SHALL return the updated rule

#### Scenario: Tencent Cloud credential error

- **WHEN** any API call fails due to invalid credentials
- **THEN** the adapter SHALL return a descriptive error
- **THEN** the adapter SHALL NOT retry (caller is responsible for retry policy)

### Requirement: Rule model

The system SHALL define a `Rule` struct used by the `FirewallAdapter` interface.

#### Scenario: Rule fields

- **WHEN** a `Rule` is defined
- **THEN** it SHALL contain: `ID`, `Protocol`, `Port`, `Source`, `Action`
