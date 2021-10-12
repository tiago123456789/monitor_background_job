package main

import (
	"encoding/json"
	"log"
	"os"
	httpClient "process-events/http"

	"process-events/config"
	"process-events/models"
	"process-events/queue"
	"time"

	"process-events/repositories"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
)

func triggerAlertTheCompanies() {
	mongo := config.NewConnection()
	companyRepository := repositories.NewCompany(mongo)
	producer := queue.NewProducer()

	uptimeTicker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-uptimeTicker.C:
			results, _ := companyRepository.GetAll()
			for _, value := range results {
				producer.Publish(os.Getenv("SQS_QUEUE_ALERT"), value)
			}
		}
	}
}

func handlerAlertTheCompany() {
	mongo := config.NewConnection()
	alertRepository := repositories.NewAlertRepository(mongo)
	lastNotificationRepository := repositories.NewLastEventNotificationRepository(mongo)

	queue.NewConsumer().Receive(os.Getenv("SQS_QUEUE_ALERT"), func(msgs ...*sqs.Message) {
		for _, msg := range msgs {
			var event models.CompanyTriggerAlert
			json.Unmarshal([]byte(*msg.Body), &event)

			alerts, err := alertRepository.FindByCompanyId(event.ID)
			if err != nil {
				log.Fatal(err)
			}

			if len(alerts) == 0 {
				continue
			}

			lastNotifications, err := lastNotificationRepository.FindByCompanyId(event.ID)
			if err != nil {
				log.Fatal(err)
			}

			mapLastNotificaiton := make(map[string]models.EventNotification)

			for _, value := range lastNotifications {
				id := value.ID.Hex()
				mapLastNotificaiton[id] = value
			}

			for i := 0; i < len(alerts); i++ {
				alert := alerts[i]
				if _, ok := mapLastNotificaiton[alert.JobId]; !ok {
					httpClient.Post(alert)
				}

				if _, ok := mapLastNotificaiton[alert.JobId]; ok {
					lastNotification := mapLastNotificaiton[alert.JobId]
					dateLastNotication := lastNotification.OccourAt.Add(
						time.Duration(alert.TimeInMinutes) * time.Minute)

					if dateLastNotication.Before(time.Now().UTC()) {
						httpClient.Post(alert)
					}
				}
			}
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

}
