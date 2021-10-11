package repositories

import (
	"context"
	"errors"

	"github.com/tiago123456789/monitor_background_job/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type CompanyRepositoryInterface interface {
	Create(company models.Company) error
	FindByName(name string) (models.Company, error)
	FindByID(id string) (models.Company, error)
}

type CompanyRepository struct {
	Client *mongo.Client
}

func NewCompanyRepostory(client *mongo.Client) *CompanyRepository {
	return &CompanyRepository{
		Client: client,
	}
}

func (c *CompanyRepository) Create(company models.Company) error {
	result, _ := c.FindByName(company.Name)
	if result != (models.Company{}) {
		return errors.New("Name can't used")
	}
	collection := c.Client.Database("monitor").Collection("companies")
	password, err := bcrypt.GenerateFromPassword([]byte(company.Password), 10)
	if err != nil {
		return err
	}

	company.Password = string(password)
	company.ID = primitive.NewObjectID()
	_, err = collection.InsertOne(context.TODO(), company)
	if err != nil {
		return err
	}
	return nil
}

func (c *CompanyRepository) FindByName(name string) (models.Company, error) {
	var filter = &bson.D{
		{"name", name},
	}

	var result models.Company
	collection := c.Client.Database("monitor").Collection("companies")
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return models.Company{}, err
	}
	return result, nil
}

func (c *CompanyRepository) FindByID(id string) (models.Company, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": bson.M{"$eq": objID}}

	var result models.Company
	collection := c.Client.Database("monitor").Collection("companies")
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return models.Company{}, err
	}
	return result, nil
}
