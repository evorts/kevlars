/**
 * @Author: steven
 * @Description:
 * @File: provider_gsm
 * @Date: 29/09/23 10.49
 */

package config

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"
	"hash/crc32"
	"strings"
)

type googleSecret struct {
	projectId    string
	resourceName string
	configType   string
	data         map[string]interface{}

	jsonCredFile string
	jsonCred     []byte
}

func (g *googleSecret) GetData() map[string]interface{} {
	return g.data
}

// Init initialize google secret manager and populate into conf object
func (g *googleSecret) Init() error {
	ctx := context.Background()
	var opts = make([]option.ClientOption, 0)
	if len(g.jsonCredFile) > 0 {
		opts = append(opts, option.WithCredentialsFile(g.jsonCredFile))
	}
	if g.jsonCred != nil && len(g.jsonCred) > 0 {
		opts = append(opts, option.WithCredentialsJSON(g.jsonCred))
	}
	client, err := secretmanager.NewClient(ctx, opts...)
	if err != nil {
		return err
	}
	defer func(client *secretmanager.Client) {
		err = client.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(client)
	secretPath := fmt.Sprintf("%s/versions/%s", g.projectId, g.resourceName)
	if strings.HasSuffix(g.projectId, "/versions") {
		secretPath = fmt.Sprintf("%s/%s", g.projectId, g.resourceName)
	}
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath,
	}
	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to access secret version: %v", err)
	}
	// Verify the data checksum.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		return fmt.Errorf("data corruption detected")
	}
	if g.configType == "json" {
		return json.Unmarshal(result.Payload.Data, &g.data)
	}
	return yaml.Unmarshal(result.Payload.Data, &g.data)
}

// NewGoogleSecretManager instantiate google secret manager with the given projectId and secret name
// The resource name of the [Secret][google.cloud.secretmanager.v1.Secret], in the format `projects/*/secrets/*`.
func NewGoogleSecretManager(projectId, resourceName, configType string, opts ...GsmOption) Provider {
	p := &googleSecret{
		projectId:    projectId,
		resourceName: resourceName,
		configType:   configType,
	}
	for _, opt := range opts {
		opt.apply(p)
	}
	return p
}
