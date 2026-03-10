package sqsaws

import (
	"context"
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type sqsAWS struct {
	client          *sqs.Client
	lgs             adapters.Logs
	SqsQueueURL     string
	Endpoint        string
	RetrieMode      string
	RoleARN         string
	ExternaID       string
}

func NewSqsAWS(
	MaxAttempts int64,
	Region string,
	AccessKeyID string,
	SecretAccessKey string,
	SqsQueueURL string,
	Endpoint string,
	RetrieMode string,
	RoleARN string,
	ExternaID string,
	lgs adapters.Logs,
) adapters.MemoryClient {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(Region),
		config.WithRetryMaxAttempts(int(MaxAttempts)),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				AccessKeyID,
				SecretAccessKey,
				"",
			),
		),
	)

	if err != nil {
		panic(err)
	}

	client := sqs.NewFromConfig(cfg)

	return &sqsAWS{
		SqsQueueURL:     SqsQueueURL,
		Endpoint:        Endpoint,
		RetrieMode:      RetrieMode,
		RoleARN:         RoleARN,
		ExternaID:       ExternaID,
		lgs:             lgs,
		client:          client,
	}
}

func (s *sqsAWS) SendMessage(ctx context.Context, q *entities.Queue) error {

	return nil
}

func (s *sqsAWS) DeleteMessage(ctx context.Context, ID string) error {

	return nil
}

func (s *sqsAWS) GetMessage(ctx context.Context, key string) (*entities.Queue, error) {

	return nil, nil
}

func (s *sqsAWS) GetKeys(ctx context.Context, mem entities.Memory) error {
	return nil
}
