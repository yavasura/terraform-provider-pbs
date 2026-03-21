# Integration Tests

This directory contains integration and test-support material for the PBS Terraform provider.

## Prerequisites

1. Install Go 1.24+ and Terraform 1.6+.
2. Build the provider binary from the project root:

```bash
go build .
```

3. Have a reachable Proxmox Backup Server instance.

## Required Environment Variables

```bash
export PBS_ENDPOINT="https://your-pbs-server:8007"
export PBS_USERNAME="root@pam"
export PBS_PASSWORD="your-password"
export PBS_INSECURE="true"  # Optional, useful for self-signed certs
```

## Running Tests

From the project root:

```bash
# Unit tests only
make test-unit

# Integration tests / acceptance tests
PBS_ENDPOINT=https://pbs.local:8007 \
PBS_USERNAME=root@pam \
PBS_PASSWORD=secret \
PBS_INSECURE=true \
./testacc
```

`./testacc` is the canonical acceptance/integration entrypoint. Other helpers should delegate to it or exist only for special local workflows.

Run a specific integration test directly:

```bash
go test -v -timeout 30m ./test/integration -run TestDatastore
```

## Docker Workflow

For local development, you can use the helper script:

```bash
./scripts/run-integration-tests.sh
```

Useful variations:

```bash
./scripts/run-integration-tests.sh TestQuickSmoke
```

The Docker-based setup is intended for local development convenience. The standard acceptance/integration harness still uses the `PBS_*` variables listed above.

## Optional S3 Test Variables

Real S3 provider coverage can be enabled with provider-specific credentials, for example:

```bash
export AWS_ACCESS_KEY_ID="your-aws-key"
export AWS_SECRET_ACCESS_KEY="your-aws-secret"
export AWS_REGION="us-east-1"
```

Equivalent variables can be supplied for other S3-compatible providers used by the test suite.
