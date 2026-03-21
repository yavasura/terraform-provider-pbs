---
page_title: "pbs_metrics_servers Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Lists PBS metrics server configurations.
---

# pbs_metrics_servers

Lists Proxmox Backup Server metrics server configurations.

## Example Usage

```terraform
data "pbs_metrics_servers" "all" {}
```

## Schema

### Read-Only

- `servers` (Attributes List) List of metrics servers returned by PBS.

### Nested Schema for `servers`

- `name` (String) Unique metrics server name.
- `type` (String) Metrics server type.
- `url` (String) Full URL for an HTTP metrics server.
- `server` (String) Server address parsed from the PBS metrics server configuration.
- `port` (Number) Server port parsed from the PBS metrics server configuration.
- `enable` (Boolean) Whether the metrics server is enabled.
- `mtu` (Number) MTU for the metrics connection.
- `organization` (String) InfluxDB organization for HTTP metrics servers.
- `bucket` (String) InfluxDB bucket for HTTP metrics servers.
- `max_body_size` (Number) Maximum HTTP request body size in bytes.
- `verify_tls` (Boolean) Whether TLS certificates are verified for HTTPS connections.
- `comment` (String) Description for the metrics server.
