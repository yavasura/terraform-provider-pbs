package namespaces

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pbsdatastores "github.com/yavasura/terraform-provider-pbs/pbs/datastores"
)

func TestNamespacePathValidation(t *testing.T) {
	assert.True(t, pbsdatastores.IsValidNamespacePath("production"))
	assert.True(t, pbsdatastores.IsValidNamespacePath("production/vms"))
	assert.True(t, pbsdatastores.IsValidNamespacePath("team_a/project-1/releases"))

	assert.False(t, pbsdatastores.IsValidNamespacePath(""))
	assert.False(t, pbsdatastores.IsValidNamespacePath("/production"))
	assert.False(t, pbsdatastores.IsValidNamespacePath("production/"))
	assert.False(t, pbsdatastores.IsValidNamespacePath("bad space"))
	assert.False(t, pbsdatastores.IsValidNamespacePath("a/b/c/d/e/f/g/h"))
}

func TestNamespaceResourceSchema(t *testing.T) {
	r := &namespaceResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	require.False(t, resp.Diagnostics.HasError())
	require.Contains(t, resp.Schema.Attributes, "store")
	require.Contains(t, resp.Schema.Attributes, "namespace")
	require.Contains(t, resp.Schema.Attributes, "comment")
}
