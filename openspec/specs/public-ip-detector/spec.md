## ADDED Requirements

### Requirement: Detect public IP via HTTP

The system SHALL detect the machine's current public IP by making an HTTP GET request to a configurable external IP echo service.

#### Scenario: Successful detection

- **WHEN** the system makes an HTTP GET request to the configured IP echo service
- **THEN** the system SHALL parse the response body as a plain-text IPv4 address
- **THEN** the system SHALL return the detected IP string

#### Scenario: Detection service unavailable

- **WHEN** the HTTP request to the IP echo service fails (timeout, connection refused, non-200 status)
- **THEN** the system SHALL return an error
- **THEN** the system SHALL NOT crash — the caller handles the error

### Requirement: Configurable detection interval

The system SHALL run the public IP detection on a configurable interval.

#### Scenario: Default interval

- **WHEN** no custom interval is configured
- **THEN** the system SHALL default to a 5-minute interval

#### Scenario: Custom interval

- **WHEN** a custom interval is set in configuration
- **THEN** the system SHALL use the configured value

### Requirement: Configurable echo service URL

The system SHALL allow configuration of the IP echo service URL.

#### Scenario: Custom URL

- **WHEN** a custom URL is configured
- **THEN** the system SHALL use that URL for detection
- **WHEN** no URL is configured
- **THEN** the system SHALL default to `https://checkip.amazonaws.com`
