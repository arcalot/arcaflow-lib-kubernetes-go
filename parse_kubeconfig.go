package arcaflow_lib_kubernetes

import (
	"arcaflow-lib-kubernetes/internal/util"
	"encoding/base64"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"os"
	"strings"
	"time"
)

func ParseKubeConfig(data string) (KubeConfig, error) {
	var kubeconfig = KubeConfig{}
	if err := yaml.Unmarshal([]byte(data), &kubeconfig); err != nil {
		return kubeconfig, err
	}
	return kubeconfig, nil
}

func KubeConfigToConnection(kubeconfig KubeConfig, inlineFiles bool) (ConnectionParameters, error) {
	if kubeconfig.CurrentContext == nil {
		return ConnectionParameters{}, errors.New("unusable KubeConfig: no current context is set")
	}
	var context *KubeConfigContext
	for _, ctx := range kubeconfig.Contexts {
		if ctx.Name == *kubeconfig.CurrentContext {
			context = &ctx
		}
	}
	if context == nil {
		return ConnectionParameters{}, fmt.Errorf("current context %s not found in kubeconfig file", *kubeconfig.CurrentContext)
	}
	currentCluster := context.Context.Cluster
	currentUser := context.Context.User

	var cluster *KubeConfigCluster

	for _, clstr := range kubeconfig.Clusters {
		if clstr.Name == currentCluster {
			cluster = &clstr
		}
	}
	if cluster == nil {
		return ConnectionParameters{}, fmt.Errorf("current cluster %s not found in kubeconfig file", currentCluster)
	}

	var user *KubeConfigUser
	for _, usr := range kubeconfig.Users {
		if usr.Name == currentUser {
			user = &usr
		}
	}

	if user == nil {
		return ConnectionParameters{}, fmt.Errorf("current user %s not found in kubeconfig file", currentUser)
	}

	if len(cluster.Cluster.Server) == 0 {
		return ConnectionParameters{}, errors.New("no cluster host found in connection")
	}

	connectionParams := ConnectionParameters{
		Host: strings.Replace(strings.Replace(cluster.Cluster.Server, "https://", "", 1), "http://", "", 1),
	}

	if cluster.Cluster.CertificateAuthority != nil {
		if inlineFiles {
			data, err := os.ReadFile(*cluster.Cluster.CertificateAuthority)
			if err != nil {
				return ConnectionParameters{}, err
			}
			connectionParams.CAData = string(data)
		} else {
			connectionParams.CAFile = *cluster.Cluster.CertificateAuthority
		}
	}

	connectionParams.Insecure = cluster.Cluster.InsecureSkipTLSVerify

	if cluster.Cluster.CertificateAuthorityData != nil {
		connectionParams.CAData = util.Base64Decode(*cluster.Cluster.CertificateAuthorityData)
	}

	if user.User.ClientCertificate != nil {
		if inlineFiles {
			data, err := os.ReadFile(*user.User.ClientCertificate)
			if err != nil {
				return ConnectionParameters{}, err
			}
			connectionParams.CertData = string(data)
		} else {
			connectionParams.CertFile = *user.User.ClientCertificate
		}
	}

	if user.User.ClientCertificateData != nil {
		connectionParams.CertData = util.Base64Decode(*user.User.ClientCertificateData)
	}

	if user.User.ClientKey != nil {
		if inlineFiles {
			data, err := os.ReadFile(*user.User.ClientKey)
			if err != nil {
				return ConnectionParameters{}, err
			}
			connectionParams.KeyData = string(data)
		} else {
			connectionParams.KeyFile = *user.User.ClientKey
		}
	}
	if user.User.ClientKeyData != nil {
		connectionParams.KeyData = util.Base64Decode(*user.User.ClientKeyData)
	}

	if user.User.Username != nil {
		connectionParams.Username = *user.User.Username
	}
	if user.User.Password != nil {
		connectionParams.Password = *user.User.Password
	}
	if user.User.Token != nil {
		connectionParams.BearerToken = *user.User.Token
	}

	if err := ConnectionParametersSchema().Validate(connectionParams); err != nil {
		return ConnectionParameters{}, err
	}
	return connectionParams, nil
}

