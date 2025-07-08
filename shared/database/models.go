package database

import "time"

type Mail struct {
	ID         int
	Sender     string
	Recipients []string
	RawData    string
	ReceivedAt time.Time
}
