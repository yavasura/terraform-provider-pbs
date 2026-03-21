# Terraform HCL Tests - Quick Start

## What Are These?

These are **native Terraform tests** using the `.tftest.hcl` format introduced in Terraform v1.6.0.

## Prerequisites

- **Terraform v1.6.0+** (HCL tests require v1.6+)
- **Built provider binary** (`go build .`)
- **PBS server** running and accessible

## Quick Start

### 1. Set Environment Variables

```bash
export TF_VAR_pbs_endpoint="https://pbs.example.com:8007"
export TF_VAR_pbs_username="root@pam"
export TF_VAR_pbs_password="your-password"
```

### 2. Run Tests

```bash
# Run all HCL tests
./test/tftest/run_hcl_tests.sh

# Run specific test
./test/tftest/run_hcl_tests.sh prune_job_datasource
```

## Test Coverage

✅ **Datastores Data Source** (`datastores_datasource/`)
- Creates multiple datastores
- Lists all datastores via data source
- Verifies test datastores appear in list
- Checks datastore attributes

✅ **Prune Job Data Source** (`prune_job_datasource/`)
- Creates datastore and prune job
- Reads job via data source
- Verifies all attributes match

✅ **Prune Jobs Data Source** (`prune_jobs_datasource/`)
- Creates multiple datastores and prune jobs
- Lists all jobs via data source
- Filters jobs by store
- Verifies filtering works correctly

✅ **Sync Job Data Source** (`sync_job_datasource/`)
- Creates datastore, remote, and sync job
- Reads job via data source  
- Verifies all attributes match

## Documentation

- **[README.md](README.md)** - Detailed usage instructions

## CI Integration

These tests run automatically in GitHub Actions after the Go integration tests:

```yaml
- name: Run Terraform HCL tests
  run: |
    (cd test/tftest/datastores_datasource && terraform test)
    (cd test/tftest/prune_job_datasource && terraform test)
    (cd test/tftest/prune_jobs_datasource && terraform test)
    (cd test/tftest/sync_job_datasource && terraform test)
```

## Debugging

Enable detailed logging:

```bash
export TF_LOG=DEBUG
terraform test -chdir=test/tftest/prune_job_datasource
```

## Status

- ✅ All 4 HCL tests passing reliably in CI
- ✅ Cover the same core datasource scenarios as the skipped Go tests
- ✅ Use the native Terraform test framework
