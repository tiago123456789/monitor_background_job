package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Alert struct {
	ID            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	JobId         string             `bson:"jobId" json:"jobId"`
	CompanyId     string             `bson:"companyId" json:"companyId"`
	TimeInMinutes int16              `bson:"timeInMinutes" json:"timeInMinutes"`
	Url           string             `bson:"url" json:"url"`
	Payload       string             `bson:"payload" json:"payload"`
}
