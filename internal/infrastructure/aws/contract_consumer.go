package aws

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/francisco/distributed-job-platform/internal/domain/contract"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContractZipConsumer struct {
	sqsConsumer   *SQSConsumer
	workerService *contract.ContractWorkerService
}

func NewContractZipConsumer(sqsConsumer *SQSConsumer, service *contract.ContractWorkerService) *ContractZipConsumer {
	return &ContractZipConsumer{
		sqsConsumer:   sqsConsumer,
		workerService: service,
	}
}

func (c *ContractZipConsumer) Start(ctx context.Context) {
	log.Println("Contract Consumer started, waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Contract Consumer...")
			return
		default:
			messages, err := c.sqsConsumer.ReceiveMessages(ctx)
			if err != nil {
				log.Printf("Error receiving messages from SQS: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for _, msg := range messages {
				c.processMessage(ctx, msg.Body, msg.ReceiptHandle, msg.MessageId)
			}
		}
	}
}

func (c *ContractZipConsumer) processMessage(ctx context.Context, body, receiptHandle, messageID *string) {
	log.Printf("Processing message ID: %s", *messageID)

	// 1. Unmarshal SNS Envelope
	var snsEnvelope struct {
		Message string `json:"Message"`
	}
	if err := json.Unmarshal([]byte(*body), &snsEnvelope); err != nil {
		log.Printf("Error unmarshaling SNS envelope: %v", err)
		return
	}

	// 2. Unmarshal Domain Event
	var event contract.ZippingRequestedEvent
	if err := json.Unmarshal([]byte(snsEnvelope.Message), &event); err != nil {
		log.Printf("Error unmarshaling domain event: %v", err)
		return
	}

	clientID, err := primitive.ObjectIDFromHex(event.ClientID)
	if err != nil {
		log.Printf("Invalid ClientID in event: %v", err)
		return
	}

	// 3. Invoke Service Logic
	log.Printf("Executing zipping job for Client: %s", clientID.Hex())
	if err := c.workerService.ProcessZippingRequest(ctx, clientID); err != nil {
		log.Printf("Job failed for Client %s: %v", clientID.Hex(), err)
		return
	}

	// 4. Delete message on success
	if err := c.sqsConsumer.DeleteMessage(ctx, *receiptHandle); err != nil {
		log.Printf("Failed to delete message %s from SQS: %v", *messageID, err)
	}

	log.Printf("Successfully processed and deleted message for client: %s", clientID.Hex())
}
