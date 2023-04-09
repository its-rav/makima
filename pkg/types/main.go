package types

import "time"

type PublishMessage struct {
	Channel   string
	Message   string
	Timestamp time.Time
}
