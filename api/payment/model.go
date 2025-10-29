package payment

import "time"

type Payment struct {
	CorrelationId string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
	RequestedAt   string  `json:"requestedAt"`
}

type Consumer struct {
	name   string
	count  int
	before time.Time
}