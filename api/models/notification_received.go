package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationReceived struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	JobId     string             `bson:"jobid" json: "name"`
	CompanyId string             `bson:"companyid" json:"companyId"`
	OccourAt  time.Time          `bson:"occourat" json:occourAt`
}
