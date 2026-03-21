---
page_title: "pbs_notification_endpoints Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Lists PBS notification endpoints.
---

# pbs_notification_endpoints

Lists notification endpoints configured in Proxmox Backup Server. Returned endpoints can be of type `gotify`, `smtp`, `sendmail`, or `webhook`.

## Example Usage

```terraform
data "pbs_notification_endpoints" "all" {}
```

## Schema

### Read-Only

- `endpoints` (Attributes List) List of notification endpoints returned by PBS.

### Nested Schema for `endpoints`

- `name` (String) Unique notification endpoint name.
- `type` (String) Endpoint type.
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
- `body` (String) Webhook request body template.
- `method` (String) Webhook HTTP method.
- `headers` (Map of String) Custom webhook headers.
