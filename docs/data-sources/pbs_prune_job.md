---
page_title: "pbs_prune_job Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Reads a PBS prune job.
---

# pbs_prune_job

Reads an existing Proxmox Backup Server prune job by ID.

## Example Usage

```terraform
data "pbs_prune_job" "daily" {
  id = "daily-prune"
}
```

## Schema

### Required

- `id` (String) Unique prune job identifier.

### Read-Only

- `store` (String) Datastore name where pruning is performed.
- `schedule` (String) Systemd calendar schedule for the job.
- `keep_last` (Number) Number of most recent backups to keep.
- `keep_hourly` (Number) Number of hourly backups to keep.
- `keep_daily` (Number) Number of daily backups to keep.
- `keep_weekly` (Number) Number of weekly backups to keep.
- `keep_monthly` (Number) Number of monthly backups to keep.
- `keep_yearly` (Number) Number of yearly backups to keep.
- `max_depth` (Number) Maximum namespace traversal depth.
- `namespace` (String) Namespace filter.
- `comment` (String) Description for the prune job.
- `disable` (Boolean) Whether the prune job is disabled.
- `digest` (String) Opaque digest returned by PBS.
