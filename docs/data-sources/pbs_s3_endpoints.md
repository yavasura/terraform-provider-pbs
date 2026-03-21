---
page_title: "pbs_s3_endpoints Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Lists PBS S3 endpoint configurations.
---

# pbs_s3_endpoints

Lists Proxmox Backup Server S3 endpoint configurations.

## Example Usage

```terraform
data "pbs_s3_endpoints" "all" {}
```

## Schema

### Read-Only

- `endpoints` (Attributes List) List of S3 endpoints returned by PBS.

### Nested Schema for `endpoints`

- `id` (String) Unique S3 endpoint identifier.
- `access_key` (String, Sensitive) Access key for the S3 object store.
- `endpoint` (String) S3 service endpoint URL or host.
- `region` (String) S3 region.
- `fingerprint` (String) X.509 certificate fingerprint.
- `port` (Number) Port used to access the S3 object store.
- `path_style` (Boolean) Whether path-style bucket addressing is used.
- `provider_quirks` (Set of String) Provider-specific compatibility flags.
