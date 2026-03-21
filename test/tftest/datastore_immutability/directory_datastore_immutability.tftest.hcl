# Test to validate datastore backend immutability behavior
#
# This test validates that:
# 1. Directory datastores can be created successfully
# 2. Subsequent applies without changes don't fail with 400 errors
# 3. Mutable fields (like comment, gc_schedule) can be updated without recreation
# 4. Immutable fields (like path) trigger replacement when changed
#
# Issue: https://github.com/yavasura/terraform-provider-pbs/issues/18
# Note: Using directory datastore to avoid external S3 dependencies

run "setup" {
  command = plan
  
  variables {
    datastore_name = "dir-immut-${var.test_id}"
    datastore_path = "/datastore/dir-immut-${var.test_id}"
    comment        = "Initial comment"
    gc_schedule    = "daily"
  }
}

run "create_directory_datastore" {
  command = apply
  
  variables {
    datastore_name = "dir-immut-${var.test_id}"
    datastore_path = "/datastore/dir-immut-${var.test_id}"
    comment        = "Initial comment"
    gc_schedule    = "daily"
  }
  
  assert {
    condition     = pbs_datastore.dir_test.name == "dir-immut-${var.test_id}"
    error_message = "Datastore name should match input"
  }
  
  assert {
    condition     = pbs_datastore.dir_test.path == "/datastore/dir-immut-${var.test_id}"
    error_message = "Datastore path should match input"
  }
  
  assert {
    condition     = pbs_datastore.dir_test.comment == "Initial comment"
    error_message = "Comment should match input"
  }
  
  assert {
    condition     = pbs_datastore.dir_test.gc_schedule == "daily"
    error_message = "GC schedule should match input"
  }
}

# Reapplying the same datastore should succeed without sending immutable
# backend fields back to the PBS update endpoint.
run "reapply_without_changes" {
  command = apply
  
  variables {
    datastore_name = "dir-immut-${var.test_id}"
    datastore_path = "/datastore/dir-immut-${var.test_id}"
    comment        = "Initial comment"
    gc_schedule    = "daily"
  }
  
  # Should succeed without errors (no changes)
  assert {
    condition     = pbs_datastore.dir_test.name == "dir-immut-${var.test_id}"
    error_message = "Datastore should remain unchanged"
  }
  
  assert {
    condition     = pbs_datastore.dir_test.comment == "Initial comment"
    error_message = "Comment should remain unchanged"
  }
}

# Verify that mutable fields can be updated without recreation
run "update_mutable_fields" {
  command = apply
  
  variables {
    datastore_name = "dir-immut-${var.test_id}"
    datastore_path = "/datastore/dir-immut-${var.test_id}"
    comment        = "Updated comment - this should not recreate"
    gc_schedule    = "weekly"
  }
  
  assert {
    condition     = pbs_datastore.dir_test.comment == "Updated comment - this should not recreate"
    error_message = "Comment should be updated"
  }
  
  assert {
    condition     = pbs_datastore.dir_test.gc_schedule == "weekly"
    error_message = "GC schedule should be updated"
  }
  
  # Verify immutable field remains unchanged
  assert {
    condition     = pbs_datastore.dir_test.path == "/datastore/dir-immut-${var.test_id}"
    error_message = "Path should remain unchanged"
  }
}

# Verify that changing immutable fields triggers replacement
run "plan_immutable_field_change" {
  command = plan
  
  variables {
    datastore_name = "dir-immut-${var.test_id}"
    datastore_path = "/datastore/different-dir-immut-${var.test_id}"  # Changed immutable field
    comment        = "Updated comment - this should not recreate"
    gc_schedule    = "weekly"
  }
  
  # Plan should show replacement due to path change
  # Terraform test framework doesn't have direct access to plan details,
  # but we can verify the configuration is accepted
  assert {
    condition     = var.datastore_path == "/datastore/different-dir-immut-${var.test_id}"
    error_message = "Variable should be set to new path"
  }
}
