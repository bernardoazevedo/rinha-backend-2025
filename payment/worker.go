package payment

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	paymentqueue "github.com/bernardoazevedo/rinha-de-backend-2025/paymentQueue"
)

func PaymentWorker() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		paymentJson := paymentqueue.Pop()

		_, err := postPayment(paymentJson)
		if err != nil {
			fmt.Println("error: " + err.Error())
		}
	}()

	fmt.Println("Monitoring services health...")
	<-sigchan

	fmt.Println("Killed, shutting down")
}