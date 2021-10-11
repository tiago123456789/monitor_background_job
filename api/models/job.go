package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type JobModel struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name      string             `bson:"name" json: "name"`
	Interval  int16              `bson:"interval" json: "interval"`
	Companyid string             `bson:"companyid" json: "companyId"`
}
