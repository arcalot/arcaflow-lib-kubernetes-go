package arcaflow_lib_kubernetes

import (
	"encoding/json"
	"fmt"
	"regexp"

	"arcaflow-lib-kubernetes/internal/util"
	"go.flow.arcalot.io/pluginsdk/schema"
)

// ConnectionParameters describes how to connect to the Kubernetes API.
type ConnectionParameters struct {
	Host    string `json:"host"`
	APIPath string `json:"path"`

	Username string `json:"username"`
	Password string `json:"password"`

	ServerName string `json:"serverName"`

	CertData string `json:"cert"`
	CertFile string `json:"certFile"`
	KeyData  string `json:"key"`
	KeyFile  string `json:"keyFile"`
	CAData   string `json:"cacert"`
	CAFile   string `json:"cacertFile"`

	BearerToken     string `json:"bearerToken"`
	BearerTokenFile string `json:"bearerTokenFile"`
	Insecure        bool   `json:"insecure"`
}

// UnmarshalJSON uses the Arcaflow schema system to unmarshal JSON data when called via json.Unmarshal on the
// ConnectionParameters struct. This prevents accidentally using the wrong unmarshalling method.
func (c *ConnectionParameters) UnmarshalJSON(data []byte) error {
	return c.UnmarshalYAML(func(a any) error {
		return json.Unmarshal(data, a)
	})
}

// MarshalJSON uses the Arcaflow schema system to marshal JSON data when called via json.Marshal on the
// ConnectionParameters struct. This prevents accidentally using the wrong unmarshalling method.
func (c *ConnectionParameters) MarshalJSON() ([]byte, error) {
	serializedData, err := c.MarshalYAML()
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
func (c *ConnectionParameters) UnmarshalYAML(unmarshaller func(interface{}) error) error {
	temp := map[string]any{}
	if err := unmarshaller(&temp); err != nil {
		return fmt.Errorf("failed to JSON unmarshal data (%w)", err)
	}
	unserializedData, err := connectionParametersSchema.UnserializeType(temp) //nolint:all
	if err != nil {
		return fmt.Errorf("failed to unserialize data (%w)", err)
	}
	// This assignment has no effect!  We disable the linter's objection to it
	// pending investigation.
	// see https://github.com/arcalot/arcaflow-lib-kubernetes-go/issues/18
	c = &unserializedData //nolint:all
	return nil
}

// MarshalYAML uses the Arcaflow schema system to marshal JSON data when called via json.Marshal on the
// ConnectionParameters struct. This prevents accidentally using the wrong unmarshalling method.
func (c *ConnectionParameters) MarshalYAML() (any, error) {
	serializedData, err := connectionParametersSchema.Serialize(c)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize connection parameters (%w)", err)
	}
	return serializedData, nil
}

