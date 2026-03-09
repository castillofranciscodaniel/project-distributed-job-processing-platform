package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSConsumer struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSConsumer(ctx context.Context, region, queueURL string) (*SQSConsumer, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &SQSConsumer{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}

func (c *SQSConsumer) ReceiveMessages(ctx context.Context) ([]types.Message, error) {
	output, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: 1,  // Procesamos de a 1 para no colapsar con zips pesados
		WaitTimeSeconds:     20, // Long polling
	})
	if err != nil {
		return nil, fmt.Errorf("failed to receive messages: %w", err)
	}

	return output.Messages, nil
}

func (c *SQSConsumer) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}
