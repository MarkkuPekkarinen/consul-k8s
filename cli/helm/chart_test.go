package helm

import (
	"embed"
	"fmt"
	"testing"

	"github.com/hashicorp/consul-k8s/cli/test/mock"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/client-go/kubernetes/fake"
)

// Embed a test chart to test against.
//go:embed test_fixtures/consul/* test_fixtures/consul/templates/_helpers.tpl
var testChartFiles embed.FS

func TestLoadChart(t *testing.T) {
	directory := "fixtures/consul"

	expectedApiVersion := "v2"
	expectedName := "Foo"
	expectedVersion := "0.1.0"
	expectedDescription := "Mock Helm Chart for testing."
	expectedValues := map[string]interface{}{
		"key": "value",
	}

	actual, err := LoadChart(testChartFiles, directory)
	require.NoError(t, err)
	require.Equal(t, expectedApiVersion, actual.Metadata.APIVersion)
	require.Equal(t, expectedName, actual.Metadata.Name)
	require.Equal(t, expectedVersion, actual.Metadata.Version)
	require.Equal(t, expectedDescription, actual.Metadata.Description)
	require.Equal(t, expectedValues, actual.Values)
}

func TestFetchChartValues(t *testing.T) {
	namespace := "default"
	name := "consul"
	settings := mock.CreateMockEnvSettings(t, namespace)
	logger := mock.CreateLogger(t)

	expected := map[string]interface{}{}

	actual, err := FetchChartValues(namespace, name, settings, logger)
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}

func TestA(t *testing.T) {
	k8s := fake.NewSimpleClientset()

	config := new(action.Configuration)
	config.KubeClient = mock.FakeClient{
		K8sClient: k8s,
	}
	d := driver.NewSecrets(k8s.CoreV1().Secrets("default"))
	d.Log = mock.CreateLogger(t)
	r := &release.Release{
		Name: "consul",
	}
	d.Create("consul", r)
	store := storage.Init(d)
	config.Releases = store

	actual, err := a(config, "consul")
	require.NoError(t, err)

	fmt.Println(actual)
}

func TestReadChartFiles(t *testing.T) {
	directory := "fixtures/consul"
	expectedFiles := map[string]string{
		"Chart.yaml":             "# This is a mock Helm Chart.yaml file used for testing.\napiVersion: v2\nname: Foo\nversion: 0.1.0\ndescription: Mock Helm Chart for testing.",
		"values.yaml":            "# This is a mock Helm values.yaml file used for testing.\nkey: value",
		"templates/_helpers.tpl": "helpers",
		"templates/foo.yaml":     "foo: bar\n",
	}

	files, err := readChartFiles(testChartFiles, directory)
	require.NoError(t, err)

	actualFiles := make(map[string]string, len(files))
	for _, f := range files {
		actualFiles[f.Name] = string(f.Data)
	}

	for expectedName, expectedContents := range expectedFiles {
		actualContents, ok := actualFiles[expectedName]
		require.True(t, ok, "Expected file %s not found", expectedName)
		require.Equal(t, expectedContents, actualContents)
	}
}
