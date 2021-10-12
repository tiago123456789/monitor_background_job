package queue

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type ConsumerInterface interface {
	Receive() error
}

type Consumer struct {
}

func (c *Consumer) getQueueURL(svc *sqs.SQS, queue *string) (*sqs.GetQueueUrlOutput, error) {
	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}
	return urlResult, nil
}

func (c *Consumer) getMessages(svc *sqs.SQS, queueURL *string) (*sqs.ReceiveMessageOutput, error) {
	var waitTime int64 = 20
	msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            queueURL,
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     &waitTime,
	})

	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

func (c *Consumer) deleteMessage(sess *sqs.SQS, queue *string, msgs ...*sqs.Message) error {
	delParams := sqs.DeleteMessageInput{
		QueueUrl: queue,
	}

	for _, msg := range msgs {
		delParams.ReceiptHandle = msg.ReceiptHandle
		_, err := sess.DeleteMessage(&delParams)
		if err != nil {
			return err
		}

	}

	return nil
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Receive(callback func(msgs ...*sqs.Message)) {
	queueName := os.Getenv("SQS_QUEUE")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	urlResult, err := c.getQueueURL(svc, &queueName)
	if err != nil {
		fmt.Println(err)
		return
	}

	queueURL := urlResult.QueueUrl
	for {
		msgResult, err := c.getMessages(svc, queueURL)
		callback(msgResult.Messages...)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(msgResult.Messages) > 0 {
			c.deleteMessage(svc, queueURL, msgResult.Messages...)
		}
	}
}
