# Test Suites

This repository currently keeps release automation in GitHub Actions and runs provider tests primarily through local or self-managed environments.

## Canonical test entrypoints

- Unit and package tests:
  ```bash
  make test-unit
  ```

- Acceptance / integration tests:
  ```bash
  PBS_ENDPOINT=https://pbs.local:8007 \
  PBS_USERNAME=root@pam \
  PBS_PASSWORD=secret \
  PBS_INSECURE=true \
  ./testacc
  ```

- Terraform HCL tests:
  ```bash
  ./test/tftest/run_hcl_tests.sh
  ```

- Docker-assisted local integration workflow:
  ```bash
  ./scripts/run-integration-tests.sh
  ```

## Notes

- `./testacc` is the canonical acceptance/integration entrypoint.
- `./scripts/run-integration-tests.sh` is a convenience wrapper for local Docker-backed testing.
- There is currently no checked-in GitHub Actions integration-test workflow in this repository.
