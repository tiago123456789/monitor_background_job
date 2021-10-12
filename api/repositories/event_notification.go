package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/tiago123456789/monitor_background_job/config"
	"github.com/tiago123456789/monitor_background_job/models"
	"github.com/tiago123456789/monitor_background_job/queue"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventNotificationRepositoryInterface interface {
	Get(id string) (string, error)
	StoreLast(eventNotification models.EventNotification) error
	GetByJobID(jobId string) ([]models.NotificationReceived, error)
}

type EventNotificationRepository struct {
	Cache    *config.Cache
	Producer queue.ProducerInterface
	Client   *mongo.Client
}

func NewEventNotificationRepository(
	cache *config.Cache,
	producer queue.ProducerInterface,
	client *mongo.Client) *EventNotificationRepository {

	return &EventNotificationRepository{
		Cache: cache, Producer: producer, Client: client,
	}
}

func (e *EventNotificationRepository) Get(id string) (string, error) {
	value, err := e.Cache.Get(context.TODO(), id)
	if err != nil {
		return "", err
	}

	jsonEventNotification := models.EventNotification{}
	err = json.Unmarshal([]byte(value), &jsonEventNotification)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (e *EventNotificationRepository) StoreLast(id string, jobName string) error {
	v, err := e.Get(id)
	jsonEventNotification := make(map[string]models.EventNotification)
	if err == nil {
		err = json.Unmarshal([]byte(v), &jsonEventNotification)
		if err != nil {
			return err
		}
	}

	var events map[string]models.EventNotification
	events = make(map[string]models.EventNotification)
	events = jsonEventNotification

	date := time.Now().UTC()
	events[jobName] = models.EventNotification{
		ID:        id,
		JobName:   jobName,
		CreatedAt: date,
	}

	valueStored, err := json.Marshal(&events)
	if err != nil {
		return err
	}

	e.Cache.Set(context.TODO(), id, string(valueStored))
	e.Cache.Client.Expire(context.TODO(), id, 10*time.Minute)

	eventNotificationMessage := &models.EventNotificationMessage{
		CompanyId: id,
		JobId:     jobName,
		OccourAt:  date,
	}
	err = e.Producer.Publish(eventNotificationMessage)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventNotificationRepository) GetByJobID(jobId string) ([]models.NotificationReceived, error) {
	filter := &bson.D{
		{"jobid", jobId},
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"occourat", -1}})

	var results []models.NotificationReceived
	collection := e.Client.Database(os.Getenv("DATABASE_NAME")).Collection("notifications_received")
	ctx := context.TODO()
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("Empty results")
	}
	return results, nil
}
