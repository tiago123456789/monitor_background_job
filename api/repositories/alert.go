package repositories

import (
	"context"
	"os"

	"github.com/tiago123456789/monitor_background_job/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlertRepositoryInterface interface {
	Create(alert models.Alert) error
}

type AlertRepository struct {
	Client *mongo.Client
}

func NewAlertRepository(client *mongo.Client) *AlertRepository {
	return &AlertRepository{
		Client: client,
	}
}

func (a *AlertRepository) Create(alert models.Alert) error {
	collection := a.Client.Database(os.Getenv("DATABASE_NAME")).Collection("alerts")
	alert.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(context.TODO(), alert)
	if err != nil {
		return err
	}
	return nil
}
