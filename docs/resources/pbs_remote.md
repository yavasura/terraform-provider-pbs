---
page_title: "pbs_remote Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a PBS remote configuration.
---

# pbs_remote

Manages a Proxmox Backup Server remote configuration.

Remotes are typically referenced by sync jobs to pull or replicate backup data from another PBS instance.

## Example Usage

```terraform
resource "pbs_remote" "secondary_site" {
  name        = "secondary-site"
  host        = "pbs-secondary.example.com"
  port        = 8007
  auth_id     = "sync-user@pbs!replication"
  password    = var.remote_token_secret
  fingerprint = "AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99"
  comment     = "Secondary PBS used for replication"
}
```

## Schema

### Required

- `name` (String) Unique remote identifier.
- `host` (String) Hostname or IP address of the remote PBS server.
- `auth_id` (String) Authentication ID, for example `user@pam` or `user@pbs!token`.
- `password` (String, Sensitive) Password or authentication token for the remote.

### Optional

- `port` (Number) Port number for the remote PBS server. PBS commonly uses `8007`.
- `fingerprint` (String) X509 certificate fingerprint in SHA256 colon-separated format.
- `comment` (String) Description for the remote.

### Read-Only

- `digest` (String) Opaque digest returned by PBS for optimistic locking.

## Notes

The PBS API does not return the remote password on reads. Terraform still stores the configured password as a sensitive value in state so updates can be applied consistently.

## Import

Import a remote by name:

```bash
terraform import pbs_remote.secondary_site secondary-site
```
