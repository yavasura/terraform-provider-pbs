# Terraform HCL Tests for PBS Provider
#
# This directory contains native Terraform tests using the `.tftest.hcl` format.

## Directory Structure

```
test/tftest/
├── prune_job_datasource/
│   ├── main.tf              # Terraform configuration
│   └── test.tftest.hcl      # Test assertions
└── sync_job_datasource/
    ├── main.tf              # Terraform configuration
    └── test.tftest.hcl      # Test assertions
```

## Running Tests

### Prerequisites

1. **Terraform CLI v1.6.0+** (HCL tests require Terraform 1.6+)
2. **Built provider binary** in project root (`go build .`)
3. **PBS server running** and accessible
4. **Environment variables** set (see below)

### Environment Variables

Set these environment variables before running tests:

```bash
export TF_VAR_pbs_endpoint="https://pbs.example.com:8007"
export TF_VAR_pbs_username="root@pam"
export TF_VAR_pbs_password="your-password"
```

Alternatively, create a `terraform.tfvars` file in each test directory:

```hcl
pbs_endpoint = "https://pbs.example.com:8007"
pbs_username = "root@pam"
pbs_password = "your-password"
```

### Run All Tests

From the project root:

```bash
# Run all HCL tests (from project root)
(cd test/tftest/datastores_datasource && terraform test)
(cd test/tftest/prune_job_datasource && terraform test)
(cd test/tftest/prune_jobs_datasource && terraform test)
(cd test/tftest/sync_job_datasource && terraform test)

# Or run individual test by changing to its directory
cd test/tftest/prune_job_datasource
terraform test
```

### Run Specific Test

```bash
# Change to test directory and run
cd test/tftest/prune_job_datasource
terraform test

# Or from project root
(cd test/tftest/prune_job_datasource && terraform test)
```

### Debug Mode

Enable detailed logging:

```bash
export TF_LOG=DEBUG
terraform test -chdir=test/tftest/prune_job_datasource
```

## Test Coverage

These tests verify:

### Datastores Data Source (`datastores_datasource/`)
- ✅ Multiple datastore creation
- ✅ Datastores data source lists all datastores
- ✅ Test datastores appear in the list
- ✅ Datastore attributes correctly populated

### Prune Job Data Source (`prune_job_datasource/`)
- ✅ Datastore creation
- ✅ Prune job creation with datastore reference
- ✅ Prune job data source reads job correctly
- ✅ All attributes match between resource and data source

### Prune Jobs Data Source (`prune_jobs_datasource/`)
- ✅ Multiple datastore creation
- ✅ Multiple prune job creation
- ✅ Prune jobs data source lists all jobs
- ✅ Prune jobs data source filters by store
- ✅ Filtered results contain only specified store's jobs

### Sync Job Data Source (`sync_job_datasource/`)
- ✅ Datastore creation
- ✅ Remote creation
- ✅ Sync job creation with datastore and remote references
- ✅ Sync job data source reads job correctly
- ✅ All attributes match between resource and data source

## CI Integration

### GitHub Actions

Add an equivalent step to your CI workflow:

```yaml
- name: Run Terraform HCL Tests
  env:
    TF_VAR_pbs_endpoint: ${{ env.PBS_ENDPOINT }}
    TF_VAR_pbs_username: ${{ env.PBS_USERNAME }}
    TF_VAR_pbs_password: ${{ env.PBS_PASSWORD }}
  run: |
    echo "Running Terraform HCL tests..."
    terraform test -chdir=test/tftest/prune_job_datasource
    terraform test -chdir=test/tftest/sync_job_datasource
```

## Local Provider Setup

The tests expect the provider binary to be available. Two options:

### Option 1: Install to Terraform Plugin Directory

```bash
go build .
make install  # Copies to ~/.terraform.d/plugins/
```

### Option 2: Use terraform_override.tf (Development)

Create `terraform_override.tf` in each test directory:

```hcl
terraform {
  required_providers {
    pbs = {
      source = "registry.terraform.io/yavasura/pbs"
    }
  }
}

provider "pbs" {
  # Provider will be loaded from dev_overrides in ~/.terraformrc
}
```

And add to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/yavasura/pbs" = "/path/to/terraform-provider-pbs"
  }
  direct {}
}
```

## Future Tests

Additional tests to consider converting:
- Sync jobs data source (list with filters)
- Verify job data source
- Verify jobs data source
- Metrics server data source
- Metrics servers data source
- Notification endpoint/matcher data sources
- Remote stores/groups/namespaces data sources

## References

- [Terraform Test Framework](https://developer.hashicorp.com/terraform/language/tests)
- [Provider Development Testing](https://developer.hashicorp.com/terraform/plugin/testing)