func ConnectionToKubeConfig(connection ConnectionParameters) (KubeConfig, error) {
	defaultStr := "default"
	clusterParams := KubeConfigClusterParams{}
	if len(connection.Host) == 0 {
		return KubeConfig{}, errors.New("no cluster host found in connection")
	}
	clusterParams.Server = fmt.Sprintf("https://%s", connection.Host)
	if len(connection.CAData) > 0 {
		caData := base64.StdEncoding.EncodeToString([]byte(connection.CAData))
		clusterParams.CertificateAuthorityData = &caData
	}
	if len(connection.CAFile) > 0 {
		clusterParams.CertificateAuthority = &connection.CAFile
	}
	clusterParams.InsecureSkipTLSVerify = connection.Insecure
	cluster := KubeConfigCluster{
		Cluster: clusterParams,
		Name:    defaultStr,
	}

	contextParams := KubeConfigContextParameters{}
	contextParams.User = connection.Username
	contextParams.Cluster = defaultStr
	contextParams.Namespace = defaultStr
	context := KubeConfigContext{
		Context: contextParams,
		Name:    defaultStr,
	}

	userParams := KubeConfigUserParameters{}
	userParams.Username = &connection.Username
	userParams.Password = &connection.Password

	if len(connection.KeyData) > 0 {
		keyData := base64.StdEncoding.EncodeToString([]byte(connection.KeyData))
		userParams.ClientKeyData = &keyData
	}
	if len(connection.KeyFile) > 0 {
		userParams.ClientKey = &connection.KeyFile
	}
	if len(connection.CertData) > 0 {
		certData := base64.StdEncoding.EncodeToString([]byte(connection.CertData))
		userParams.ClientCertificateData = &certData
	}
	if len(connection.CertFile) > 0 {
		userParams.ClientCertificate = &connection.CertFile
	}

	if len(connection.BearerToken) > 0 {
		userParams.Token = &connection.BearerToken
	}

	if len(connection.BearerTokenFile) > 0 {
		tokenData, err := os.ReadFile(connection.BearerTokenFile)
		if err != nil {
			return KubeConfig{}, err
		}
		token := string(tokenData)
		userParams.Token = &token
	}

	user := KubeConfigUser{
		User: userParams,
		Name: connection.Username,
	}

	kubeconfig := KubeConfig{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       []KubeConfigCluster{cluster},
		Contexts:       []KubeConfigContext{context},
		Users:          []KubeConfigUser{user},
		CurrentContext: &defaultStr,
		Preferences:    map[interface{}]interface{}{}, // default type in sdk parser, used to compare the struct in unit tests
	}

	return kubeconfig, nil

}

func ConnectionToRestConfig(connection ConnectionParameters) (*restclient.Config, error) {
	const defaultTimeOut = 10 * time.Second

	if len(connection.Host) == 0 {
		return nil, errors.New("no cluster host found in connection")
	}

	clientConfig := restclient.Config{
		Host:    connection.Host,
		APIPath: connection.APIPath,
		ContentConfig: restclient.ContentConfig{
			GroupVersion:         &core.SchemeGroupVersion,
			NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
		},
		Username:    connection.Username,
		Password:    connection.Password,
		BearerToken: connection.BearerToken,
		Impersonate: restclient.ImpersonationConfig{},
		TLSClientConfig: restclient.TLSClientConfig{
			ServerName: connection.ServerName,
			CertData:   []byte(connection.CertData),
			CertFile:   connection.CertFile,
			KeyData:    []byte(connection.KeyData),
			KeyFile:    connection.KeyFile,
			CAData:     []byte(connection.CAData),
			CAFile:     connection.CAFile,
			Insecure:   connection.Insecure,
		},
		UserAgent: "Arcaflow",
		QPS:       restclient.DefaultQPS,
		Burst:     restclient.DefaultBurst,
		Timeout:   defaultTimeOut,
	}
	return &clientConfig, nil
}

func Client(connection ConnectionParameters) (*kubernetes.Clientset, error) {
	config, err := ConnectionToRestConfig(connection)
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

func RESTClient(connection ConnectionParameters) (*restclient.RESTClient, error) {
	clientConfig, err := ConnectionToRestConfig(connection)
	if err != nil {
		return nil, err
	}
	client, err := restclient.RESTClientFor(clientConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}
