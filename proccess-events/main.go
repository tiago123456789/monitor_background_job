package main

import (
	"encoding/json"
	"log"
	"process-events/config"
	"process-events/models"

	"process-events/repositories"

	"process-events/queue"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongo := config.NewConnection()
	eventNotificationRepository := repositories.NewEventNotificationRepository(mongo)

	queue.NewConsumer().Receive(func(msgs ...*sqs.Message) {
		for _, msg := range msgs {
			var event models.EventNotification
			json.Unmarshal([]byte(*msg.Body), &event)
			eventNotificationRepository.Create(event)
		}
	})
}
