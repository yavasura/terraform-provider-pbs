---
page_title: "pbs_webhook_notification Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a PBS webhook notification target.
---

# pbs_webhook_notification

Manages a Proxmox Backup Server webhook notification target.

## Example Usage

```terraform
resource "pbs_webhook_notification" "incident_hook" {
  name    = "incident-hook"
  url     = "https://hooks.example.com/pbs"
  method  = "post"
  headers = {
    X-Service = "pbs"
  }
  secret  = var.webhook_secret
  comment = "Webhook for incident routing"
}
```

## Schema

### Required

- `name` (String) Unique webhook notification target name.
- `url` (String) Webhook URL.

### Optional

- `body` (String) Custom request body template.
- `method` (String) HTTP method. One of `post` or `put`. Defaults to `post`.
- `headers` (Map of String) Custom HTTP headers.
- `secret` (String, Sensitive) Secret used for HMAC-SHA256 signing.
- `comment` (String) Description for the notification target.
- `disable` (Boolean) Disable the target without deleting it.

### Read-Only

- `origin` (String) Origin reported by PBS for this configuration.

## Import

Import a webhook notification target by name:

```bash
terraform import pbs_webhook_notification.incident_hook incident-hook
```
