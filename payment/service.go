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
	"github.com/bernardoazevedo/rinha-de-backend-2025/health"
	"github.com/bernardoazevedo/rinha-de-backend-2025/key"
	"github.com/bernardoazevedo/rinha-de-backend-2025/logger"
	paymentqueue "github.com/bernardoazevedo/rinha-de-backend-2025/paymentQueue"
)

func postPayment(paymentJson string) (string, bool, error) {
	var statusCode int

	url, _ := key.Get("url")

	postBody := bytes.NewBufferString(paymentJson)

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
		errResponse, err := io.ReadAll(response.Body)
		if err != nil {
			return "", false, errors.New("error parsing body")
		}
		message := fmt.Sprintf("[%d] status != 200 - response: "+string(errResponse), statusCode)
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

	// _, err = retryQueue.AddConsumer("retryConsumer", NewConsumer("retryConsumer", 1))
	// if err != nil {
	// 	panic(err)
	// }

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
	// var payment Payment

	paymentJson := delivery.Payload()

	_, alreadyExistsPayment, err := postPayment(paymentJson)
	consumer.count++
	if alreadyExistsPayment {
		err = delivery.Reject()
		if err != nil {
			logger.Add("\t\terror acking: " + err.Error())
		}

	} else if err != nil {
		_, err := health.CheckSetReturnUrl()
		if err != nil {
			logger.Add("" + err.Error())
		}

		deliveryErr := delivery.Push()
		if deliveryErr != nil {
			logger.Add("\t\terror pushing: " + deliveryErr.Error())
		}

	} else {
		// success!
		err = delivery.Ack()
		if err != nil {
			logger.Add("\t\terror acking: " + err.Error())
		}
	}
}
