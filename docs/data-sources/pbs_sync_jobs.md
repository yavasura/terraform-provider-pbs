---
page_title: "pbs_sync_jobs Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Lists PBS sync jobs.
---

# pbs_sync_jobs

Lists Proxmox Backup Server sync jobs, optionally filtered by target datastore or remote.

## Example Usage

### List all sync jobs

```terraform
data "pbs_sync_jobs" "all" {}
```

### Filter by datastore

```terraform
data "pbs_sync_jobs" "backup_store" {
  store = "backup-store"
}
```

### Filter by remote

```terraform
data "pbs_sync_jobs" "remote_a" {
  remote = "remote-a"
}
```

## Schema

### Optional

- `store` (String) Filter sync jobs by target datastore name.
- `remote` (String) Filter sync jobs by remote name.

### Read-Only

- `jobs` (Attributes List) List of sync jobs returned by PBS.

### Nested Schema for `jobs`

- `id` (String) Unique sync job identifier.
- `store` (String) Target datastore name.
- `schedule` (String) Systemd calendar schedule for the job.
- `remote` (String) Remote name.
- `remote_store` (String) Remote datastore name.
- `remote_namespace` (String) Remote namespace.
- `namespace` (String) Local namespace for synced backups.
- `max_depth` (Number) Maximum namespace traversal depth.
- `group_filter` (List of String) Backup group filters applied to the sync job.
- `remove_vanished` (Boolean) Whether vanished backups are removed locally.
- `comment` (String) Description for the sync job.
- `digest` (String) Opaque digest returned by PBS.
