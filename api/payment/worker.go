package payment

import (
	"fmt"

	paymentqueue "github.com/bernardoazevedo/rinha-backend-2025/api/paymentQueue"
)

func PaymentWorker() {
	for {
		message := paymentqueue.Get()
		if len(message) > 0 {
			_, err := postPayment([]byte(message))
			if err != nil {
				fmt.Println("erro ao enviar: " + err.Error())
			}
		}
	}
}