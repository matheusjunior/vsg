package messaging

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"matheus.com/vgs/internal/logger"
)

type (
	SQSListener struct {
		queueUrl string
		client   *sqs.Client
		stop     chan struct{}
		messages chan RawMessage
	}

	SQSPublisher struct {
		queueUrl string
		client   *sqs.Client
	}

	RawMessage string
)

func NewSQSPublisher(queueUrl string) *SQSPublisher {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://localhost:4566",
			SigningRegion: "us-east-1",
		}, nil
	})
	config, err := config.LoadDefaultConfig(context.Background(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		logger.Logger().Fatal(err)
	}

	return &SQSPublisher{
		queueUrl: queueUrl,
		client:   sqs.NewFromConfig(config),
	}
}

func (p *SQSPublisher) Publish(msg interface{}) {
	payload, err := json.Marshal(msg)
	if err != nil {
		logger.Logger().Error("could not send message")
	}
	input := sqs.SendMessageInput{
		DelaySeconds: 2,
		MessageBody:  aws.String(string(payload)),
		QueueUrl:     &p.queueUrl,
	}
	resp, err := p.client.SendMessage(context.TODO(), &input)
	if err != nil {
		logger.Logger().Error(err)
	}
	logger.Logger().Info("Sent message with ID:", *resp.MessageId)
}

func NewSQSListener(queueUrl string) *SQSListener {
	// TODO make it abstract
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://localhost:4566",
			SigningRegion: "us-east-1",
		}, nil
	})
	config, err := config.LoadDefaultConfig(context.Background(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		logger.Logger().Fatal(err)
	}

	return &SQSListener{
		queueUrl: queueUrl,
		client:   sqs.NewFromConfig(config),
		stop:     make(chan struct{}),
		messages: make(chan RawMessage),
	}
}

func (l *SQSListener) Start() <-chan RawMessage {
	ctx := context.Background()
	go func() {
		for {
			select {
			case <-l.stop:
				logger.Logger().Warn("client stopped: closing message channel")
				close(l.messages)
				close(l.stop)
				return
			default:
				response, err := l.fetch(ctx)
				if err != nil {
					logger.Logger().Error(err)
					continue
				}

				l.dispatch(ctx, response)
			}
		}
	}()

	return l.messages
}

func (l *SQSListener) Stop() {
	l.stop <- struct{}{}
}

func (l *SQSListener) fetch(ctx context.Context) (*sqs.ReceiveMessageOutput, error) {
	request := sqs.ReceiveMessageInput{QueueUrl: &l.queueUrl, WaitTimeSeconds: 2}
	return l.client.ReceiveMessage(ctx, &request)
}

func (l *SQSListener) dispatch(ctx context.Context, response *sqs.ReceiveMessageOutput) {
	for _, m := range response.Messages {
		if m.Body != nil {
			l.messages <- RawMessage(*m.Body)
		}

		// TODO do it in batch
		_, err := l.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{QueueUrl: &l.queueUrl, ReceiptHandle: m.ReceiptHandle})
		if err != nil {
			logger.Logger().Error("failed deleting message")
		}
	}
}
