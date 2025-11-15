package secret

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/bytedance/sonic"
)

type GSM struct {
	projectID string
}

func NewSecretGSM(projectID string) Interface {
	return &GSM{
		projectID: projectID,
	}
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
		Name: fmt.Sprintf("projects/%s/secrets/domea/versions/latest", v.projectID),
	})
	if err != nil {
		return err
	}
	byteCfg = resp.GetPayload().GetData()

	err = sonic.Unmarshal(byteCfg, &target)
	if err != nil {
		return err
	}

	return nil
}
