#!/bin/bash
set -euo pipefail

QUEUE_NAME=${QUEUE_NAME:-demo-queue}
REGION=${AWS_REGION:-us-east-1}

awslocal sqs create-queue --queue-name "$QUEUE_NAME" --attributes VisibilityTimeout=30

echo "Queue '$QUEUE_NAME' criada"
