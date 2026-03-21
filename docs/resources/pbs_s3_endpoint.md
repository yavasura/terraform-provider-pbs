---
page_title: "pbs_s3_endpoint Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a PBS S3 endpoint configuration.
---

# pbs_s3_endpoint

Manages a Proxmox Backup Server S3 endpoint configuration.

S3 endpoints are referenced by `pbs_datastore` when creating S3-backed datastores.

## Example Usage

```terraform
resource "pbs_s3_endpoint" "backup" {
  id         = "backup-s3"
  endpoint   = "https://s3.example.com"
  access_key = var.s3_access_key
  secret_key = var.s3_secret_key
  region     = "eu-west-1"
  port       = 443
}
```

## Schema

### Required

- `id` (String) Unique S3 endpoint identifier.
- `access_key` (String, Sensitive) Access key for the S3 object store.
- `secret_key` (String, Sensitive) Secret key for the S3 object store.
- `endpoint` (String) Endpoint used to access the S3 object store.

### Optional

- `region` (String) Region name for the S3 object store.
- `fingerprint` (String) X509 certificate fingerprint in SHA256 colon-separated format.
- `port` (Number) Port used to access the S3 object store.
- `path_style` (Boolean) Whether to use path-style bucket addressing instead of virtual-host style.
- `provider_quirks` (Set of String) Provider-specific quirks. Use `skip-if-none-match-header` for Backblaze B2 compatibility.

## Import

Import an S3 endpoint by ID:

```bash
terraform import pbs_s3_endpoint.backup backup-s3
```
