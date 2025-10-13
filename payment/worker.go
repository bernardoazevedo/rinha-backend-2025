package payment

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	paymentqueue "github.com/bernardoazevedo/rinha-backend-2025/paymentQueue"
)

func PaymentWorker() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			paymentJson := paymentqueue.Pop()
			payment := []byte(paymentJson)

			if paymentJson != "" {
				_, err := postPayment(payment)
				if err != nil {
					fmt.Println("error: " + err.Error())
				}
			} else {
				time.Sleep(time.Second / 2)
			}
		}
	}()

	fmt.Println("Monitoring payments...")
	<-sigchan

	fmt.Println("Killed, shutting down")
}
