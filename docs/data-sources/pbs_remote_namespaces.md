---
page_title: "pbs_remote_namespaces Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Lists namespaces in a remote PBS datastore.
---

# pbs_remote_namespaces

Scans a datastore on a configured remote Proxmox Backup Server and returns the available namespace paths.

## Example Usage

```terraform
data "pbs_remote_namespaces" "remote_a_backup" {
  remote_name = "remote-a"
  store       = "backup-store"
}
```

## Schema

### Required

- `remote_name` (String) Name of the configured remote to scan.
- `store` (String) Datastore name on the remote server.

### Read-Only

- `namespaces` (List of String) Namespace paths available in the remote datastore.
