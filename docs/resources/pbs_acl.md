---
page_title: "pbs_acl Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a PBS access control list entry.
---

# pbs_acl

Manages a PBS ACL entry for a user, token, or group on a PBS path.

## Example Usage

### User ACL

```terraform
resource "pbs_acl" "backup_admin" {
  path      = "/datastore/backups"
  ugid      = "backup-operator@pbs"
  role_id   = "DatastoreAdmin"
  propagate = true
}
```

### Token ACL

```terraform
resource "pbs_acl" "terraform_token" {
  path      = "/remote/secondary"
  ugid      = "automation@pbs!terraform"
  role_id   = "RemoteAdmin"
  propagate = false
}
```

## Role Reference

| Role | Typical Use |
| --- | --- |
| `Admin` | Full administrative access |
| `Audit` | Read-only auditing |
| `DatastoreAdmin` | Full control of a datastore |
| `DatastoreBackup` | Backup/restore oriented access |
| `DatastorePowerUser` | Elevated datastore operations |
| `DatastoreReader` | Read-only datastore access |
| `NoAccess` | Explicit deny |
| `RemoteAdmin` | Manage remotes |
| `RemoteSyncOperator` | Operate remote sync workflows |
| `TapeAdmin` | Full tape administration |
| `TapeAudit` | Read-only tape audit access |
| `TapeOperator` | Tape operation access |

## Path Examples

- `/`
- `/datastore`
- `/datastore/backups`
- `/remote/secondary`
- `/tape`

## Security Notes

- Prefer least-privilege roles such as `DatastoreReader` or `RemoteSyncOperator` over broad administrative roles.
- More specific paths override broader ones.
- `propagate = false` stops inheritance to child paths.
- `NoAccess` can be used as an explicit deny.

## Schema

### Required

- `path` (String) PBS ACL path.
- `ugid` (String) User, token, or group identifier.
- `role_id` (String) PBS role to assign.
- `propagate` (Boolean) Whether the ACL propagates to child paths.

## Import

Import an ACL by `{path}:{ugid}`:

```bash
terraform import pbs_acl.backup_admin /datastore/backups:backup-operator@pbs
```
