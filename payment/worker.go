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

// func PaymentWorker() {
// 	sigchan := make(chan os.Signal, 1)
// 	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

// 	go func() {
// 		var err error

// 		for {
// 			payment := paymentqueue.Get()
// 			if payment != "" {
// 				_, err = postPayment([]byte(payment))
// 				if err != nil {
// 					fmt.Println("erro: " + err.Error())
// 				}
// 			} else {
// 				time.Sleep(time.Second)
// 			}
// 		}
// 	}()

// 	fmt.Println("Monitoring payments...")
// 	<-sigchan

// 	fmt.Println("Killed, shutting down")
// }

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
			// msg, err := pubsub.ReceiveMessage(ctx)			
			// if err != nil {
			// 	fmt.Println("erro ao receber: " + err.Error())
			// }

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
