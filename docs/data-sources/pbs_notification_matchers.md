---
page_title: "pbs_notification_matchers Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Lists PBS notification matchers.
---

# pbs_notification_matchers

Lists notification matchers configured in Proxmox Backup Server.

## Example Usage

```terraform
data "pbs_notification_matchers" "all" {}
```

## Schema

### Read-Only

- `matchers` (Attributes List) List of notification matchers returned by PBS.

### Nested Schema for `matchers`

- `name` (String) Unique notification matcher name.
- `targets` (List of String) Notification endpoint names targeted by this matcher.
- `match_severity` (List of String) Matched severities.
- `match_field` (List of String) Matched `field=value` expressions.
- `match_calendar` (List of String) Matched calendar event IDs.
- `mode` (String) Matching mode, such as `all` or `any`.
- `invert_match` (Boolean) Whether the match logic is inverted.
- `comment` (String) Description for the matcher.
- `disable` (Boolean) Whether the matcher is disabled.
- `origin` (String) Origin of the matcher configuration.
