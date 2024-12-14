package keyclient

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-acme/lego/log"
)

// Implements SecretManagerKeyClientInterface
type AWSSecretManagerKeyClient struct {
	secretClient *sm.Client
	projectID    string
	secretID     string
}

// New create a new SecretManagerKeyService.
func NewAwsClient(ctx context.Context, projectID, secretID string) (*AWSSecretManagerKeyClient, error) {

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
	)
	if err != nil {
		log.Fatal("Could not get config %v", err)
	}
	//resolver := localstack.NewSecretsManagerResolverV2(l)
	cli := sm.NewFromConfig(cfg, func(o *sm.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
	})
	client := &AWSSecretManagerKeyClient{
		secretClient: cli,
		projectID:    projectID,
		secretID:     secretID,
	}
	return client, nil
}

func (awssm *AWSSecretManagerKeyClient) ServiceSigningPrivateKeyset(ctx context.Context) ([]byte, error) {
	name := fmt.Sprintf("proejctid/%s/secretid/%s/%s", awssm.projectID, awssm.secretID, "signingKey")
	output, err := awssm.secretClient.GetSecretValue(ctx, &sm.GetSecretValueInput{
		SecretId: &name,
	})
	if err != nil {
		return nil, err
	}
	return output.SecretBinary, nil
}

func (awssm *AWSSecretManagerKeyClient) ServiceEncryptionPrivateKey(ctx context.Context) ([]byte, error) {
	name := fmt.Sprintf("proejctid/%s/secretid/%s/%s", awssm.projectID, awssm.secretID, "encryptionKey")
	output, err := awssm.secretClient.GetSecretValue(ctx, &sm.GetSecretValueInput{
		SecretId: &name,
	})
	if err != nil {
		return nil, err
	}
	return output.SecretBinary, nil
}

func (awssm *AWSSecretManagerKeyClient) AddKey(ctx context.Context, key string, payload []byte) error {
	name := fmt.Sprintf("proejctid/%s/secretid/%s/%s", awssm.projectID, awssm.secretID, key)
	_, err := awssm.secretClient.CreateSecret(ctx, &sm.CreateSecretInput{
		Name:         &name,
		SecretBinary: payload,
	})
	if err != nil {
		return err
	}
	return nil
}

func (awssm *AWSSecretManagerKeyClient) Close() {
	return
}
