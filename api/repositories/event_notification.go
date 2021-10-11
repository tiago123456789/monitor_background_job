package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tiago123456789/monitor_background_job/config"
	"github.com/tiago123456789/monitor_background_job/models"
	"github.com/tiago123456789/monitor_background_job/queue"
)

type EventNotificationRepositoryInterface interface {
	Get(id string) (string, error)
	StoreLast(eventNotification models.EventNotification) error
}

type EventNotificationRepository struct {
	Cache    *config.Cache
	Producer queue.ProducerInterface
}

func NewEventNotificationRepository(
	cache *config.Cache,
	producer queue.ProducerInterface) *EventNotificationRepository {

	return &EventNotificationRepository{
		Cache: cache, Producer: producer,
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
