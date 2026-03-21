# Terraform Provider for Proxmox Backup Server

[![License](https://img.shields.io/badge/License-MPL--2.0-blue.svg)](https://opensource.org/licenses/MPL-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/yavasura/terraform-provider-pbs)](https://goreportcard.com/report/github.com/yavasura/terraform-provider-pbs)

Terraform / OpenTofu provider for managing [Proxmox Backup Server](https://www.proxmox.com/en/proxmox-backup-server).

## Compatibility

- Terraform >= 1.0
- Go >= 1.24 for local builds
- Proxmox Backup Server >= 2.0

The provider currently supports core PBS configuration domains including datastores, S3 endpoints, remotes, notifications, metrics, jobs, and users.

## Install

### From the registry

```hcl
terraform {
  required_providers {
    pbs = {
      source  = "yavasura/pbs"
      version = "~> 1.0"
    }
  }
}
```

### Provider configuration

```hcl
provider "pbs" {
  endpoint = "https://pbs.example.com:8007"
  username = "admin@pbs"
  password = var.pbs_password
  insecure = false
}
```

Environment variables are also supported:

```bash
export PBS_ENDPOINT="https://pbs.example.com:8007"
export PBS_USERNAME="admin@pbs"
export PBS_PASSWORD="your-password"
export PBS_API_TOKEN="user@realm!token=secret"
export PBS_INSECURE="false"
```

## Example

```hcl
resource "pbs_datastore" "backup" {
  name    = "backup-store"
  path    = "/mnt/datastore/backup"
  comment = "Primary backup datastore"
}
```

More examples are available in:

- [examples/provider/provider.tf](examples/provider/provider.tf)
- [examples/resources](examples/resources)
- [examples/data-sources](examples/data-sources)

## Local Development

### Build and install

```bash
make build
make install
```

This installs the provider into:

`~/.terraform.d/plugins/registry.terraform.io/yavasura/pbs/${VERSION}/${GOOS}_${GOARCH}/`

### Terraform CLI config

For local provider development, copy [example.tfrc](example.tfrc) to `~/.terraformrc` and adjust the path if needed.

### Testing

```bash
# Unit and package tests
make test-unit

# Acceptance / integration tests
PBS_ENDPOINT=https://pbs.local:8007 \
PBS_USERNAME=root@pam \
PBS_PASSWORD=secret \
PBS_INSECURE=true \
./testacc

# Native Terraform HCL tests
./test/tftest/run_hcl_tests.sh
```

Additional integration testing notes live in:

- [docs/INTEGRATION_TESTS.md](docs/INTEGRATION_TESTS.md)
- [test/README.md](test/README.md)
- [test/tftest/README.md](test/tftest/README.md)

## Release Workflow

Release automation is scaffolded with:

- `.github/workflows/release-please.yml`
- `release-please-config.json`
- `.release-please-manifest.json`
- `.goreleaser.yml`

`release-please` prepares version PRs and changelog updates, and the tag-based release workflow builds and publishes release artifacts.

## Resources

- `pbs_acl`
- `pbs_api_token`
- `pbs_datastore`
- `pbs_namespace`
- `pbs_s3_endpoint`
- `pbs_remote`
- `pbs_metrics_server`
- `pbs_smtp_notification`
- `pbs_gotify_notification`
- `pbs_sendmail_notification`
- `pbs_webhook_notification`
- `pbs_notification_matcher`
- `pbs_prune_job`
- `pbs_sync_job`
- `pbs_verify_job`
- `pbs_user`

Full resource docs are available in [docs/resources](docs/resources).

## Data Sources

- `pbs_datastore`
- `pbs_datastores`
- `pbs_namespaces`
- `pbs_prune_job`
- `pbs_prune_jobs`
- `pbs_sync_job`
- `pbs_sync_jobs`
- `pbs_verify_job`
- `pbs_verify_jobs`
- `pbs_metrics_server`
- `pbs_metrics_servers`
- `pbs_s3_endpoint`
- `pbs_s3_endpoints`
- `pbs_remote_stores`
- `pbs_remote_namespaces`
- `pbs_remote_groups`
- `pbs_notification_endpoint`
- `pbs_notification_endpoints`
- `pbs_notification_matcher`
- `pbs_notification_matchers`

Full data source docs are available in [docs/data-sources](docs/data-sources).

For the Terraform Registry view, see:

https://registry.terraform.io/providers/yavasura/pbs/latest/docs

## Support

- Issues: https://github.com/yavasura/terraform-provider-pbs/issues
- Discussions: https://github.com/yavasura/terraform-provider-pbs/discussions

## License

Mozilla Public License 2.0. See [LICENSE](LICENSE).
