# PBS 4.0 Compatibility

This provider is intended to work with Proxmox Backup Server 4.0 and later.

## Current Status

The codebase and documentation already reflect the major PBS 4.0 API changes, including:

- metrics API path and schema changes
- notification endpoint path and payload changes
- verify, prune, and sync job schema updates
- datastore-level handling of garbage collection configuration

## Important Notes

- `pbs_metrics_server` uses the PBS 4.0 metrics API shape.
- notification resources and data sources use the PBS 4.0 `endpoints` and `matchers` APIs.
- `pbs_verify_job`, `pbs_prune_job`, and `pbs_sync_job` follow the PBS 4.0 job schemas, including digest handling where applicable.
- garbage collection is configured at the datastore level in PBS 4.0; there is no separate GC job resource in this provider.

## Related Documentation

- Detailed schema notes: [PBS_4.0_SCHEMA_CHANGES.md](./PBS_4.0_SCHEMA_CHANGES.md)
- Provider usage: [index.md](./index.md)
- Integration workflow: [INTEGRATION_TESTS.md](./INTEGRATION_TESTS.md)

## Maintainer Note

This document intentionally avoids historical failure reports. If compatibility regresses, prefer updating this file with the current state rather than appending old investigation logs.
