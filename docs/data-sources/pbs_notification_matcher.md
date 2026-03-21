---
page_title: "pbs_notification_matcher Data Source - Proxmox Backup Server"
subcategory: ""
description: |-
  Reads a PBS notification matcher.
---

# pbs_notification_matcher

Reads a notification matcher from Proxmox Backup Server. Matchers define routing rules that send notifications to selected endpoints based on severity, field matches, or calendar matches.

## Example Usage

```terraform
data "pbs_notification_matcher" "errors_only" {
  name = "errors-only"
}
```

## Schema

### Required

- `name` (String) Unique notification matcher name.

### Read-Only

- `targets` (List of String) Notification endpoint names targeted by this matcher.
- `match_severity` (List of String) Matched severities.
- `match_field` (List of String) Matched `field=value` expressions.
- `match_calendar` (List of String) Matched calendar event IDs.
- `mode` (String) Matching mode, such as `all` or `any`.
- `invert_match` (Boolean) Whether the match logic is inverted.
- `comment` (String) Description for the matcher.
- `disable` (Boolean) Whether the matcher is disabled.
- `origin` (String) Origin of the matcher configuration.
