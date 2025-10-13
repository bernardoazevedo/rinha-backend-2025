package paymentqueue

var paymentQueue [][]byte

func Push(payment []byte) {
	paymentQueue = append(paymentQueue, payment)
}

func Pop() []byte {
	var payment []byte

	if len(paymentQueue) > 0 {
		payment = paymentQueue[0]
		paymentQueue = paymentQueue[1:]  
	}

	return payment
}
