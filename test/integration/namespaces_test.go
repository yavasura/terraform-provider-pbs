package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pbsdatastores "github.com/yavasura/terraform-provider-pbs/pbs/datastores"
)

func TestNamespacesHierarchyIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := SetupTest(t)
	defer tc.DestroyTerraform(t)

	store := getenvDefault("PBS_TEST_NAMESPACE_STORE", "backup")
	parent := fmt.Sprintf("tf-ns-%s", GenerateTestName("parent"))
	child := parent + "/vms"

	datastoreClient := pbsdatastores.NewClient(tc.APIClient)

	defer func() {
		_ = datastoreClient.DeleteNamespace(context.Background(), store, child)
		_ = datastoreClient.DeleteNamespace(context.Background(), store, parent)
	}()

	config := fmt.Sprintf(`
resource "pbs_namespace" "parent" {
  store     = "%s"
  namespace = "%s"
}

resource "pbs_namespace" "child" {
  store      = "%s"
  namespace  = "%s"
  depends_on = [pbs_namespace.parent]
}
`, store, parent, store, child)

	tc.WriteMainTF(t, config)
	tc.ApplyTerraform(t)

	ns, err := datastoreClient.GetNamespace(context.Background(), store, child)
	require.NoError(t, err)
	assert.Equal(t, child, ns.Path)
}

func TestNamespaceImportIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := SetupTest(t)
	defer tc.DestroyTerraform(t)

	store := getenvDefault("PBS_TEST_NAMESPACE_STORE", "backup")
	namespace := fmt.Sprintf("tf-ns-%s", GenerateTestName("import"))
	datastoreClient := pbsdatastores.NewClient(tc.APIClient)

	err := datastoreClient.CreateNamespace(context.Background(), store, namespace, "")
	require.NoError(t, err)
	defer func() {
		_ = datastoreClient.DeleteNamespace(context.Background(), store, namespace)
	}()

	config := fmt.Sprintf(`
resource "pbs_namespace" "imported" {
  store     = "%s"
  namespace = "%s"
}
`, store, namespace)

	tc.WriteMainTF(t, config)
	tc.ImportResource(t, "pbs_namespace.imported", fmt.Sprintf("%s:%s", store, namespace))
	tc.ApplyTerraform(t)
}

func getenvDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
