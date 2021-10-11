package models

import "time"

type EventNotificationMessage struct {
	CompanyId string
	JobId     string
	OccourAt  time.Time
}
