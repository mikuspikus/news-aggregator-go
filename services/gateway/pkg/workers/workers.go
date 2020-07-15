package workers

import "time"

type Message struct {
	UserUID   string    `json:"user_uid"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
}
