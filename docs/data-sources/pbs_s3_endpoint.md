---
page_title: "pbs_s3_endpoint Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Reads a PBS S3 endpoint configuration.
---

# pbs_s3_endpoint

Reads a specific Proxmox Backup Server S3 endpoint configuration by ID.

## Example Usage

```terraform
data "pbs_s3_endpoint" "object_store" {
  id = "backup-s3"
}
```

## Schema

### Required

- `id` (String) Unique S3 endpoint identifier.

### Read-Only

- `access_key` (String, Sensitive) Access key for the S3 object store.
- `endpoint` (String) S3 service endpoint URL or host.
- `region` (String) S3 region.
- `fingerprint` (String) X.509 certificate fingerprint.
- `port` (Number) Port used to access the S3 object store.
- `path_style` (Boolean) Whether path-style bucket addressing is used.
- `provider_quirks` (Set of String) Provider-specific compatibility flags.
