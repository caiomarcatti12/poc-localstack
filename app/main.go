package main

import (
	"context"
	"fmt"
	"log"
	"time"

	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/caiomarcatti12/poc-localstack/internal/env"
	"github.com/caiomarcatti12/poc-localstack/internal/sqs"
)

// ------------------------------ Main ----------------------------------------
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cfg := env.LoadEnvConfig()

	awsCfg, err := sqs.BuildAWSConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("config AWS: %v", err)
	}

	client := awssqs.NewFromConfig(awsCfg)

	queueURL, err := sqs.EnsureQueue(ctx, client, cfg.QueueName)
	if err != nil {
		log.Fatalf("garantindo fila: %v", err)
	}
	log.Printf("Fila pronta: %s (%s)", cfg.QueueName, queueURL)

	body := fmt.Sprintf("Olá LocalStack! %s", time.Now().Format(time.RFC3339))
	if err := sqs.SendMessage(ctx, client, queueURL, body); err != nil {
		log.Fatalf("enviando mensagem: %v", err)
	}
	log.Printf("Mensagem enviada: %s", body)

	msg, err := sqs.ReceiveOneAndDelete(ctx, client, queueURL)
	if err != nil {
		log.Fatalf("recebendo mensagem: %v", err)
	}
	if msg != nil {
		log.Printf("Mensagem processada e removida: ID=%s Body=%s", msg.ID, msg.Body)
	} else {
		log.Printf("Nenhuma mensagem disponível")
	}
}
