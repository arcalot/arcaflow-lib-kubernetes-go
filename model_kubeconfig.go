package arcaflow_lib_kubernetes

import (
	"encoding/json"
	"fmt"
	"go.flow.arcalot.io/pluginsdk/schema"
)

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
	Extensions               any     `json:"extensions"`
}

type KubeConfigCluster struct {
	Name    string                  `json:"name"`
	Cluster KubeConfigClusterParams `json:"cluster"`
}

type KubeConfigContextParameters struct {
	Cluster    string `json:"cluster"`
	User       string `json:"user"`
	Namespace  string `json:"namespace"`
	Extensions any    `json:"extensions"`
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

func (k *KubeConfig) UnmarshalJSON(data []byte) error {
	return k.UnmarshalYAML(func(a any) error {
		return json.Unmarshal(data, a)
	})
}

// MarshalJSON uses the Arcaflow schema system to marshal JSON data when called via json.Marshal on the
// ConnectionParameters struct. This prevents accidentally using the wrong unmarshalling method.
func (k *KubeConfig) MarshalJSON() ([]byte, error) {
	serializedData, err := k.MarshalYAML()
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(serializedData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data (%w)", err)
	}
	return data, nil
}

// UnmarshalYAML uses the Arcaflow schema system to unmarshal YAML data when called via yaml.Unmarshal on the
// ConnectionParameters struct. This prevents accidentally using the wrong unmarshalling method.
func (k *KubeConfig) UnmarshalYAML(unmarshaller func(interface{}) error) error {
	temp := map[string]any{}
	if err := unmarshaller(&temp); err != nil {
		return fmt.Errorf("failed to JSON unmarshal data (%w)", err)
	}
	unserializedData, err := kubeconfigSchema.UnserializeType(temp)
	if err != nil {
		return fmt.Errorf("failed to unserialize data (%w)", err)
	}
	k.APIVersion = unserializedData.APIVersion
	k.Kind = unserializedData.Kind
	k.CurrentContext = unserializedData.CurrentContext
	k.Preferences = unserializedData.Preferences
	k.Users = unserializedData.Users
	k.Clusters = unserializedData.Clusters
	k.Contexts = unserializedData.Contexts

	return nil
}

// MarshalYAML uses the Arcaflow schema system to marshal JSON data when called via json.Marshal on the
// ConnectionParameters struct. This prevents accidentally using the wrong unmarshalling method.
func (k *KubeConfig) MarshalYAML() (any, error) {
	serializedData, err := connectionParametersSchema.Serialize(k)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize connection parameters (%w)", err)
	}
	return serializedData, nil
}

