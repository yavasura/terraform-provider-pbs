package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRemoteDataSources tests the remote data sources for stores, namespaces, and groups
func TestRemoteDataSources(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := SetupTest(t)
	defer tc.DestroyTerraform(t)

	remoteName := GenerateTestName("ds-remote")

	// Note: These data sources require a real remote PBS server to scan
	// For testing purposes, we create the remote configuration but don't expect
	// the scan to succeed (would require actual remote server)
	config := fmt.Sprintf(`
resource "pbs_remote" "data_source_test" {
  name     = "%s"
  host     = "pbs.example.com"
  auth_id  = "sync@pbs!test-token"
  password = "test-password"
}

# These data sources would work with a real remote server
# In testing, they may fail if the remote is not accessible
# data "pbs_remote_stores" "test_stores" {
#   remote_name = pbs_remote.data_source_test.name
# }
#
# data "pbs_remote_namespaces" "test_namespaces" {
#   remote_name = pbs_remote.data_source_test.name
#   store       = "datastore1"
# }
#
# data "pbs_remote_groups" "test_groups" {
#   remote_name = pbs_remote.data_source_test.name
#   store       = "datastore1"
#   namespace   = "production"
# }

output "remote_name" {
  value = pbs_remote.data_source_test.name
}
`, remoteName)

	tc.WriteMainTF(t, config)
	tc.ApplyTerraform(t)

	// Verify remote was created
	resource := tc.GetResourceFromState(t, "pbs_remote.data_source_test")
	assert.Equal(t, remoteName, resource.AttributeValues["name"])

	t.Log("INFO: Data source tests skipped - require real remote PBS server")
}
