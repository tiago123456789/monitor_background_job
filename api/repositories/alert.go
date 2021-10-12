package repositories

import (
	"context"
	"os"

	"github.com/tiago123456789/monitor_background_job/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlertRepositoryInterface interface {
	Create(alert models.Alert) error
	FindByCompanyIdAndJobId(companyId string, jobId string) ([]models.Alert, error)
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

func (a *AlertRepository) FindByCompanyIdAndJobId(companyId string, jobId string) ([]models.Alert, error) {
	filter := &bson.D{
		{"companyId", companyId},
		{"jobId", jobId},
	}

	var results []models.Alert
	collection := a.Client.Database(os.Getenv("DATABASE_NAME")).Collection("alerts")
	ctx := context.TODO()
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return []models.Alert{}, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return []models.Alert{}, err
	}

	return results, nil
}
