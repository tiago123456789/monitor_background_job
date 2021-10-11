package models

import "time"

type EventNotification struct {
	ID        string    `json:"id"`
	JobName   string    `json:"jobName"`
	CreatedAt time.Time `json:"createdAt"`
}
