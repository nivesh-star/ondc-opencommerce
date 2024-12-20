package localstackclient

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func NewSNSClient(ctx context.Context) (*sns.Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "test")),
	)
	if err != nil {
		panic(err)
	}

	return sns.NewFromConfig(cfg, func(o *sns.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("AWS_ENDPOINT"))
	}), nil
}
