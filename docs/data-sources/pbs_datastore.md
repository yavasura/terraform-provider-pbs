---
page_title: "pbs_datastore Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Reads a PBS datastore configuration.
---

# pbs_datastore

Reads an existing Proxmox Backup Server datastore configuration by name.

## Example Usage

```terraform
data "pbs_datastore" "backup" {
  name = "backup-store"
}
```

## Schema

### Required

- `name` (String) Unique datastore name.

### Read-Only

- `path` (String) Filesystem path for the datastore.
- `removable` (Boolean) Whether the datastore is backed by a removable device.
- `backing_device` (String) UUID of the filesystem partition for a removable datastore.
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
- `notify_user` (String) User to send notifications to.
- `notify_level` (String) Notification level.
- `notification_mode` (String) Notification delivery mode.
- `verify_new` (Boolean) Whether new snapshots are verified immediately after backup.
- `fingerprint` (String) Certificate fingerprint, typically used for secure remote or S3-related access.
- `digest` (String) Opaque digest returned by PBS.
- `s3_client` (String) S3 endpoint ID for an S3-backed datastore.
- `s3_bucket` (String) S3 bucket name for an S3-backed datastore.

### Read-Only Nested Blocks

#### `notify`

- `gc` (String) Garbage collection notification level.
- `prune` (String) Prune job notification level.
- `sync` (String) Sync job notification level.
- `verify` (String) Verification job notification level.

#### `maintenance_mode`

- `type` (String) Maintenance mode type.
- `message` (String) Optional maintenance message.

#### `tuning`

- `chunk_order` (String) Chunk iteration order.
- `gc_atime_cutoff` (Number) Garbage collection access time cutoff in seconds.
- `gc_atime_safety_check` (Boolean) Whether atime safety checks are enabled.
- `gc_cache_capacity` (Number) Garbage collection cache capacity.
- `sync_level` (String) Datastore fsync level.