var kubeconfigSchema = schema.NewTypedObject[KubeConfig](
	"KubeConfig",
	map[string]*schema.PropertySchema{
		"apiVersion": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("APIVersion"),
				schema.PointerTo("API Version"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"v1"`),
			nil,
		).TreatEmptyAsDefaultValue(),
		"kind": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Kind"),
				schema.PointerTo("kubernetes Resource Kind"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"Config"`),
			nil,
		).TreatEmptyAsDefaultValue(),
		"current-context": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("CurrentContext"),
				schema.PointerTo("active context"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"Config"`),
			nil,
		).TreatEmptyAsDefaultValue(),
		"clusters": schema.NewPropertySchema(
			schema.NewListSchema(clusterSchema, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Clusters"),
				schema.PointerTo("kubeconfig clusters"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"Clusters"`),
			nil,
		),
		"contexts": schema.NewPropertySchema(
			schema.NewListSchema(contextSchema, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Contexts"),
				schema.PointerTo("kubeconfig contexts"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"Contexts"`),
			nil,
		),
		"users": schema.NewPropertySchema(
			schema.NewListSchema(userSchema, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Users"),
				schema.PointerTo("kubeconfig users"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"Users"`),
			nil,
		),
		"preferences": schema.NewPropertySchema(
			schema.NewAnySchema(),
			schema.NewDisplayValue(
				schema.PointerTo("Preferences"),
				schema.PointerTo("Kubeconfig preferences"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"Preferences"`),
			nil,
		),
	},
)

var clusterSchema = schema.NewTypedObject[KubeConfigCluster](
	"KubeConfigCluster",
	map[string]*schema.PropertySchema{
		"name": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Name"),
				schema.PointerTo("cluster name"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"cluster": schema.NewPropertySchema(
			clusterParamsSchema,
			schema.NewDisplayValue(
				schema.PointerTo("Cluster"),
				schema.PointerTo("cluster parameters"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		),
	},
)

var clusterParamsSchema = schema.NewTypedObject[KubeConfigClusterParams](
	"KubeConfigClusterParameters",
	map[string]*schema.PropertySchema{
		"server": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Server"),
				schema.PointerTo("host name and port of the kubernetes server"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"certificate-authority": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("CertificateAuthority"),
				schema.PointerTo("cluster CA certificate path"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"certificate-authority-data": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("CertificateAuthorityData"),
				schema.PointerTo("cluster CA certificate base64 encoded"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"insecure-skip-tls-verify": schema.NewPropertySchema(
			schema.NewBoolSchema(),
			schema.NewDisplayValue(
				schema.PointerTo("InsecureSkipTLSVerify"),
				schema.PointerTo("toggles TLS verification"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"extensions": schema.NewPropertySchema(
			schema.NewAnySchema(),
			schema.NewDisplayValue(
				schema.PointerTo("Extensions"),
				schema.PointerTo("minikube kube config section "+
					"introduced to avoid local tests issues"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		),
	},
)

var contextSchema = schema.NewTypedObject[KubeConfigContext](
	"KubeConfigContext",
	map[string]*schema.PropertySchema{
		"name": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Name"),
				schema.PointerTo("context name"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"context": schema.NewPropertySchema(
			contextParamsSchema,
			schema.NewDisplayValue(
				schema.PointerTo("Context"),
				schema.PointerTo("context parameters object"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
	},
)
var contextParamsSchema = schema.NewTypedObject[KubeConfigContextParameters](
	"KubeConfigContextParameters",
	map[string]*schema.PropertySchema{
		"cluster": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Cluster"),
				schema.PointerTo("cluster name of the context"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"user": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("User"),
				schema.PointerTo("user name of the context"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"namespace": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Namespace"),
				schema.PointerTo("default namespace of the context"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"extensions": schema.NewPropertySchema(
			schema.NewAnySchema(),
			schema.NewDisplayValue(
				schema.PointerTo("Extensions"),
				schema.PointerTo("minikube kube config section "+
					"introduced to avoid local tests issues"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		),
	},
)

var userSchema = schema.NewTypedObject[KubeConfigUser](
	"KubeConfigUser",
	map[string]*schema.PropertySchema{
		"name": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Name"),
				schema.PointerTo("user name"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"user": schema.NewPropertySchema(
			userParamsSchema,
			schema.NewDisplayValue(
				schema.PointerTo("User"),
				schema.PointerTo("user parameters object"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
	},
)

var userParamsSchema = schema.NewTypedObject[KubeConfigUserParameters](
	"KubeConfigUserParameters",
	map[string]*schema.PropertySchema{
		"username": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Username"),
				schema.PointerTo("user username for basic authentication"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"password": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Password"),
				schema.PointerTo("user password for basic authentication"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"token": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Token"),
				schema.PointerTo("user bearer token"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"client-certificate": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("ClientCertificate"),
				schema.PointerTo("client certificate path"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"client-certificate-data": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("ClientCertificateData"),
				schema.PointerTo("client certificate data base64 encoded"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"client-key": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("ClientKey"),
				schema.PointerTo("client private key path"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"client-key-data": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("ClientKeyData"),
				schema.PointerTo("client private key data base64 encoded"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
	},
)

func KubeconfigSchema() *schema.ObjectSchema {
	return &kubeconfigSchema.ObjectSchema
}
