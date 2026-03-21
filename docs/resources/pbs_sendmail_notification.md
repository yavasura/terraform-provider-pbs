---
page_title: "pbs_sendmail_notification Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a PBS Sendmail notification target.
---

# pbs_sendmail_notification

Manages a Proxmox Backup Server Sendmail notification target.

## Example Usage

```terraform
resource "pbs_sendmail_notification" "local_mail" {
  name         = "local-mail"
  from_address = "pbs@example.com"
  mailto       = ["ops@example.com"]
  comment      = "Local sendmail notifications"
}
```

## Schema

### Required

- `name` (String) Unique Sendmail notification target name.
- `from_address` (String) Sender email address.

### Optional

- `mailto` (List of String) Recipient email addresses.
- `mailto_user` (List of String) PBS user IDs that should receive notifications.
- `author` (String) Author name shown in email headers.
- `comment` (String) Description for the notification target.
- `disable` (Boolean) Disable the target without deleting it.

### Read-Only

- `origin` (String) Origin reported by PBS for this configuration.

## Import

Import a Sendmail notification target by name:

```bash
terraform import pbs_sendmail_notification.local_mail local-mail
```
