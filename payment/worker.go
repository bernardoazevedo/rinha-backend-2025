package payment

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adjust/rmq/v5"
	paymentqueue "github.com/bernardoazevedo/rinha-backend-2025/paymentQueue"
)

func PaymentWorker() {
	mainWorkersQuantity := 3
	retryWorkersQuantity := 1

	connection, err := paymentqueue.GetNewConnection()
	if err != nil {
		panic(err)
	}

	// all payments are stored here
	mainQueue, err := connection.OpenQueue("payment")
	if err != nil {
		panic(err)
	}

	// if a consume fail, we move the item to here
	retryQueue, err := connection.OpenQueue("retryPayment")
	if err != nil {
		panic(err)
	}
	mainQueue.SetPushQueue(retryQueue)

	err = mainQueue.StartConsuming(10, time.Second)
	if err != nil {
		panic(err)
	}

	err = retryQueue.StartConsuming(10, time.Second)
	if err != nil {
		panic(err)
	}

	for i := 1; i <= mainWorkersQuantity; i++ {
		_, err = mainQueue.AddConsumer("mainConsumer", NewConsumer("mainConsumer", i))
		if err != nil {
			panic(err)
		}
	}

	for i := 1; i <= retryWorkersQuantity; i++ {
		_, err = retryQueue.AddConsumer("retryConsumer", NewConsumer("retryConsumer", i))
		if err != nil {
			panic(err)
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	defer signal.Stop(signals)

	<-signals
	go func() {
		<-signals
		os.Exit(1)
	}()

	<-connection.StopAllConsuming()
}

func NewConsumer(name string, tag int) *Consumer {
	return &Consumer{
		name:   fmt.Sprintf(name+"%d", tag),
		count:  0,
		before: time.Now(),
	}
}

func (consumer *Consumer) Consume(delivery rmq.Delivery) {
	payment := []byte(delivery.Payload())

	alreadyExistsPayment, err := postPayment(payment)
	consumer.count++
	if alreadyExistsPayment {
		err = delivery.Reject()
		if err != nil {
			fmt.Println("\t\terror acking: " + err.Error())
		}

	} else if err != nil {
		fmt.Println("mandando de " + consumer.name + " pra próxima fila")
		deliveryErr := delivery.Push()
		if deliveryErr != nil {
			fmt.Println("\t\terror pushing: " + deliveryErr.Error())
		}

	} else {
		// success!
		err = delivery.Ack()
		if err != nil {
			fmt.Println("\t\terror acking: " + err.Error())
		}
	}
}
