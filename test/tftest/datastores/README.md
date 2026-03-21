# Datastore HCL Tests

This directory contains HCL-based Terraform tests for datastore resources.

## Tests

### directory_datastore.tftest.hcl

Tests directory-backed datastore lifecycle:
- Create directory datastore
- Update mutable fields (comment, gc_schedule, etc.)
- Verify datastore configuration

## Related Tests

For datastore backend immutability coverage, see: `test/tftest/datastore_immutability/`
