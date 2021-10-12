package queue

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type ProducerInterface interface {
	Publish(message interface{}) error
}

type Producer struct {
}

func NewProducer() *Producer {
	return &Producer{}
}

func (p *Producer) getQueueURL(svc *sqs.SQS, queue *string) (*sqs.GetQueueUrlOutput, error) {
	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}
	return urlResult, nil
}

func (p *Producer) Publish(queue string, message interface{}) error {
	messageString, err := json.Marshal(message)
	if err != nil {
		return err
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	urlResult, err := p.getQueueURL(svc, &queue)
	if err != nil {
		return err
	}
	queueURL := urlResult.QueueUrl

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(messageString)),
		QueueUrl:    queueURL,
	})
	if err != nil {
		return err
	}

	return nil
}
