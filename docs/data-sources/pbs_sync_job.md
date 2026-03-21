---
page_title: "pbs_sync_job Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Reads a PBS sync job.
---

# pbs_sync_job

Reads an existing Proxmox Backup Server sync job by ID.

## Example Usage

```terraform
data "pbs_sync_job" "mirror" {
  id = "daily-sync"
}
```

## Schema

### Required

- `id` (String) Unique sync job identifier.

### Read-Only

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
