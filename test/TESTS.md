# Test Suite Documentation

This directory contains the Go-based integration tests for the provider.

## Primary commands

From the project root:

```bash
# Unit and package tests
make test-unit

# Acceptance / integration tests
PBS_ENDPOINT=https://pbs.local:8007 \
PBS_USERNAME=root@pam \
PBS_PASSWORD=secret \
PBS_INSECURE=true \
./testacc
```

For local Docker-backed integration runs:

```bash
./scripts/run-integration-tests.sh
```

For Terraform native HCL tests:

```bash
./test/tftest/run_hcl_tests.sh
```

## Scope

The current Go integration suite covers:

- datastores
- metrics
- remotes
- users
- cleanup helpers for external test resources

The Terraform-native HCL suite under `test/tftest/` covers the broader resource and data source matrix, including jobs, notifications, remotes, data sources, and S3 provider scenarios.

## Notes

- The Go integration tests use the shared test context in `test/integration/setup.go`.
- The HCL suite under `test/tftest/` complements the Go tests with native Terraform test execution.
- Prefer the current resource and data source set in the codebase over historical test docs.
