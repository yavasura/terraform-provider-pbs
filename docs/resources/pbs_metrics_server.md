---
page_title: "pbs_metrics_server Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a PBS metrics server configuration.
---

# pbs_metrics_server

Manages a Proxmox Backup Server metrics server configuration.

This resource supports:
- `influxdb-udp`
- `influxdb-http`

## Example Usage

### InfluxDB HTTP

```terraform
resource "pbs_metrics_server" "http_metrics" {
  name         = "influx-http"
  type         = "influxdb-http"
  url          = "https://influxdb.example.com:8086"
  organization = "ops"
  bucket       = "pbs"
  token        = var.influxdb_token
  verify_tls   = true
}
```

### InfluxDB UDP

```terraform
resource "pbs_metrics_server" "udp_metrics" {
  name     = "influx-udp"
  type     = "influxdb-udp"
  server   = "influxdb.example.com"
  port     = 8089
  enable   = true
  mtu      = 1500
  comment  = "UDP metrics export"
}
```

## Type-Specific Notes

- `influxdb-http` commonly uses `url`, `organization`, `bucket`, and `token`
- `influxdb-udp` commonly uses `server` and `port`
- `url` takes precedence over `server` and `port`

## Schema

### Required

- `name` (String) Unique metrics server name.
- `type` (String) Metrics server type. One of `influxdb-udp` or `influxdb-http`.

### Optional

- `url` (String) Full URL for InfluxDB HTTP. Takes precedence over `server` and `port`.
- `server` (String) Hostname or IP address for the metrics server.
- `port` (Number) Port for the metrics server.
- `enable` (Boolean) Whether this metrics server is enabled.
- `mtu` (Number) MTU for the metrics connection.
- `protocol` (String) Protocol for UDP metrics transport. One of `udp` or `tcp`.
- `organization` (String) InfluxDB organization for HTTP mode.
- `bucket` (String) InfluxDB bucket for HTTP mode.
- `token` (String, Sensitive) InfluxDB API token for HTTP mode.
- `max_body_size` (Number) Maximum HTTP request body size in bytes.
- `verify_tls` (Boolean) Whether to verify TLS certificates for HTTPS connections.
- `timeout` (Number) HTTP request timeout in seconds.
- `comment` (String) Description for the metrics server.

## Import

Import a metrics server by name:

```bash
terraform import pbs_metrics_server.http_metrics influx-http
```
