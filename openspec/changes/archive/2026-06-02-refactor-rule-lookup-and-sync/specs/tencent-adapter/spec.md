## ADDED Requirements

### Requirement: FirewallAdapter SHALL expose GetRuleByDescription

The `FirewallAdapter` interface SHALL include `GetRuleByDescription(ctx, templateID, description) (*Rule, error)` to retrieve a single template rule by its description field.

#### Scenario: Matching rule found

- **WHEN** `GetRuleByDescription` is called with a description that matches an existing rule in the template
- **THEN** the adapter SHALL return the matching `Rule`
- **THEN** the adapter SHALL NOT return an error

#### Scenario: No matching rule found

- **WHEN** `GetRuleByDescription` is called with a description that does not match any rule
- **THEN** the adapter SHALL return `nil` for the rule
- **THEN** the adapter SHALL NOT return an error

#### Scenario: Pagination handled transparently

- **WHEN** the template contains more rules than a single API page
- **THEN** the adapter SHALL paginate through all pages internally
- **THEN** the adapter SHALL return the matching rule from any page

### Requirement: FirewallAdapter SHALL expose GetInstanceRuleByDescription

The `FirewallAdapter` interface SHALL include `GetInstanceRuleByDescription(ctx, instanceID, description) (*Rule, error)` to retrieve a single instance firewall rule by its description field.

#### Scenario: Matching rule found

- **WHEN** `GetInstanceRuleByDescription` is called with a description that matches an existing rule on the instance
- **THEN** the adapter SHALL return the matching `Rule`

#### Scenario: No matching rule found

- **WHEN** `GetInstanceRuleByDescription` is called with a description that does not match any rule on the instance
- **THEN** the adapter SHALL return `nil` for the rule
- **THEN** the adapter SHALL NOT return an error

### Requirement: ListRules SHALL paginate

`ListRules` SHALL handle API pagination to ensure all rules are returned.

#### Scenario: Paginated API response

- **WHEN** the Tencent API returns a paginated response for `DescribeFirewallTemplateRules`
- **THEN** the adapter SHALL iterate through all pages using `Offset` and `Limit` parameters
- **THEN** the adapter SHALL return the complete set of rules

### Requirement: ListInstanceRules SHALL paginate

`ListInstanceRules` SHALL handle API pagination to ensure all rules are returned.

#### Scenario: Paginated API response

- **WHEN** the Tencent API returns a paginated response for `DescribeFirewallRules`
- **THEN** the adapter SHALL iterate through all pages using `Offset` and `Limit` parameters
- **THEN** the adapter SHALL return the complete set of rules
