package paymentqueue

var paymentQueue []string

func Push(paymentJson string) {
	paymentQueue = append(paymentQueue, paymentJson)
}

func Pop() string {
	if len(paymentQueue) > 0 {
		paymentJson := paymentQueue[0]
		paymentQueue = paymentQueue[1:]  
		return paymentJson
	} else {
		return ""
	}
}
