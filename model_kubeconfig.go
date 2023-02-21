package arcaflow_lib_kubernetes

type KubeConfig struct {
	Kind           string              `json:"kind" yaml:"kind"`
	APIVersion     string              `json:"apiVersion" yaml:"apiVersion"`
	Clusters       []KubeConfigCluster `json:"clusters" yaml:"clusters"`
	Contexts       []KubeConfigContext `json:"contexts" yaml:"contexts"`
	Users          []KubeConfigUser    `json:"users" yaml:"users"`
	CurrentContext *string             `json:"current-context" yaml:"current-context"`
	Preferences    any                 `json:"preferences" yaml:"preferences"`
}

type KubeConfigClusterParams struct {
	Server                   string  `json:"server" yaml:"server"`
	CertificateAuthority     *string `json:"certificate-authority" yaml:"certificate-authority"`
	CertificateAuthorityData *string `json:"certificate-authority-data" yaml:"certificate-authority-data"`
	InsecureSkipTLSVerify    bool    `json:"insecure-skip-tls-verify" yaml:"insecure-skip-tls-verify"`
}

type KubeConfigCluster struct {
	Name    string                  `json:"name" yaml:"name"`
	Cluster KubeConfigClusterParams `json:"cluster" yaml:"cluster"`
}

type KubeConfigContextParameters struct {
	Cluster   string `json:"cluster" yaml:"cluster"`
	User      string `json:"user" yaml:"user"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

type KubeConfigContext struct {
	Name    string                      `json:"name" yaml:"name"`
	Context KubeConfigContextParameters `json:"context" yaml:"context"`
}

type KubeConfigUserParameters struct {
	Username              *string `json:"username" yaml:"username"`
	Password              *string `json:"password" yaml:"password"`
	Token                 *string `json:"token" yaml:"token"`
	ClientCertificate     *string `json:"client-certificate" yaml:"client-certificate"`
	ClientCertificateData *string `json:"client-certificate-data" yaml:"client-certificate-data"`
	ClientKey             *string `json:"client-key" yaml:"client-key"`
	ClientKeyData         *string `json:"client-key-data" yaml:"client-key-data"`
}

type KubeConfigUser struct {
	Name string                   `json:"name" yaml:"name"`
	User KubeConfigUserParameters `json:"user" yaml:"user"`
}
