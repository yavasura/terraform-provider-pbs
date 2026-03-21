---
page_title: "pbs_verify_jobs Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Lists PBS verify jobs.
---

# pbs_verify_jobs

Lists Proxmox Backup Server verify jobs, optionally filtered by datastore.

## Example Usage

### List all verify jobs

```terraform
data "pbs_verify_jobs" "all" {}
```

### Filter by datastore

```terraform
data "pbs_verify_jobs" "backup_store" {
  store = "backup-store"
}
```

## Schema

### Optional

- `store` (String) Filter verify jobs by datastore name.

### Read-Only

- `jobs` (Attributes List) List of verify jobs returned by PBS.

### Nested Schema for `jobs`

- `id` (String) Unique verify job identifier.
- `store` (String) Datastore name where verification is performed.
- `schedule` (String) Systemd calendar schedule for the job.
- `ignore_verified` (Boolean) Whether snapshots verified after `outdated_after` are skipped.
- `outdated_after` (Number) Days after which snapshots are considered outdated for verification.
- `namespace` (String) Namespace filter.
- `max_depth` (Number) Maximum namespace traversal depth.
- `comment` (String) Description for the verify job.
- `digest` (String) Opaque digest returned by PBS.
