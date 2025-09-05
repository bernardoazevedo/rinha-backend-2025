package payment

type Payment struct {
	CorrelationId string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
	RequestedAt   string  `json:"requestedAt"`
}