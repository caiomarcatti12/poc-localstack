package sqs

import (
	"context"
	"fmt"

	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/caiomarcatti12/poc-localstack/internal/types"
	"github.com/caiomarcatti12/poc-localstack/internal/utils"
)

// ReceiveOneAndDelete realiza um receive de no máximo 1 mensagem e deleta se encontrada.
func ReceiveOneAndDelete(ctx context.Context, client *awssqs.Client, queueURL string) (*types.ReceivedMessage, error) {
	out, err := client.ReceiveMessage(ctx, &awssqs.ReceiveMessageInput{
		QueueUrl:            &queueURL,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     2,
		MessageAttributeNames: []string{
			"All",
		},
	})
	if err != nil {
		return nil, err
	}
	if len(out.Messages) == 0 {
		return nil, nil
	}
	m := out.Messages[0]
	_, derr := client.DeleteMessage(ctx, &awssqs.DeleteMessageInput{QueueUrl: &queueURL, ReceiptHandle: m.ReceiptHandle})
	if derr != nil {
		return nil, fmt.Errorf("falha ao deletar mensagem %s: %w", utils.AwsStringValue(m.MessageId), derr)
	}
	return &types.ReceivedMessage{ID: utils.AwsStringValue(m.MessageId), Body: utils.AwsStringValue(m.Body)}, nil
}
