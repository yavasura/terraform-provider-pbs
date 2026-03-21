---
page_title: "pbs_verify_job Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Reads a PBS verify job.
---

# pbs_verify_job

Reads an existing Proxmox Backup Server verify job by ID.

## Example Usage

```terraform
data "pbs_verify_job" "weekly" {
  id = "weekly-verify"
}
```

## Schema

### Required

- `id` (String) Unique verify job identifier.

### Read-Only

- `store` (String) Datastore name where verification is performed.
- `schedule` (String) Systemd calendar schedule for the job.
- `ignore_verified` (Boolean) Whether snapshots verified after `outdated_after` are skipped.
- `outdated_after` (Number) Days after which snapshots are considered outdated for verification.
- `namespace` (String) Namespace filter.
- `max_depth` (Number) Maximum namespace traversal depth.
- `comment` (String) Description for the verify job.
- `digest` (String) Opaque digest returned by PBS.
