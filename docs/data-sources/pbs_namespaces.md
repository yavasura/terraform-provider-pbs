---
page_title: "pbs_namespaces Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Lists namespaces in a Proxmox Backup Server datastore.
---

# pbs_namespaces

Lists namespaces in a Proxmox Backup Server datastore.

## Example Usage

```terraform
data "pbs_namespaces" "production" {
  store     = "backups"
  prefix    = "production"
  max_depth = 2
}
```

## Schema

### Required

- `store` (String) Datastore to inspect.

### Optional

- `prefix` (String) Optional namespace prefix filter.
- `max_depth` (Number) Optional maximum namespace depth to request from PBS.

### Read-Only

- `namespaces` (List of Object) Matching namespaces.
- `namespaces[].namespace` (String) Namespace path.
- `namespaces[].comment` (String) Namespace comment.
