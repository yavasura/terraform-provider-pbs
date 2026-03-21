---
page_title: "pbs_metrics_server Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Reads a PBS metrics server configuration.
---

# pbs_metrics_server

Reads a specific Proxmox Backup Server metrics server configuration by name and type.

## Example Usage

```terraform
data "pbs_metrics_server" "influx" {
  name = "influx-main"
  type = "influxdb-http"
}
```

## Schema

### Required

- `name` (String) Unique metrics server name.
- `type` (String) Metrics server type. Valid values are `influxdb-http` and `influxdb-udp`.

### Read-Only

- `url` (String) Full URL for an HTTP metrics server.
- `server` (String) Server address parsed from the PBS metrics server configuration.
- `port` (Number) Server port parsed from the PBS metrics server configuration.
- `enable` (Boolean) Whether the metrics server is enabled.
- `mtu` (Number) MTU for the metrics connection, typically used for UDP.
- `organization` (String) InfluxDB organization for HTTP metrics servers.
- `bucket` (String) InfluxDB bucket for HTTP metrics servers.
- `max_body_size` (Number) Maximum HTTP request body size in bytes.
- `verify_tls` (Boolean) Whether TLS certificates are verified for HTTPS connections.
- `comment` (String) Description for the metrics server.
