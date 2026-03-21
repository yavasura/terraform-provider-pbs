---
page_title: "pbs_smtp_notification Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a PBS SMTP notification target.
---

# pbs_smtp_notification

Manages a Proxmox Backup Server SMTP notification target.

SMTP targets can be referenced directly by PBS notification matchers and endpoint groups.

## Example Usage

```terraform
resource "pbs_smtp_notification" "ops_mail" {
  name         = "ops-mail"
  server       = "smtp.example.com"
  port         = 587
  mode         = "starttls"
  from_address = "pbs@example.com"
  mailto       = ["ops@example.com"]
  username     = "pbs@example.com"
  password     = var.smtp_password
  comment      = "Primary SMTP notification target"
}
```

## Schema

### Required

- `name` (String) Unique SMTP notification target name.
- `server` (String) SMTP server hostname or IP address.
- `from_address` (String) Sender email address.

### Optional

- `port` (Number) SMTP server port. Defaults to `25`.
- `mode` (String) SMTP connection mode. One of `insecure`, `starttls`, `tls`. Defaults to `insecure`.
- `mailto` (List of String) Recipient email addresses.
- `mailto_user` (List of String) PBS user IDs that should receive notifications.
- `username` (String) SMTP authentication username.
- `password` (String, Sensitive) SMTP authentication password.
- `author` (String) Author name shown in email headers.
- `comment` (String) Description for the notification target.
- `disable` (Boolean) Disable the target without deleting it.

### Read-Only

- `origin` (String) Origin reported by PBS for this configuration.

## Import

Import an SMTP notification target by name:

```bash
terraform import pbs_smtp_notification.ops_mail ops-mail
```
