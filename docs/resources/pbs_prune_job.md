---
page_title: "pbs_prune_job Resource - Proxmox Backup Server"
subcategory: ""
description: |-
Manages a PBS prune job.

Use [`pbs_namespace`](./pbs_namespace.md) to create hierarchical namespace targets before applying namespace-scoped prune rules.
---

# pbs_prune_job

Manages a Proxmox Backup Server prune job.

Prune jobs enforce retention policies by removing older backup snapshots on a schedule.

## Example Usage

```terraform
resource "pbs_prune_job" "daily" {
  id          = "daily-prune"
  store       = pbs_datastore.backup.name
  schedule    = "daily"
  keep_last   = 7
  keep_daily  = 14
  keep_weekly = 8
  comment     = "Default retention policy"
}
```

## Schema

### Required

- `id` (String) Unique prune job identifier.
- `store` (String) Datastore name where pruning will run.
- `schedule` (String) When to run the job in systemd calendar format, for example `daily` or `Mon..Fri *-*-* 02:00:00`.

### Optional

- `keep_last` (Number) Keep the last N snapshots regardless of age.
- `keep_hourly` (Number) Keep hourly snapshots for the last N hours.
- `keep_daily` (Number) Keep daily snapshots for the last N days.
- `keep_weekly` (Number) Keep weekly snapshots for the last N weeks.
- `keep_monthly` (Number) Keep monthly snapshots for the last N months.
- `keep_yearly` (Number) Keep yearly snapshots for the last N years.
- `max_depth` (Number) Maximum namespace traversal depth for pruning.
- `namespace` (String) Namespace filter applied when pruning.
- `comment` (String) Description for the prune job.
- `disable` (Boolean) Disable the prune job without deleting it.

### Read-Only

- `digest` (String) Opaque digest returned by PBS for optimistic locking.

## Import

Import a prune job by ID:

```bash
terraform import pbs_prune_job.daily daily-prune
```
