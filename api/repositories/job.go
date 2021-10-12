package repositories

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/tiago123456789/monitor_background_job/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobsRepositoryInterface interface {
	Create(job models.JobModel) error
	FindByName(name string) (models.JobModel, error)
	FindByCompanyId(companyId string) ([]models.JobModel, error)
}

type JobRepository struct {
	Client            *mongo.Client
	CompanyRepository CompanyRepositoryInterface
}

func NewJobRepository(client *mongo.Client, companyRepository CompanyRepositoryInterface) *JobRepository {
	return &JobRepository{
		Client:            client,
		CompanyRepository: companyRepository,
	}
}

func (j *JobRepository) Create(job models.JobModel) error {
	result, _ := j.FindByName(job.Name)
	if result != (models.JobModel{}) {
		return errors.New("Name can't used")
	}

	company, _ := j.CompanyRepository.FindByID(job.Companyid)
	if company == (models.Company{}) {
		return errors.New("Company informated not exist")
	}

	collection := j.Client.Database(os.Getenv("DATABASE_NAME")).Collection("jobs")
	job.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(context.TODO(), job)
	if err != nil {
		return err
	}
	return nil
}

func (j *JobRepository) FindByName(name string) (models.JobModel, error) {
	var filter = &bson.D{
		{"name", name},
	}

	var result models.JobModel
	collection := j.Client.Database(os.Getenv("DATABASE_NAME")).Collection("jobs")
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return models.JobModel{}, err
	}
	return result, nil
}

func (j *JobRepository) FindByCompanyId(companyId string) ([]models.JobModel, error) {
	filter := &bson.D{
		{"companyid", companyId},
	}

	var results []models.JobModel
	collection := j.Client.Database(os.Getenv("DATABASE_NAME")).Collection("jobs")
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
