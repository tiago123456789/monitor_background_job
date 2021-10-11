package queue

import (
	"encoding/json"
	"os"

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

func (p *Producer) Publish(message interface{}) error {
	messageString, err := json.Marshal(message)
	if err != nil {
		return err
	}

	queue := os.Getenv("SQS_QUEUE")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	result, err := GetQueueURL(sess, &queue)
	if err != nil {
		return err
	}
	queueURL := result.QueueUrl

	svc := sqs.New(sess)

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(messageString)),
		QueueUrl:    queueURL,
	})
	if err != nil {
		return err
	}

	return nil
}
