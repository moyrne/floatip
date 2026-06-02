## 1. Project Setup

- [x] 1.1 Create package directories: `internal/adapter/`, `internal/adapter/tencent/`, `internal/detector/`, `internal/syncer/`
- [x] 1.2 Add configuration struct with fields (sync interval, IP echo URL, Tencent Cloud credentials/region, template ID, rule defaults)

## 2. Adapter Interface & Model

- [x] 2.1 Define `Rule` struct in `internal/adapter/` with fields: ID, Protocol, Port, Source, Action
- [x] 2.2 Define `FirewallAdapter` interface in `internal/adapter/` with methods: ListRules, CreateRule, UpdateRule

## 3. Tencent Cloud Adapter

- [x] 3.1 Implement `NewTencentAdapter` constructor in `internal/adapter/tencent/` (takes credentials, region)
- [x] 3.2 Implement `ListRules` — calls DescribeFirewallTemplateRules API and maps response to []Rule
- [x] 3.3 Implement `CreateRule` — calls CreateFirewallTemplateRules API
- [x] 3.4 Implement `UpdateRule` — calls ReplaceFirewallTemplateRule API

## 4. Public IP Detector

- [x] 4.1 Implement `DetectPublicIP` function in `internal/detector/` — HTTP GET to configurable URL, parse plain-text IPv4
- [x] 4.2 Handle error cases (timeout, non-200, invalid response)

## 5. Firewall Sync Logic

- [x] 5.1 Implement `Syncer` struct in `internal/syncer/` that holds a `FirewallAdapter` and configuration
- [x] 5.2 Implement `Sync` method: detect IP → list rules → compare → create or update as needed
- [x] 5.3 Implement periodic loop: `Run(ctx)` that ticks on interval, calls detect then sync

## 6. Integration & Entry Point

- [x] 6.1 Wire configuration loading (env vars / config file via viper)
- [x] 6.2 Wire main func: create adapter, create syncer, start periodic loop with graceful shutdown
- [x] 6.3 Add build target to Makefile

## 7. Testing

- [x] 7.1 Write unit tests for `syncer.Sync` with mock `FirewallAdapter`
- [x] 7.2 Write unit tests for detector with a test HTTP server
- [x] 7.3 Write unit tests for tencent adapter (integration or mocked SDK)
