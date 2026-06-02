# floatip

[**中文**](README.md) | [**English**](README.en.md)

Automatically sync the current host's public IP to a cloud firewall template and apply it to target instances.

## How It Works

1. **Detect IP** — Periodically query an IP echo service to obtain the current public IP.
2. **Sync Rule** — Find or update a rule in the firewall template with the current public IP as the source.
3. **Apply Template** — Apply the template to target instances to keep their firewall rules in sync.

## Quick Start

```bash
make build    # Build
make run      # Run (requires config.yaml)
make test     # Test
```

## Docker

```bash
docker build -t floatip .
docker run -v $(pwd)/config.yaml:/etc/floatip/config.yaml floatip
```

## Deployment

Kubernetes manifests are available in the `deploy/` directory.

## Supported Providers

| Provider | Status |
|----------|--------|
| Tencent Cloud Lighthouse | ✅ Supported |
