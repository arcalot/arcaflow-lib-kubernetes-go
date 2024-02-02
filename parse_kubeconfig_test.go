package arcaflow_lib_kubernetes

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	CACERTPATH = "testdata/ca.crt"
	CERTPATH   = "testdata/client.crt"
	KEYPATH    = "testdata/client.key"
)

type testFixtures struct {
	caCert               string
	clientCrt            string
	clientKey            string
	kubeconfig           string
	kubeconfigNoData     string
	kubeconfigNoHost     string
	kubeconfigNoContext  string
	kubeconfigSkipTLS    string
	tokenFile            string
	kubeconfigExtensions string
}

func NewFixtures(t *testing.T) testFixtures {

	caCrt, err := os.ReadFile(CACERTPATH)
	assert.Nil(t, err)
	clientCrt, err := os.ReadFile(CERTPATH)
	assert.Nil(t, err)
	clientKey, err := os.ReadFile(KEYPATH)
	assert.Nil(t, err)
	kubeNodata, err := os.ReadFile("testdata/kubeconfig-nodata.yaml")
	assert.Nil(t, err)
	kube, err := os.ReadFile("testdata/kubeconfig-data.yaml")
	assert.Nil(t, err)
	kubeNoHost, err := os.ReadFile("testdata/kubeconfig-nohost.yaml")
	assert.Nil(t, err)
	kubeNoCtx, err := os.ReadFile("testdata/kubeconfig-nocontext.yaml")
	assert.Nil(t, err)
	kubeTLSSkip, err := os.ReadFile("testdata/kubeconfig-tlsskip.yaml")
	assert.Nil(t, err)
	tokenFile, err := os.ReadFile("testdata/tokenfile")
	assert.Nil(t, err)
	kubeExtensions, err := os.ReadFile("testdata/kubeconfig-extensions.yaml")
	assert.Nil(t, err)

	return testFixtures{
		caCert:               string(caCrt),
		clientCrt:            string(clientCrt),
		clientKey:            string(clientKey),
		kubeconfig:           string(kube),
		kubeconfigNoData:     string(kubeNodata),
		kubeconfigNoHost:     string(kubeNoHost),
		kubeconfigNoContext:  string(kubeNoCtx),
		kubeconfigSkipTLS:    string(kubeTLSSkip),
		tokenFile:            string(tokenFile),
		kubeconfigExtensions: string(kubeExtensions),
	}
}

func TestParseKubeConfig(t *testing.T) {
	fixtures := NewFixtures(t)
	kubeconf, err := ParseKubeConfig(fixtures.kubeconfig)
	assert.Nil(t, err)
	assert.NotNil(t, kubeconf)

	kubeconfExtensions, err := ParseKubeConfig(fixtures.kubeconfigExtensions)
	assert.Nil(t, err)
	assert.NotNil(t, kubeconfExtensions)
}

func TestKubeConfigToConnection(t *testing.T) {
	// test with cert inlining
	fixtures := NewFixtures(t)
	kubeconf, err := ParseKubeConfig(fixtures.kubeconfigNoData)
	assert.Nil(t, err)
	connection, err := KubeConfigToConnection(kubeconf, true)
	assert.Nil(t, err)
	assert.NotNil(t, connection)
	assert.Equal(t, connection.KeyData, fixtures.clientKey)
	assert.Equal(t, connection.CAData, fixtures.caCert)
	assert.Equal(t, connection.CertData, fixtures.clientCrt)
	// test that by default insecure-skip-tls-verify is false
	assert.False(t, connection.Insecure)

	// test without inlining
	kubeconf, err = ParseKubeConfig(fixtures.kubeconfigNoData)
	assert.Nil(t, err)
	connection, err = KubeConfigToConnection(kubeconf, false)
	assert.Nil(t, err)
	assert.Equal(t, connection.KeyFile, KEYPATH)
	assert.Equal(t, connection.CertFile, CERTPATH)
	assert.Equal(t, connection.CAFile, CACERTPATH)

	// test failure on empty host
	kubeconf, err = ParseKubeConfig(fixtures.kubeconfigNoHost)
	assert.NoError(t, err)
	assert.NotNil(t, kubeconf)
	connection, err = KubeConfigToConnection(kubeconf, true)
	assert.NotNil(t, err)

	// test failure on empty default context
	kubeconf, err = ParseKubeConfig(fixtures.kubeconfigNoContext)
	assert.NoError(t, err)
	assert.NotNil(t, kubeconf)
	connection, err = KubeConfigToConnection(kubeconf, true)
	assert.NotNil(t, err)

	// test success on insecure-skip-tls-verify: true
	kubeconf, err = ParseKubeConfig(fixtures.kubeconfigSkipTLS)
	assert.NoError(t, err)
	assert.NotNil(t, kubeconf)
	connection, err = KubeConfigToConnection(kubeconf, true)
	assert.True(t, connection.Insecure)
	assert.Nil(t, err)
}

func TestConnectionToKubeConfig(t *testing.T) {
	// test parsing without file inlining
	fixtures := NewFixtures(t)
	kubeconf, err := ParseKubeConfig(fixtures.kubeconfigNoData)
	assert.Nil(t, err)
	connection, err := KubeConfigToConnection(kubeconf, false)
	assert.NoError(t, err)
	kubeconfBack, err := ConnectionToKubeConfig(connection)
	assert.Nil(t, err)
	assert.Equal(t, kubeconf, kubeconfBack)
	// test parsing with file inlining
	kubeconf, err = ParseKubeConfig(fixtures.kubeconfig)
	assert.Nil(t, err)
	connection, err = KubeConfigToConnection(kubeconf, false)
	assert.NoError(t, err)
	kubeconfBack, err = ConnectionToKubeConfig(connection)
	assert.Nil(t, err)
	assert.Equal(t, kubeconf, kubeconfBack)

	// test parsing without inlining with token file in connection
	kubeconf, err = ParseKubeConfig(fixtures.kubeconfigNoData)
	assert.Nil(t, err)
	connection, err = KubeConfigToConnection(kubeconf, false)
	assert.NoError(t, err)
	connection.BearerToken = ""
	connection.BearerTokenFile = "testdata/tokenfile"
	kubeconfBack, err = ConnectionToKubeConfig(connection)
	assert.Nil(t, err)
	assert.Equal(t, kubeconf, kubeconfBack)

}
