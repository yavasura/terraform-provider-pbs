# Running Integration Tests

## Canonical Runner

The canonical acceptance and integration entrypoint is:

```bash
PBS_ENDPOINT=https://pbs.local:8007 \
PBS_USERNAME=root@pam \
PBS_PASSWORD=secret \
PBS_INSECURE=true \
./testacc
```

Use the lower-level commands in this document only for targeted debugging.

## Current Test Split

### Go integration tests

The Go suite under `test/integration/` covers live PBS-backed behavior for:

- datastores
- metrics
- remotes
- users
- cleanup helpers for external test resources

These tests use the shared Terraform CLI harness in `test/integration/setup.go`.

### Terraform-native HCL tests

The HCL suite under `test/tftest/` covers the broader resource and data source matrix, including:

- datastores and datastore immutability
- jobs and job data sources
- notifications and notification data sources
- remotes
- S3 endpoints and S3-backed datastore cases

## Prerequisites

Required:

- PBS 4.0+ instance
- admin credentials
- Go 1.24+
- Terraform in `PATH`

Optional for broader coverage:

- S3 credentials
- InfluxDB instance
- additional PBS test infrastructure

## Common Commands

### Package tests

```bash
go test ./...
```

### Provider-focused tests

```bash
go test ./internal/provider/...
```

### Go integration tests

```bash
export PBS_ENDPOINT="https://pbs-server:8007"
export PBS_USERNAME="root@pam"
export PBS_PASSWORD="secret"
export PBS_INSECURE="true"

go test -v -timeout 30m ./test/integration/...
```

### Focused Go integration runs

```bash
go test -v -timeout 30m ./test/integration/... -run "Datastore|Metrics|Remote|User"
```

### HCL tests

```bash
./test/tftest/run_hcl_tests.sh
```

## Recommendations

- For day-to-day work, start with package tests and targeted HCL tests.
- Run Go integration tests when changing live PBS interactions, import flows, or environment-sensitive behavior.
- Before release, run the full suite against a real PBS instance.

## FAQ

- **Are package tests enough for normal development?**
  Usually yes, especially together with HCL tests.

- **What do Go integration tests add?**
  Real API validation, Terraform CLI interactions, and end-to-end lifecycle coverage against PBS.

- **When should I prefer HCL tests?**
  When validating resource and data source behavior that is already modeled well in Terraform-native test cases.
