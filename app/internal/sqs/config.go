package sqs

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/caiomarcatti12/poc-localstack/internal/types"
)

// BuildAWSConfig constrói configuração AWS com endpoint LocalStack customizado.
func BuildAWSConfig(ctx context.Context, cfg types.EnvConfig) (aws.Config, error) {
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
		if cfg.Endpoint != "" && service == awssqs.ServiceID {
			return aws.Endpoint{URL: cfg.Endpoint, SigningRegion: cfg.Region, HostnameImmutable: true}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
	return config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithEndpointResolverWithOptions(resolver),
	)
}

// EnsureQueue garante que a fila exista e retorna a URL.
func EnsureQueue(ctx context.Context, client *awssqs.Client, queueName string) (string, error) {
	getOut, err := client.GetQueueUrl(ctx, &awssqs.GetQueueUrlInput{QueueName: &queueName})
	if err == nil && getOut != nil && getOut.QueueUrl != nil {
		return *getOut.QueueUrl, nil
	}
	createOut, cErr := client.CreateQueue(ctx, &awssqs.CreateQueueInput{
		QueueName: &queueName,
		Attributes: map[string]string{
			string(sqstypes.QueueAttributeNameVisibilityTimeout): "30",
		},
	})
	if cErr != nil {
		return "", fmt.Errorf("criando fila %s: %w", queueName, cErr)
	}
	if createOut.QueueUrl == nil {
		return "", errors.New("CreateQueue retornou URL nula")
	}
	return *createOut.QueueUrl, nil
}
