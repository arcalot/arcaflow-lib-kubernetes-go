package arcaflow_lib_kubernetes

type KubeConfig struct {
    Kind           string              `json:"kind"`
    APIVersion     string              `json:"apiVersion"`
    Clusters       []KubeConfigCluster `json:"clusters"`
    Contexts       []KubeConfigContext `json:"contexts"`
    Users          []KubeConfigUser    `json:"users"`
    CurrentContext *string             `json:"current-context"`
    Preferences    any                 `json:"preferences"`
}

type KubeConfigClusterParams struct {
    Server                   string  `json:"server"`
    CertificateAuthority     *string `json:"certificate-authority"`
    CertificateAuthorityData *string `json:"certificate-authority-data"`
    InsecureSkipTLSVerify    bool    `json:"insecure-skip-tls-verify"`
}

type KubeConfigCluster struct {
    Name    string                  `json:"name"`
    Cluster KubeConfigClusterParams `json:"cluster"`
}

type KubeConfigContextParameters struct {
    Cluster   string `json:"cluster"`
    User      string `json:"user"`
    Namespace string `json:"namespace"`
}

type KubeConfigContext struct {
    Name    string                      `json:"name"`
    Context KubeConfigContextParameters `json:"context"`
}

type KubeConfigUserParameters struct {
    Username              *string `json:"username"`
    Password              *string `json:"password"`
    Token                 *string `json:"token"`
    ClientCertificate     *string `json:"client-certificate"`
    ClientCertificateData *string `json:"client-certificate-data"`
    ClientKey             *string `json:"client-key"`
    ClientKeyData         *string `json:"client-key-data"`
}

type KubeConfigUser struct {
    Name string                   `json:"name"`
    User KubeConfigUserParameters `json:"user"`
}
