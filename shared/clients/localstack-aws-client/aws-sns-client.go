package localstackclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func NewSNSClient(ctx context.Context) (*sns.Client, error) {

	// l, err := localstack.NewInstance()
	// if err != nil {
	// 	log.Fatal("Could not connect to Docker %v", err)
	// }
	// if err := l.Start(); err != nil {
	// 	log.Fatal("Could not start localstack %v", err)
	// }

	// cfg, err := awsconfig.LoadDefaultConfig(ctx,
	// 	awsconfig.WithRegion("us-east-1"),
	// 	awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
	// )
	// if err != nil {
	// 	log.Fatal("Could not get config %v", err)
	// }
	// resolver := localstack.NewSnsResolverV2(l)
	// client := sns.NewFromConfig(cfg, sns.WithEndpointResolverV2(resolver))
	// return client, nil

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "test")),
	)
	if err != nil {
		panic(err)
	}

	return sns.NewFromConfig(cfg, func(o *sns.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
	}), nil
}
