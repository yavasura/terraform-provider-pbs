---
page_title: "pbs_api_token Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a PBS API token for an existing user.
---

# pbs_api_token

Manages a PBS API token for an existing user through the `/access/users/<userid>/token/...` API.

PBS returns the token secret only when the token is created. This resource exposes that one-time
secret as the sensitive `value` attribute. Imported tokens do not have a retrievable secret, so
`value` remains unset for imported resources.

## Example Usage

```terraform
resource "pbs_user" "automation" {
  userid  = "automation@pbs"
  comment = "Automation account"
  enable  = true
}

resource "pbs_api_token" "terraform" {
  userid     = pbs_user.automation.userid
  token_name = "terraform"
  comment    = "Terraform automation token"
  enable     = true
}
```

## Schema

### Required

- `userid` (String) PBS user ID that owns the token, for example `automation@pbs`.
- `token_name` (String) Token name segment, for example `terraform` in `automation@pbs!terraform`.

### Optional

- `comment` (String) Optional token comment.
- `enable` (Boolean) Whether the token is enabled.
- `expire` (Number) Token expiration time as a Unix timestamp.

### Read-Only

- `digest` (String) Opaque digest returned by PBS for optimistic locking.
- `tokenid` (String) Full PBS token ID in `userid!token_name` format.
- `value` (String, Sensitive) One-time token secret returned by PBS when the token is created.

## Import

Import a PBS API token by its full token ID:

```bash
terraform import pbs_api_token.terraform automation@pbs!terraform
```
