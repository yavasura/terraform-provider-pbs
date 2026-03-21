package namespaces

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/require"
)

func TestNamespacesDataSourceSchema(t *testing.T) {
	ds := &namespacesDataSource{}
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	require.False(t, resp.Diagnostics.HasError())
	require.Contains(t, resp.Schema.Attributes, "store")
	require.Contains(t, resp.Schema.Attributes, "prefix")
	require.Contains(t, resp.Schema.Attributes, "max_depth")
	require.Contains(t, resp.Schema.Attributes, "namespaces")
}
