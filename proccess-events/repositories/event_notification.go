package repositories

import (
	"context"
	"process-events/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventNotificationRepositoryInterface interface {
	Create(event models.EventNotification) error
}

type EventNotificationRepository struct {
	Client *mongo.Client
}

func NewEventNotificationRepository(client *mongo.Client) *EventNotificationRepository {
	return &EventNotificationRepository{
		Client: client,
	}
}

func (e *EventNotificationRepository) Create(event models.EventNotification) error {
	event.ID = primitive.NewObjectID()
	collection := e.Client.Database("monitor").Collection("notifications_received")
	_, err := collection.InsertOne(context.TODO(), event)
	if err == mongo.ErrNoDocuments {
		return err
	}
	return nil
}
