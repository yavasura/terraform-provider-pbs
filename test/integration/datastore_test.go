package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDatastoreValidation tests validation scenarios
func TestDatastoreValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := SetupTest(t)
	defer tc.DestroyTerraform(t)

	// Test missing required path for directory datastore
	invalidDirConfig := `
resource "pbs_datastore" "invalid_dir" {
  name = "invalid-dir"
  # missing required path
}
`

	tc.WriteMainTF(t, invalidDirConfig)
	err := tc.ApplyTerraformWithError(t)
	assert.Error(t, err, "Should fail validation for directory datastore without path")

	// Test missing S3 bucket when client is provided
	invalidS3Config := `
resource "pbs_datastore" "invalid_s3" {
  name      = "invalid-s3"
  s3_client = "endpoint-1"
  path      = "/datastore/invalid-s3"
}
`

	tc.WriteMainTF(t, invalidS3Config)
	err = tc.ApplyTerraformWithError(t)
	assert.Error(t, err, "Should fail validation when only one S3 attribute is provided")

	// Test removable datastore without backing device
	invalidRemovableConfig := `
resource "pbs_datastore" "invalid_removable" {
  name       = "invalid-removable"
  path       = "/datastore/invalid-removable"
  removable  = true
}
`

	tc.WriteMainTF(t, invalidRemovableConfig)
	err = tc.ApplyTerraformWithError(t)
	assert.Error(t, err, "Should fail validation when removable datastore lacks backing_device")

}
