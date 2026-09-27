package r2

import (
	"context"
	"sync"

	"music-server/src/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	client   *s3.Client
	connOnce sync.Once
)

func get_client() *s3.Client {
	connOnce.Do(func() {
		//var bucket_name = config.GetEnv().R2BucketName
		var endpoint = config.Get().R2Endpoint
		var access_key_id = config.Get().R2AccessKeyID
		var secret_access_key = config.Get().R2SecretAccessKey

		cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
			awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(access_key_id, secret_access_key, "")),
			awsconfig.WithRegion("auto"),
		)

		if err != nil {
			panic(err)
		}

		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	})

	return client
}