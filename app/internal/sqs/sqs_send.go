package sqs

import (
	"context"

	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/caiomarcatti12/poc-localstack/internal/utils"
)

// SendMessage envia uma mensagem simples com atributo de origem.
func SendMessage(ctx context.Context, client *awssqs.Client, queueURL, body string) error {
	_, err := client.SendMessage(ctx, &awssqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: &body,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Origin": {DataType: utils.AwsString("String"), StringValue: utils.AwsString("poc-localstack")},
		},
	})
	return err
}
