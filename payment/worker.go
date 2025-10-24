package payment

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bernardoazevedo/rinha-backend-2025/key"
	paymentqueue "github.com/bernardoazevedo/rinha-backend-2025/paymentQueue"
)

func PaymentWorker() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()


	client := key.GetNewClient()
	pubsub := client.Subscribe(ctx, paymentqueue.QueueName)
	channel := pubsub.Channel()

	go func() {
		var err error

		for {
			msg, ok := <-channel
			if !ok {
				break
			}

			_, err = postPayment([]byte(msg.Payload))
			if err != nil {
				fmt.Println("erro ao enviar: " + err.Error())
			}
		}
	}()

	defer pubsub.Close()

	fmt.Println("Monitoring payments...")
	<-sigchan

	fmt.Println("Killed, shutting down")
}
