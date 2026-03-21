---
page_title: "pbs_namespace Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a datastore namespace in Proxmox Backup Server.
---

# pbs_namespace

Manages a datastore namespace in Proxmox Backup Server.

Namespaces are hierarchical and must be created parent-first. PBS allows up to 7 path components,
with each component matching `[A-Za-z0-9_][A-Za-z0-9._-]*`.

## Example Usage

### Hierarchical Namespace

```terraform
resource "pbs_namespace" "production" {
  store     = "backups"
  namespace = "production"
}

resource "pbs_namespace" "production_vms" {
  store     = "backups"
  namespace = "production/vms"
  depends_on = [pbs_namespace.production]
}
```

### Multi-Tenant Layout

```terraform
resource "pbs_namespace" "tenant_a" {
  store     = "backups"
  namespace = "tenants/tenant-a"
}

resource "pbs_namespace" "tenant_b" {
  store     = "backups"
  namespace = "tenants/tenant-b"
}
```

Use `pbs_namespace` to pre-create namespaces referenced by job resources such as
`pbs_prune_job`, `pbs_sync_job`, and `pbs_verify_job`.

## Schema

### Required

- `store` (String) Datastore that owns the namespace.
- `namespace` (String) Namespace path, for example `production/vms`.

### Optional

- `comment` (String) Optional namespace comment.

## Import

Import a namespace by `store:namespace`:

```bash
terraform import pbs_namespace.production_vms backups:production/vms
```
