---
page_title: "pbs_remote_stores Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Lists available datastores on a remote PBS server.
---

# pbs_remote_stores

Scans a configured remote Proxmox Backup Server and returns the datastore names available on that remote.

## Example Usage

```terraform
data "pbs_remote_stores" "remote_a" {
  remote_name = "remote-a"
}
```

## Schema

### Required

- `remote_name` (String) Name of the configured remote to scan.

### Read-Only

- `stores` (List of String) Datastore names available on the remote server.
