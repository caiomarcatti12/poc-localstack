package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func main() {
	ctx := context.Background()
	endpoint := os.Getenv("LOCALSTACK_ENDPOINT")
	queueName := os.Getenv("QUEUE_NAME")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if endpoint != "" && service == sqs.ServiceID { // ServiceID is "SQS"
			return aws.Endpoint{URL: endpoint, PartitionID: "aws", SigningRegion: os.Getenv("AWS_REGION")}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatalf("erro carregando config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	queueUrl, err := ensureQueue(ctx, sqsClient, queueName)
	if err != nil {
		log.Fatalf("erro garantindo fila: %v", err)
	}
	log.Printf("Usando fila %s (%s)", queueName, queueUrl)

	// Envia uma mensagem
	body := fmt.Sprintf("Olá LocalStack! %s", time.Now().Format(time.RFC3339))
	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: &body,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Origin": {DataType: awsString("String"), StringValue: awsString("poc-localstack")},
		},
	})
	if err != nil {
		log.Fatalf("erro enviando mensagem: %v", err)
	}
	log.Printf("Mensagem enviada: %s", body)

	// Recebe mensagem
	out, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     2,
		MessageAttributeNames: []string{
			"All",
		},
	})
	if err != nil {
		log.Fatalf("erro recebendo mensagem: %v", err)
	}
	for _, m := range out.Messages {
		log.Printf("Mensagem recebida ID=%s Body=%s", awsStringValue(m.MessageId), awsStringValue(m.Body))
		// Deleta
		_, derr := sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{QueueUrl: &queueUrl, ReceiptHandle: m.ReceiptHandle})
		if derr != nil {
			log.Printf("erro deletando mensagem: %v", derr)
		} else {
			log.Printf("Mensagem deletada: %s", awsStringValue(m.MessageId))
		}
	}
}

func ensureQueue(ctx context.Context, client *sqs.Client, name string) (string, error) {
	// tenta obter URL primeiro
	getOut, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: &name})
	if err == nil && getOut.QueueUrl != nil {
		return *getOut.QueueUrl, nil
	}
	// se não existe, cria
	var notFound bool
	if err != nil {
		var nf *types.QueueDoesNotExist
		if errors.As(err, &nf) {
			notFound = true
		}
	}
	if notFound {
		_, cErr := client.CreateQueue(ctx, &sqs.CreateQueueInput{QueueName: &name})
		if cErr != nil {
			return "", fmt.Errorf("falha criando fila: %w", cErr)
		}
		getOut, err = client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: &name})
		if err != nil || getOut.QueueUrl == nil {
			return "", fmt.Errorf("não foi possível obter URL da fila após criação: %w", err)
		}
		return *getOut.QueueUrl, nil
	}
	if err != nil { // outro erro
		return "", fmt.Errorf("erro obtendo fila: %w", err)
	}
	return "", fmt.Errorf("estado inesperado ao garantir fila")
}

func awsString(s string) *string { return &s }
func awsStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
