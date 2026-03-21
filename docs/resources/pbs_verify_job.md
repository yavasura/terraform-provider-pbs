---
page_title: "pbs_verify_job Resource - Proxmox Backup Server"
subcategory: ""
description: |-
Manages a PBS verify job.

Use [`pbs_namespace`](./pbs_namespace.md) to pre-create namespace paths before applying namespace-scoped verification jobs.
---

# pbs_verify_job

Manages a Proxmox Backup Server verify job.

Verify jobs periodically check backup integrity so corruption or storage issues can be detected before restore time.

## Example Usage

```terraform
resource "pbs_verify_job" "weekly" {
  id              = "weekly-verify"
  store           = pbs_datastore.backup.name
  schedule        = "weekly"
  ignore_verified = true
  outdated_after  = 7
  comment         = "Weekly integrity verification"
}
```

## Schema

### Required

- `id` (String) Unique verify job identifier.
- `store` (String) Datastore name where backups will be verified.
- `schedule` (String) When to run the job in systemd calendar format.

### Optional

- `ignore_verified` (Boolean) Skip backups that were verified recently.
- `outdated_after` (Number) Number of days after which a verified backup should be considered due for re-verification.
- `namespace` (String) Namespace to verify.
- `max_depth` (Number) Maximum recursion depth when traversing namespaces.
- `comment` (String) Description for the verify job.

### Read-Only

- `digest` (String) Opaque digest returned by PBS for optimistic locking.

## Import

Import a verify job by ID:

```bash
terraform import pbs_verify_job.weekly weekly-verify
```
