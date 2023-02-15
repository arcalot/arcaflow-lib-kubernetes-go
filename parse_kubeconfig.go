package arcaflow_lib_kubernetes

import (
    "k8s.io/client-go/kubernetes"
    restclient "k8s.io/client-go/rest"
)

func ParseKubeConfig(data string) (KubeConfig, error) {

}

func KubeConfigToConnection(kubeconfig KubeConfig) (ConnectionParameters, error) {

}

func ConnectionToKubeConfig(connection ConnectionParameters) (KubeConfig, error) {

}

func Client(connection ConnectionParameters) (*kubernetes.Clientset, error) {
}

func RESTClient(connection ConnectionParameters) (*restclient.RESTClient, error) {

}
