package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSPublisher struct {
	client   *sns.Client
	topicArn string
}

func NewSNSPublisher(ctx context.Context, region, topicArn string) (*SNSPublisher, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &SNSPublisher{
		client:   sns.NewFromConfig(cfg),
		topicArn: topicArn,
	}, nil
}

func (p *SNSPublisher) Publish(ctx context.Context, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = p.client.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(string(body)),
		TopicArn: aws.String(p.topicArn),
	})

	if err != nil {
		return fmt.Errorf("failed to publish to SNS: %w", err)
	}

	return nil
}
