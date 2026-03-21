---
page_title: "pbs_remote_groups Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Lists backup groups in a remote PBS namespace.
---

# pbs_remote_groups

Scans a namespace in a datastore on a configured remote Proxmox Backup Server and returns the available backup groups.

## Example Usage

```terraform
data "pbs_remote_groups" "remote_a_prod" {
  remote_name = "remote-a"
  store       = "backup-store"
  namespace   = "production"
}
```

## Schema

### Required

- `remote_name` (String) Name of the configured remote to scan.
- `store` (String) Datastore name on the remote server.

### Optional

- `namespace` (String) Namespace path to scan. Defaults to the root namespace.

### Read-Only

- `groups` (List of String) Backup group identifiers in `type/id` format.
