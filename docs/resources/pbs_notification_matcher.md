---
page_title: "pbs_notification_matcher Resource - Proxmox Backup Server"
subcategory: ""
description: |-
  Manages a PBS notification matcher.
---

# pbs_notification_matcher

Manages a Proxmox Backup Server notification matcher.

Notification matchers define routing rules that decide which targets receive which notification events.

## Example Usage

```terraform
resource "pbs_notification_matcher" "critical_only" {
  name           = "critical-only"
  targets        = ["ops-mail", "incident-hook"]
  match_severity = ["error", "warning"]
  mode           = "any"
  comment        = "Route higher-severity notifications to ops"
}
```

## Schema

### Required

- `name` (String) Unique notification matcher name.

### Optional

- `targets` (List of String) Notification target names that should receive matching events.
- `match_severity` (List of String) Severity levels to match. Supported values are `info`, `notice`, `warning`, `error`.
- `match_field` (List of String) Field match expressions in `field=value` form.
- `match_calendar` (List of String) Calendar IDs for time-based routing.
- `mode` (String) Matching mode. One of `all` or `any`. Defaults to `all`.
- `invert_match` (Boolean) Invert the match logic.
- `comment` (String) Description for the matcher.
- `disable` (Boolean) Disable the matcher without deleting it.

### Read-Only

- `origin` (String) Origin reported by PBS for this configuration.

## Import

Import a notification matcher by name:

```bash
terraform import pbs_notification_matcher.critical_only critical-only
```
