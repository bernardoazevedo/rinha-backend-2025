package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/bernardoazevedo/rinha-de-backend-2025/key"
	"github.com/bernardoazevedo/rinha-de-backend-2025/logger"
	paymentqueue "github.com/bernardoazevedo/rinha-de-backend-2025/paymentQueue"
)

func postPayment(payment Payment) (string, bool, error) {
	var statusCode int

	paymentJson, err := json.Marshal(payment)
	if err != nil {
		return "", false, errors.New("error parsing payment")
	}

	url, _ := key.Get("url")

	postBody := bytes.NewBuffer(paymentJson)

	response, err := http.Post(url+"/payments", "application/json", postBody)
	if response != nil {
		statusCode = response.StatusCode
	} else {
		statusCode = 400
	}

	if err != nil {
		message := fmt.Sprintf("[%d] "+err.Error(), statusCode)
		return "", false, errors.New(message)

	} else if statusCode == http.StatusUnprocessableEntity {
		message := fmt.Sprintf("[%d] payment already exists", statusCode)
		return "", true, errors.New(message)

	} else if statusCode != 200 {
		message := fmt.Sprintf("[%d] status != 200", statusCode)
		return "", false, errors.New(message)

	} else {
		//success
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", false, errors.New("error parsing body")
	}

	return string(responseBody), false, nil
}

func queuePayment(payment Payment) error {
	paymentJson, err := json.Marshal(payment)
	if err != nil {
		return errors.New("error parsing payment")
	}

	err = paymentqueue.Add(string(paymentJson))
	if err != nil {
		return errors.New("error adding to queue")
	}

	return nil
}

func PaymentWorker() {
	logger.Add("starting payment worker...")
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

	_, err = mainQueue.AddConsumer("mainConsumer", NewConsumer("mainConsumer", 1))
	if err != nil {
		panic(err)
	}

	_, err = retryQueue.AddConsumer("retryConsumer", NewConsumer("retryConsumer", 1))
	if err != nil {
		panic(err)
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
	var payment Payment
	var message string

	paymentJson := delivery.Payload()

	err := json.Unmarshal([]byte(paymentJson), &payment)
	if err != nil {
		logger.Add("error parsing payment: " + paymentJson)
		return
	}

	logger.Add(consumer.name + " posting -> " + payment.CorrelationId)

	result, alreadyExistsPayment, err := postPayment(payment)
	consumer.count++
	if alreadyExistsPayment {
		message = "\terror: " + err.Error()

		err = delivery.Ack()
		if err != nil {
			logger.Add("\t\terror acking: " + err.Error())
		}

	} else if err != nil {
		message = "\terror: " + err.Error()

		deliveryErr := delivery.Push()
		if deliveryErr != nil {
			logger.Add("\t\terror pushing: " + deliveryErr.Error())
		}

	} else {
		// success!
		message = "\t" + payment.CorrelationId + ": " + result

		err = delivery.Ack()
		if err != nil {
			logger.Add("\t\terror acking: " + err.Error())
		}

		logger.Add(fmt.Sprintf("\t%s: %d", consumer.name, consumer.count))
	}

	logger.Add(message)
}
