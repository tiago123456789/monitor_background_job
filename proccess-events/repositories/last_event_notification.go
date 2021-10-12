package repositories

import (
	"context"
	"os"
	"process-events/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LastEventNotificationRepositoryInterface interface {
	Register(event models.EventNotification) error
	FindByCompanyId(companyId string) ([]models.EventNotification, error)
}

type LastEventNotificationRepository struct {
	Client *mongo.Client
}

func NewLastEventNotificationRepository(client *mongo.Client) *LastEventNotificationRepository {
	return &LastEventNotificationRepository{
		Client: client,
	}
}

func (e *LastEventNotificationRepository) isExist(jobId string, companyId string) bool {
	filter := &bson.D{
		{"companyid", companyId},
		{"jobid", jobId},
	}

	var results []models.EventNotification
	collection := e.Client.Database(os.Getenv("DATABASE_NAME")).Collection("last_notifications_received")
	ctx := context.TODO()
	err := collection.FindOne(ctx, filter).Decode(&results)
	if err == mongo.ErrNoDocuments {
		return false
	}

	return true
}

func (e *LastEventNotificationRepository) update(
	jobId string, companyId string, occourAt time.Time) error {
	filter := &bson.D{
		{"companyid", companyId},
		{"jobid", jobId},
	}

	dataModified := bson.M{
		"$set": bson.M{
			"occourat": occourAt,
		},
	}

	collection := e.Client.Database(os.Getenv("DATABASE_NAME")).Collection("last_notifications_received")
	ctx := context.TODO()
	_, err := collection.UpdateOne(ctx, filter, dataModified)
	if err == mongo.ErrNoDocuments {
		return err
	}

	return nil
}

func (e *LastEventNotificationRepository) Register(event models.EventNotification) error {
	if e.isExist(event.JobId, event.CompanyId) {
		e.update(event.JobId, event.CompanyId, event.OccourAt)
		return nil
	}

	event.ID = primitive.NewObjectID()
	collection := e.Client.Database(os.Getenv("DATABASE_NAME")).Collection("last_notifications_received")
	_, err := collection.InsertOne(context.TODO(), event)
	if err == mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func (e *LastEventNotificationRepository) FindByCompanyId(companyId string) ([]models.EventNotification, error) {
	filter := &bson.D{
		{"companyId", companyId},
	}

	var results []models.EventNotification
	collection := e.Client.Database(os.Getenv("DATABASE_NAME")).Collection("last_notifications_received")
	ctx := context.TODO()
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return []models.EventNotification{}, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return []models.EventNotification{}, err
	}

	return results, nil
}
