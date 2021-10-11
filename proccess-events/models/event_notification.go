package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventNotification struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	JobId     string             `json:"jobId"`
	CompanyId string             `json:"companyId"`
	OccourAt  time.Time          `json:"occourAt"`
}
