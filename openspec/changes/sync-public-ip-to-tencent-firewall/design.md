## Context

The project `floatip` already depends on `tencentcloud-sdk-go` for Lighthouse. Currently there is no automated sync between the machine's public IP and Tencent Cloud firewall rules. The user must manually update firewall templates when their public IP changes.

## Goals / Non-Goals

**Goals:**
- Periodically detect the machine's public IP
- Sync the detected IP to a configurable Tencent Cloud Lighthouse firewall template
- Use adapter pattern to separate business logic from Tencent Cloud SDK (enabling testability and future cloud provider support)
- Configurable sync interval, credentials, region, and template ID

**Non-Goals:**
- Not a HTTP daemon or server — this is a background periodic task
- No UI or CLI for manual sync (may be added later)
- No support for multiple IPs or complex firewall rule management beyond the public IP sync use case

## Decisions

### Adapter Pattern for Cloud Provider Abstraction

Introduce an interface `FirewallAdapter` in `internal/adapter/` with methods for CRUD on firewall rules. The Tencent Cloud implementation lives in `internal/adapter/tencent/`. This keeps `syncer` package testable with mocks and allows adding other cloud providers without changing business logic.

**Alternatives considered:**
- Direct SDK calls in business logic: simpler but untestable and tightly coupled
- Generic cloud abstraction library (e.g., cross-provider SDK): overkill for a single use case

### Public IP Detection via External Service

Use a simple HTTP client to query an external IP echo service (e.g., `https://checkip.amazonaws.com` or `https://api.ipify.org`). Configurable URL via configuration.

**Alternatives considered:**
- Network interface inspection: unreliable when behind NAT
- STUN protocol: too complex for this use case

### No Persistent State

Reconcile from scratch on each sync cycle — read all rules from the template, compare with current IP, create/update as needed. Avoids state drift and simplifies error recovery.

## Risks / Trade-offs

- [External IP service down] → Sync fails gracefully, logs warning, retries next cycle
- [Tencent Cloud API rate limits] → Add retry with backoff in adapter
- [Credentials expired] → Log error clearly on startup
- [Adapter abstraction overhead] → Acceptable for testability; only 3-4 methods in the interface
