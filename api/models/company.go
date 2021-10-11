package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Company struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name     string             `json: "name"`
	Email    string             `json:"email"`
	Password string             `json:password`
}
