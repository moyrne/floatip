## Why

When the public IP of a machine changes (e.g., after ISP reconnection, dynamic IP), manually updating firewall rules in Tencent Cloud Lighthouse is tedious and error-prone. This change automates the sync so firewall rules always match the current public IP, ensuring uninterrupted access while maintaining security.

## What Changes

- Add a periodic public IP detector that checks the current machine's public IP
- Add a firewall rule syncer that ensures Tencent Cloud Lighthouse firewall template rules match the detected IP
- Implement adapter pattern to decouple business logic (detect + sync) from Tencent Cloud SDK implementation
- Add configuration for sync interval and Tencent Cloud credentials/region
- Create or update firewall rules in the specified template: create if absent, update if IP mismatches

## Capabilities

### New Capabilities
- `public-ip-detector`: Detect the current public IP of the machine at a configurable interval
- `firewall-sync`: Sync public IP to Tencent Cloud Lighthouse firewall template rules using adapter pattern
- `tencent-adapter`: Adapter implementation for Tencent Cloud Lighthouse SDK (firewall rule CRUD)

### Modified Capabilities

*(None — no existing specs to modify)*

## Impact

- New dependencies: none (Tencent Cloud SDK already present)
- New Go packages under `internal/`: `detector`, `syncer`, `adapter`, `adapter/tencent`
- Configuration file or env-based config for credentials, region, template ID, and sync interval
