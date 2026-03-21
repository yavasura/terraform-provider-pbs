---
page_title: "pbs_datastores Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Lists PBS datastores.
---

# pbs_datastores

Lists Proxmox Backup Server datastores.

## Example Usage

```terraform
data "pbs_datastores" "all" {}
```

## Schema

### Read-Only

- `stores` (Attributes List) List of datastores known to PBS.

### Nested Schema for `stores`

- `name` (String) Unique datastore name.
- `path` (String) Filesystem path for the datastore.
- `removable` (Boolean) Whether the datastore is backed by a removable device.
- `backing_device` (String) UUID of the filesystem partition for a removable datastore.
- `comment` (String) Description for the datastore.
- `disabled` (Boolean) Whether the datastore is disabled.
- `gc_schedule` (String) Garbage collection schedule in cron format.
- `prune_schedule` (String) Prune schedule in cron format. Deprecated in PBS 4.0+.
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
- `fingerprint` (String) Certificate fingerprint.
- `digest` (String) Opaque digest returned by PBS.
- `s3_client` (String) S3 endpoint ID for an S3-backed datastore.
- `s3_bucket` (String) S3 bucket name for an S3-backed datastore.
