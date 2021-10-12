package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"process-events/config"
	"process-events/models"
	"process-events/queue"
	"time"

	"process-events/repositories"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
)

type Alert struct {
	TimeInMinutes int16
	Url           string
	Payload       string
}

func triggerAlertTheCompanies() {
	mongo := config.NewConnection()
	companyRepository := repositories.NewCompany(mongo)
	producer := queue.NewProducer()

	uptimeTicker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-uptimeTicker.C:
			results, _ := companyRepository.GetAll()
			for _, value := range results {
				producer.Publish(os.Getenv("SQS_QUEUE_ALERT"), value)
				fmt.Printf("Trigger alert the company %v \n", value.ID)
			}
		}
	}
}

func post(alert Alert) {
	req, err := http.NewRequest("POST", alert.Url, bytes.NewBuffer([]byte(alert.Payload)))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func handlerAlertTheCompany() {

	alertTest := &Alert{
		TimeInMinutes: 5,
		Url:           "https://webhook.site/2d1c2902-5153-4867-9bae-6f53a70a7624",
		Payload:       "{ \"message\": \"text\" }",
	}

	queue.NewConsumer().Receive(os.Getenv("SQS_QUEUE_ALERT"), func(msgs ...*sqs.Message) {
		for _, msg := range msgs {
			var event models.CompanyTriggerAlert
			json.Unmarshal([]byte(*msg.Body), &event)
			fmt.Printf("Alert company %v \n", event)
			post(*alertTest)
		}
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongo := config.NewConnection()
	eventNotificationRepository := repositories.NewEventNotificationRepository(mongo)
	lastEventNotificationRepository := repositories.NewLastEventNotificationRepository(mongo)
	go triggerAlertTheCompanies()
	go handlerAlertTheCompany()

	queue.NewConsumer().Receive(os.Getenv("SQS_QUEUE"), func(msgs ...*sqs.Message) {
		for _, msg := range msgs {
			var event models.EventNotification
			json.Unmarshal([]byte(*msg.Body), &event)
			eventNotificationRepository.Create(event)
			lastEventNotificationRepository.Register(event)
		}
	})

	// queue.NewConsumer().Receive(os.Getenv("SQS_QUEUE_ALERT"), func(msgs ...*sqs.Message) {
	// 	for _, msg := range msgs {
	// 		var event models.CompanyTriggerAlert
	// 		json.Unmarshal([]byte(*msg.Body), &event)
	// 		fmt.Printf("Trigger alert the company %v \n", event)

	// 		// eventNotificationRepository.Create(event)
	// 		// lastEventNotificationRepository.Register(event)
	// 	}
	// })

	// https://webhook.site/2d1c2902-5153-4867-9bae-6f53a70a7624
}
