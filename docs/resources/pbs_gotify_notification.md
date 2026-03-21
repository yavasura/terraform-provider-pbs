---
page_title: "pbs_gotify_notification Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a PBS Gotify notification target.
---

# pbs_gotify_notification

Manages a Proxmox Backup Server Gotify notification target.

## Example Usage

```terraform
resource "pbs_gotify_notification" "mobile_alerts" {
  name    = "mobile-alerts"
  server  = "https://gotify.example.com"
  token   = var.gotify_token
  comment = "Push alerts to the operations team"
}
```

## Schema

### Required

- `name` (String) Unique Gotify notification target name.
- `server` (String) Gotify server URL.
- `token` (String, Sensitive) Gotify application token.

### Optional

- `comment` (String) Description for the notification target.
- `disable` (Boolean) Disable the target without deleting it.

### Read-Only

- `origin` (String) Origin reported by PBS for this configuration.

## Import

Import a Gotify notification target by name:

```bash
terraform import pbs_gotify_notification.mobile_alerts mobile-alerts
```
