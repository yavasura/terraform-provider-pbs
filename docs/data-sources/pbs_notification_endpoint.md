---
page_title: "pbs_notification_endpoint Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Reads a PBS notification endpoint configuration.
---

# pbs_notification_endpoint

Reads a notification endpoint from Proxmox Backup Server. Supported endpoint types are `gotify`, `smtp`, `sendmail`, and `webhook`.

## Example Usage

```terraform
data "pbs_notification_endpoint" "smtp_admins" {
  name = "smtp-admins"
  type = "smtp"
}
```

## Schema

### Required

- `name` (String) Unique notification endpoint name.
- `type` (String) Endpoint type: `gotify`, `smtp`, `sendmail`, or `webhook`.

### Read-Only

- `disable` (Boolean) Whether the endpoint is disabled.
- `comment` (String) Description for the endpoint.
- `origin` (String) Origin of the endpoint configuration.
- `server` (String) SMTP server address.
- `port` (Number) SMTP server port.
- `from_address` (String) From email address for SMTP or Sendmail.
- `mailto` (List of String) Recipient email addresses.
- `mailto_user` (List of String) PBS user IDs to notify.
- `mode` (String) SMTP connection mode.
- `username` (String) SMTP authentication username.
- `author` (String) Email author or sender name.
- `url` (String) Gotify or webhook target URL.
- `token` (String, Sensitive) Gotify application token.
- `secret` (String, Sensitive) Webhook secret used for signing.
- `body` (String) Webhook request body template.
- `method` (String) Webhook HTTP method.
- `headers` (Map of String) Custom webhook headers.
