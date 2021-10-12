package repositories

import (
	"context"
	"os"
	"process-events/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlertRepositoryInterface interface {
	FindByCompanyId(companyId string, jobId string) ([]models.Alert, error)
}

type AlertRepository struct {
	Client *mongo.Client
}

func NewAlertRepository(client *mongo.Client) *AlertRepository {
	return &AlertRepository{
		Client: client,
	}
}

func (a *AlertRepository) FindByCompanyId(companyId string) ([]models.Alert, error) {
	filter := &bson.D{
		{"companyId", companyId},
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