var connectionParametersSchema = schema.NewTypedObject[ConnectionParameters](
	"Connection",
	map[string]*schema.PropertySchema{
		"host": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Host"),
				schema.PointerTo("Host name and port of the Kubernetes server"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"kubernetes.default.svc"`),
			nil,
		).TreatEmptyAsDefaultValue(),
		"path": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Path"),
				schema.PointerTo("Path to the API server."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"/api"`),
			nil,
		).TreatEmptyAsDefaultValue(),
		"username": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Username"),
				schema.PointerTo("Username for basic authentication."),
				nil,
			),
			false,
			[]string{"password"},
			nil,
			nil,
			nil,
			nil,
		),
		"password": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Password"),
				schema.PointerTo("Password for basic authentication."),
				nil,
			),
			false,
			[]string{"username"},
			nil,
			nil,
			nil,
			nil,
		),
		"serverName": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("TLS server name"),
				schema.PointerTo("Expected TLS server name to verify in the certificate."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		),
		"cacert": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, regexp.MustCompile(`^$|^\s*-----BEGIN CERTIFICATE-----(\s*.*\s*)*-----END CERTIFICATE-----\s*$`)),
			schema.NewDisplayValue(
				schema.PointerTo("CA certificate"),
				schema.PointerTo("CA certificate in PEM format to verify Kubernetes server certificate against."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			[]string{
				util.JSONEncode(util.Base64Decode(`LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUI0VENDQVl1Z0F3SUJBZ0lVQ0hoaGZmWTFsemV6R2F0WU1SMDJncEVKQ2hrd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1JURUxNQWtHQTFVRUJoTUNRVlV4RXpBUkJnTlZCQWdNQ2xOdmJXVXRVM1JoZEdVeElUQWZCZ05WQkFvTQpHRWx1ZEdWeWJtVjBJRmRwWkdkcGRITWdVSFI1SUV4MFpEQWVGdzB5TWpBNU1qZ3dOVEk0TVRKYUZ3MHlNekE1Ck1qZ3dOVEk0TVRKYU1FVXhDekFKQmdOVkJBWVRBa0ZWTVJNd0VRWURWUVFJREFwVGIyMWxMVk4wWVhSbE1TRXcKSHdZRFZRUUtEQmhKYm5SbGNtNWxkQ0JYYVdSbmFYUnpJRkIwZVNCTWRHUXdYREFOQmdrcWhraUc5dzBCQVFFRgpBQU5MQURCSUFrRUFycjg5ZjJrZ2dTTy95YUNCNkV3SVFlVDZacHRCb1gwWnZDTUkrRHBrQ3dxT1M1ZndSYmoxCm5FaVBuTGJ6RERnTVU4S0NQQU1oSTdKcFlSbEhuaXB4V3dJREFRQUJvMU13VVRBZEJnTlZIUTRFRmdRVWlaNkoKRHd1RjlRQ2gxdndRR1hzMk11dHVROUV3SHdZRFZSMGpCQmd3Rm9BVWlaNkpEd3VGOVFDaDF2d1FHWHMyTXV0dQpROUV3RHdZRFZSMFRBUUgvQkFVd0F3RUIvekFOQmdrcWhraUc5dzBCQVFzRkFBTkJBRllJRk0yN0JEaUc3MjVkClZraFJibGt2WnplUkhoY3d0RE9RVEM5ZDhNL0x5bU4yeTBuSFNsSkNabS9Mby9hSDh2aVNZMXZpMUdTSGZEejcKVGxmZThncz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=`)),
			},
		),
		"cacertFile": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("CA certificate file"),
				schema.PointerTo("File holding the CA certificate in PEM format to verify Kubernetes server certificate against."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(util.JSONEncode("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")),
			nil,
		),
		"cert": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, regexp.MustCompile(`^$|^\s*-----BEGIN CERTIFICATE-----(\s*.*\s*)*-----END CERTIFICATE-----\s*$`)),
			schema.NewDisplayValue(
				schema.PointerTo("Client certificate"),
				schema.PointerTo("Client certificate in PEM format to authenticate against Kubernetes with."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			[]string{
				util.JSONEncode(util.Base64Decode(`LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUI0VENDQVl1Z0F3SUJBZ0lVQ0hoaGZmWTFsemV6R2F0WU1SMDJncEVKQ2hrd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1JURUxNQWtHQTFVRUJoTUNRVlV4RXpBUkJnTlZCQWdNQ2xOdmJXVXRVM1JoZEdVeElUQWZCZ05WQkFvTQpHRWx1ZEdWeWJtVjBJRmRwWkdkcGRITWdVSFI1SUV4MFpEQWVGdzB5TWpBNU1qZ3dOVEk0TVRKYUZ3MHlNekE1Ck1qZ3dOVEk0TVRKYU1FVXhDekFKQmdOVkJBWVRBa0ZWTVJNd0VRWURWUVFJREFwVGIyMWxMVk4wWVhSbE1TRXcKSHdZRFZRUUtEQmhKYm5SbGNtNWxkQ0JYYVdSbmFYUnpJRkIwZVNCTWRHUXdYREFOQmdrcWhraUc5dzBCQVFFRgpBQU5MQURCSUFrRUFycjg5ZjJrZ2dTTy95YUNCNkV3SVFlVDZacHRCb1gwWnZDTUkrRHBrQ3dxT1M1ZndSYmoxCm5FaVBuTGJ6RERnTVU4S0NQQU1oSTdKcFlSbEhuaXB4V3dJREFRQUJvMU13VVRBZEJnTlZIUTRFRmdRVWlaNkoKRHd1RjlRQ2gxdndRR1hzMk11dHVROUV3SHdZRFZSMGpCQmd3Rm9BVWlaNkpEd3VGOVFDaDF2d1FHWHMyTXV0dQpROUV3RHdZRFZSMFRBUUgvQkFVd0F3RUIvekFOQmdrcWhraUc5dzBCQVFzRkFBTkJBRllJRk0yN0JEaUc3MjVkClZraFJibGt2WnplUkhoY3d0RE9RVEM5ZDhNL0x5bU4yeTBuSFNsSkNabS9Mby9hSDh2aVNZMXZpMUdTSGZEejcKVGxmZThncz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=`)),
			},
		),
		"certFile": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Client certificate file"),
				schema.PointerTo("File holding the client certificate in PEM format to authenticate against Kubernetes with."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		),
		"key": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, regexp.MustCompile(`^$|^[-]+BEGIN (?:.* )?PRIVATE KEY[-]+([^-]*)[-]+END (?:.* )?PRIVATE KEY[-]+\s*$`)),
			schema.NewDisplayValue(
				schema.PointerTo("Client key"),
				schema.PointerTo("Client private key in PEM format to authenticate against Kubernetes with."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			[]string{
				util.JSONEncode(util.Base64Decode(`LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUJWQUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQVQ0d2dnRTZBZ0VBQWtFQXJyODlmMmtnZ1NPL3lhQ0IKNkV3SVFlVDZacHRCb1gwWnZDTUkrRHBrQ3dxT1M1ZndSYmoxbkVpUG5MYnpERGdNVThLQ1BBTWhJN0pwWVJsSApuaXB4V3dJREFRQUJBa0J5YnUveDBNRWxjR2kydS9KMlVkd1Njc1Y3amU1VHQxMno4Mmw3VEptWkZGSjhSTG1jCnJoMDBHdmViNFZwR2hkMStjM2xaYk8xbUlUNnYzdkhNOUEwaEFpRUExNEVXNmIrOTlYWXphNys1dXdJRHVpTSsKQnozcGtLKzl0bGZWWEU3SnlLc0NJUURQbFlKNXh0YnVUK1Z2QjNYT2REL1ZXaUVxRW12RTNmbFYwNDE3UnFoYQpFUUlnYnl4d05wd3RFZ0V0Vzh1bnRCckE4M2lVMmtXTlJZL3o3YXA0TGt1Uyswc0NJR2UyRSswUm1mcVFzbGxwCmljTXZNMkU5MllueWtDTlluNlR3d0NRU0pqUnhBaUVBbzlNbWFWbEs3WWRoU01QbzUydUpZemQ5TVFaSnFocSsKbEIxWkdEeC9BUkU9Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K`)),
			},
		),
		"keyFile": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Client key file"),
				schema.PointerTo("File holding the client private key in PEM format to authenticate against Kubernetes with."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		),
		"bearerToken": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Bearer token"),
				schema.PointerTo("Bearer token to authenticate against the Kubernetes API with."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		),
		"bearerTokenFile": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Bearer token file"),
				schema.PointerTo("File holding the bearer token to authenticate against the Kubernetes API with."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(util.JSONEncode("/var/run/secrets/kubernetes.io/serviceaccount")),
			nil,
		),
		"insecure": schema.NewPropertySchema(
			schema.NewBoolSchema(),
			schema.NewDisplayValue(
				schema.PointerTo("Insecure connection"),
				schema.PointerTo("Skip TLS verification"),
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

// ConnectionParametersSchema returns a schema for serializing/unserializing connection parameters.
func ConnectionParametersSchema() *schema.ObjectSchema {
	return &connectionParametersSchema.ObjectSchema
}
