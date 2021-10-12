package repositories

import (
	"context"
	"errors"
	"log"
	"os"
	"process-events/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanyInterface interface {
	GetAll() ([]models.CompanyTriggerAlert, error)
}

type Company struct {
	Client *mongo.Client
}

func NewCompany(client *mongo.Client) *Company {
	return &Company{
		Client: client,
	}
}

func (c *Company) GetAll() ([]models.CompanyTriggerAlert, error) {
	filter := &bson.D{}

	var results []models.CompanyTriggerAlert
	collection := c.Client.Database(os.Getenv("DATABASE_NAME")).Collection("companies")
	ctx := context.TODO()
	cursor, err := collection.Find(ctx, filter)
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
