package broker

import "time"

type Message struct {
	Sequence  int       `json:"id"`
	Message   string    `json:"message"`
	Nick      string    `json:"nick"`
	Channel   string    `json:"channel"`
	Platform  string    `json:"platform"`
	Timestamp time.Time `json:"timestamp"`
}
