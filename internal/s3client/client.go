package s3client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/indor79/s3bench/internal/config"
)

func New(ctx context.Context, c config.Config) (*s3.Client, error) {
	access := strings.TrimSpace(os.Getenv(c.Auth.AccessKeyEnv))
	secret := strings.TrimSpace(os.Getenv(c.Auth.SecretKeyEnv))
	session := strings.TrimSpace(os.Getenv(c.Auth.SessionTokenEnv))
	if access == "" || secret == "" {
		return nil, fmt.Errorf("missing credentials in env vars %s/%s", c.Auth.AccessKeyEnv, c.Auth.SecretKeyEnv)
	}

	cfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(c.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(access, secret, session)),
	)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}

	cli := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.HTTPClient = httpClient
		if strings.TrimSpace(c.Endpoint) != "" {
			o.BaseEndpoint = aws.String(c.Endpoint)
		}
	})
	return cli, nil
}
