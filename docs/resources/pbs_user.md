---
page_title: "pbs_user Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a PBS user account.
---

# pbs_user

Manages a PBS user account through the `/access/users` API.

This resource manages account metadata only. It does not manage passwords or other credentials.
For PAM users, passwords are managed on the PBS host itself. For external realms such as LDAP, AD,
and OpenID, authentication continues to be controlled by the upstream identity provider.

## Example Usage

```terraform
resource "pbs_user" "backup_operator" {
  userid    = "backup-operator@ldap"
  comment   = "Backup operations team account"
  enable    = true
  firstname = "Backup"
  lastname  = "Operator"
  email     = "backup-operator@example.com"
}
```

## Schema

### Required

- `userid` (String) PBS user ID in `username@realm` format, for example `admin@pam` or `john@ldap`.

### Optional

- `comment` (String) Optional comment for the user account.
- `email` (String) Email address associated with the user.
- `enable` (Boolean) Whether the account is enabled.
- `expire` (Number) Account expiration time as a Unix timestamp.
- `firstname` (String) Given name for the user.
- `lastname` (String) Family name for the user.

### Read-Only

- `digest` (String) Opaque digest returned by PBS for optimistic locking.

## Import

Import a PBS user by its user ID:

```bash
terraform import pbs_user.backup_operator backup-operator@ldap
```
