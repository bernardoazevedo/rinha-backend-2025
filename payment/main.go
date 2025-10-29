package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bernardoazevedo/rinha-backend-2025/api/key"
	"github.com/bernardoazevedo/rinha-backend-2025/api/payment"
)

func main() {
	key.GetNewClient()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for i := 0; i < 3; i++ {
		go payment.PaymentWorker()
	}

	fmt.Println("Monitoring payments...")
	<-sigchan

	fmt.Println("Killed, shutting down")
}
