package helm

import (
	"testing"

	"github.com/hashicorp/consul-k8s/cli/test/mock"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/action"
)

func TestInitActionConfig(t *testing.T) {
	actionConfig := &action.Configuration{}
	namespace := "consul"
	settings := mock.CreateMockEnvSettings(t, namespace)
	logger := mock.CreateLogger(t)

	expected := &action.Configuration{}

	actual, err := InitActionConfig(actionConfig, namespace, settings, logger)
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}
