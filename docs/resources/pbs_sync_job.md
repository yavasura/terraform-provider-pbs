---
page_title: "pbs_sync_job Resource - Proxmox Backup Server"
subcategory: ""
description: |-
Manages a PBS sync job.

Use [`pbs_namespace`](./pbs_namespace.md) to pre-create the local target namespace when syncing into a hierarchical namespace path.
---

# pbs_sync_job

Manages a Proxmox Backup Server sync job.

Sync jobs replicate backup data between PBS instances by pulling from or pushing to a configured remote.

## Example Usage

```terraform
resource "pbs_sync_job" "replicate" {
  id           = "nightly-sync"
  store        = pbs_datastore.backup.name
  remote       = pbs_remote.secondary_site.name
  remote_store = "remote-backups"
  schedule     = "daily"

  remove_vanished = true
  verified_only   = true
  comment         = "Nightly replication from secondary PBS"
}
```

## Schema

### Required

- `id` (String) Unique sync job identifier.
- `store` (String) Local datastore name where backups are written.
- `schedule` (String) When to run the job in systemd calendar format.
- `remote` (String) Remote name configured with `pbs_remote`.
- `remote_store` (String) Datastore name on the remote PBS server.

### Optional

- `remote_namespace` (String) Remote namespace to sync from.
- `namespace` (String) Local namespace where synced backups will be stored.
- `max_depth` (Number) Maximum namespace traversal depth.
- `group_filter` (List of String) Backup group selectors. Supported forms are `group:<name>`, `type:vm`, `type:ct`, `type:host`, or `regex:<pattern>`.
- `remove_vanished` (Boolean) Remove local backups that no longer exist on the remote.
- `resync_corrupt` (Boolean) Resync snapshots whose data is corrupt.
- `encrypted_only` (Boolean) Sync only encrypted backups.
- `verified_only` (Boolean) Sync only successfully verified backups.
- `run_on_mount` (Boolean) Run immediately after the datastore is mounted.
- `transfer_last` (Number) Only transfer backups newer than the last N seconds.
- `sync_direction` (String) Sync direction. One of `pull` or `push`.
- `owner` (String) Owner of the synced backups.
- `rate_in` (String) Inbound transfer rate limit in PBS byte-size format, for example `10M`.
- `rate_out` (String) Outbound transfer rate limit in PBS byte-size format.
- `burst_in` (String) Inbound burst rate limit in PBS byte-size format.
- `burst_out` (String) Outbound burst rate limit in PBS byte-size format.
- `comment` (String) Description for the sync job.

### Read-Only

- `digest` (String) Opaque digest returned by PBS for optimistic locking.

## Import

Import a sync job by ID:

```bash
terraform import pbs_sync_job.replicate nightly-sync
```
