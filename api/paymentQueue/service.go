package paymentqueue

import (
	"github.com/bernardoazevedo/rinha-backend-2025/api/key"
)

var QueueName string = "payment" 

func Add(item []byte) error {
	err := key.Push(QueueName, string(item))
	if err != nil {
		return err
	}
	return nil
}

func Get() string {
	payment, err := key.Pop(QueueName)
	if err != nil {
		return ""
	}
	return payment
}

func AddToChannel(item []byte) error {
	err := key.Publish(QueueName, string(item))
	if err != nil {
		return err
	}
	return nil
}