package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yavasura/terraform-provider-pbs/pbs/metrics"
)

// getInfluxDBHost returns the InfluxDB host from environment or default
func getInfluxDBHost() string {
	if host := os.Getenv("TEST_INFLUXDB_HOST"); host != "" {
		return host
	}
	return "localhost"
}

// getInfluxDBPort returns the InfluxDB port from environment or default
func getInfluxDBPort() int {
	if port := os.Getenv("TEST_INFLUXDB_PORT"); port != "" {
		// Parse port from string, default to 8086 if parse fails
		var portNum int
		if _, err := fmt.Sscanf(port, "%d", &portNum); err == nil {
			return portNum
		}
	}
	return 8086
}

// getInfluxDBUDPHost returns the InfluxDB UDP host from environment or default
func getInfluxDBUDPHost() string {
	if host := os.Getenv("TEST_INFLUXDB_UDP_HOST"); host != "" {
		return host
	}
	return "localhost"
}

// getInfluxDBUDPPort returns the InfluxDB UDP port from environment or default
func getInfluxDBUDPPort() int {
	if port := os.Getenv("TEST_INFLUXDB_UDP_PORT"); port != "" {
		var portNum int
		if _, err := fmt.Sscanf(port, "%d", &portNum); err == nil {
			return portNum
		}
	}
	return 8089
}

// TestMetricsServerVerifyCertificate tests metrics server with certificate verification
func TestMetricsServerVerifyCertificate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := SetupTest(t)
	defer tc.DestroyTerraform(t)

	serverName := GenerateTestName("influxdb-tls")
	influxHost := getInfluxDBHost()
	influxPort := getInfluxDBPort()

	testConfig := fmt.Sprintf(`
resource "pbs_metrics_server" "test_tls" {
  name         = "%s"
  type         = "influxdb-http"
  url          = "http://%s:%d"
  organization = "testorg"
  bucket       = "pbs-metrics"
  token        = "test-token-123456"
  verify_tls   = false
  enable       = true
  comment      = "Metrics server with TLS verification disabled"
}
`, serverName, influxHost, influxPort)

	tc.WriteMainTF(t, testConfig)
	tc.ApplyTerraform(t)

	resource := tc.GetResourceFromState(t, "pbs_metrics_server.test_tls")
	assert.Equal(t, false, resource.AttributeValues["verify_tls"])

	metricsClient := metrics.NewClient(tc.APIClient)
	server, err := metricsClient.GetMetricsServer(context.Background(), metrics.MetricsServerTypeInfluxDBHTTP, serverName)
	require.NoError(t, err)
	assert.NotNil(t, server.VerifyTLS)
	assert.False(t, *server.VerifyTLS)
}

// TestMetricsServerMaxBodySize tests metrics server with custom max body size
func TestMetricsServerMaxBodySize(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tc := SetupTest(t)
	defer tc.DestroyTerraform(t)

	serverName := GenerateTestName("influxdb-bodysize")
	influxHost := getInfluxDBHost()
	influxPort := getInfluxDBPort()

	testConfig := fmt.Sprintf(`
resource "pbs_metrics_server" "test_bodysize" {
  name          = "%s"
  type          = "influxdb-http"
  url           = "http://%s:%d"
  organization  = "testorg"
  bucket        = "pbs-metrics"
  token         = "test-token-123456"
  max_body_size = 65536
  enable        = true
  comment       = "Metrics server with custom max body size"
}
`, serverName, influxHost, influxPort)

	tc.WriteMainTF(t, testConfig)
	tc.ApplyTerraform(t)

	resource := tc.GetResourceFromState(t, "pbs_metrics_server.test_bodysize")
	assert.Equal(t, json.Number("65536"), resource.AttributeValues["max_body_size"])

	metricsClient := metrics.NewClient(tc.APIClient)
	server, err := metricsClient.GetMetricsServer(context.Background(), metrics.MetricsServerTypeInfluxDBHTTP, serverName)
	require.NoError(t, err)
	assert.NotNil(t, server.MaxBodySize)
	assert.Equal(t, 65536, *server.MaxBodySize)
}

// TestMetricsServerTimeout tests metrics server with custom timeout
// NOTE: PBS 4.0 removed the timeout parameter, so this test is disabled
