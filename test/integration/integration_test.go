package integration

import (
	"testing"
)

// TestIntegration is the main integration test entry point that runs all provider functionality tests
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	t.Run("Datastore", func(t *testing.T) {
		t.Run("Validation", func(t *testing.T) {
			TestDatastoreValidation(t)
		})
	})

	t.Run("Remotes", func(t *testing.T) {
		t.Run("DataSources", func(t *testing.T) {
			TestRemoteDataSources(t)
		})
	})

	t.Run("Metrics", func(t *testing.T) {
		t.Run("MetricsServerVerifyCertificate", func(t *testing.T) {
			TestMetricsServerVerifyCertificate(t)
		})

		t.Run("MetricsServerMaxBodySize", func(t *testing.T) {
			TestMetricsServerMaxBodySize(t)
		})
	})
}

// TestQuickSmoke runs basic smoke tests that should always pass
func TestQuickSmoke(t *testing.T) {
	// Quick connectivity and basic function tests
	tc := SetupTest(t)
	defer tc.DestroyTerraform(t)

	// Test PBS server connectivity
	datastoreName := GenerateTestName("smoke-test")

	// Simple directory datastore creation as smoke test
	testConfig := `
resource "pbs_datastore" "smoke" {
  name = "` + datastoreName + `"
  path = "/datastore/` + datastoreName + `"
}
`

	tc.WriteMainTF(t, testConfig)
	tc.ApplyTerraform(t)

	// Verify resource exists in state
	resource := tc.GetResourceFromState(t, "pbs_datastore.smoke")
	if resource.AttributeValues["name"] != datastoreName {
		t.Errorf("Expected datastore name %s, got %v", datastoreName, resource.AttributeValues["name"])
	}

	debugLog(t, "Smoke test passed: PBS provider basic functionality works")
}
