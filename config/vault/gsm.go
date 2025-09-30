package vault

import (
	"context"
	"encoding/json"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type GSM struct {
	ProjectID string
}

func (v *GSM) GetSecret(ctx context.Context, target interface{}) error {
	var (
		byteCfg []byte
	)

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	resp, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/domea/versions/latest", v.ProjectID),
	})
	if err != nil {
		return err
	}
	byteCfg = resp.GetPayload().GetData()

	err = json.Unmarshal(byteCfg, &target)
	if err != nil {
		return err
	}

	return nil
}
