---
page_title: "pbs_prune_jobs Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Lists PBS prune jobs.
---

# pbs_prune_jobs

Lists Proxmox Backup Server prune jobs, optionally filtered by datastore.

## Example Usage

### List all prune jobs

```terraform
data "pbs_prune_jobs" "all" {}
```

### Filter by datastore

```terraform
data "pbs_prune_jobs" "backup_store" {
  store = "backup-store"
}
```

## Schema

### Optional

- `store` (String) Filter prune jobs by datastore name.

### Read-Only

- `jobs` (Attributes List) List of prune jobs returned by PBS.

### Nested Schema for `jobs`

- `id` (String) Unique prune job identifier.
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
