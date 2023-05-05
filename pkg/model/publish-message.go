package model

import "time"

type PublishMessage[T any] struct {
	Message     T
	Source      string
	Destination string
	Extras      map[string]string
	Timestamp   time.Time
}
