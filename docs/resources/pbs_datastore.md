---
page_title: "pbs_datastore Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a PBS datastore.
---

# pbs_datastore

Manages a Proxmox Backup Server datastore configuration.

This resource supports:
- directory-backed datastores
- removable datastores
- S3-backed datastores with a local cache path

## Example Usage

### Directory datastore

```terraform
resource "pbs_datastore" "backup" {
  name        = "backup-store"
  path        = "/mnt/datastore/backup"
  comment     = "Primary backup datastore"
  gc_schedule = "daily"
}
```

### Removable datastore

```terraform
resource "pbs_datastore" "usb_backup" {
  name           = "usb-backup"
  path           = "/mnt/removable/backup"
  removable      = true
  backing_device = "01234567-89ab-cdef-0123-456789abcdef"
}
```

### S3 datastore

```terraform
resource "pbs_s3_endpoint" "object_store" {
  id         = "backup-s3"
  endpoint   = "https://s3.example.com"
  access_key = var.s3_access_key
  secret_key = var.s3_secret_key
}

resource "pbs_datastore" "s3_backup" {
  name      = "s3-backup"
  path      = "/var/lib/pbs-s3-cache"
  s3_client = pbs_s3_endpoint.object_store.id
  s3_bucket = "pbs-backups"
}
```

## Backend Rules

- Directory datastore: set `name` and `path`
- Removable datastore: set `name`, `path`, `removable = true`, and `backing_device`
- S3 datastore: set `name`, `path`, `s3_client`, and `s3_bucket`

`s3_client` and `s3_bucket` must be provided together.

`backing_device` can only be used with `removable = true`.

## Schema

### Required

- `name` (String) Unique datastore name.

### Optional

- `path` (String) Filesystem path for the datastore. Required in practice for directory, removable, and S3 cache-backed datastores.
- `comment` (String) Description for the datastore.
- `disabled` (Boolean) Whether the datastore is disabled.
- `gc_schedule` (String) Garbage collection schedule in cron format.
- `prune_schedule` (String) Prune schedule in cron format. Deprecated in PBS 4.0+. Prefer `pbs_prune_job`.
- `keep_last` (Number) Number of most recent backups to keep.
- `keep_hourly` (Number) Number of hourly backups to keep.
- `keep_daily` (Number) Number of daily backups to keep.
- `keep_weekly` (Number) Number of weekly backups to keep.
- `keep_monthly` (Number) Number of monthly backups to keep.
- `keep_yearly` (Number) Number of yearly backups to keep.
- `removable` (Boolean) Whether the datastore is backed by a removable device.
- `backing_device` (String) Lowercase filesystem UUID for a removable datastore.
- `notify_user` (String) User to send notifications to.
- `notify_level` (String) Notification level. One of `info`, `notice`, `warning`, `error`.
- `notification_mode` (String) Notification delivery mode. One of `legacy-sendmail`, `notification-system`.
- `verify_new` (Boolean) Verify newly created snapshots immediately after backup.
- `reuse_datastore` (Boolean) Reuse existing datastore chunks when possible.
- `overwrite_in_use` (Boolean) Allow overwriting chunks that are currently in use.
- `tune_level` (Number) Tuning level from `0` to `4`. Deprecated in favor of `tuning.sync_level`.
- `fingerprint` (String) Certificate fingerprint for secure connections, typically for S3 datastores.
- `s3_client` (String) S3 endpoint ID for an S3-backed datastore.
- `s3_bucket` (String) S3 bucket name for an S3-backed datastore.

### Optional Nested Blocks

#### `notify`

- `gc` (String) Garbage collection notification level. One of `always`, `error`, `never`.
- `prune` (String) Prune job notification level. One of `always`, `error`, `never`.
- `sync` (String) Sync job notification level. One of `always`, `error`, `never`.
- `verify` (String) Verification job notification level. One of `always`, `error`, `never`.

#### `maintenance_mode`

- `type` (String) Required maintenance mode type. One of `offline`, `read-only`.
- `message` (String) Optional maintenance message.

#### `tuning`

- `chunk_order` (String) Chunk iteration order. One of `inode`, `none`.
- `gc_atime_cutoff` (Number) Garbage collection access time cutoff in seconds.
- `gc_atime_safety_check` (Boolean) Enable garbage collection access time safety checks.
- `gc_cache_capacity` (Number) Garbage collection cache capacity.
- `sync_level` (String) Datastore fsync level. One of `none`, `filesystem`, `file`.

### Read-Only

- `digest` (String) Opaque digest returned by PBS for optimistic locking.

## Import

Import a datastore by name:

```bash
terraform import pbs_datastore.backup backup-store
```
