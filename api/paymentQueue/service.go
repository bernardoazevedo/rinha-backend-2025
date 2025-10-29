package paymentqueue

import (
	"sync"

	"github.com/bernardoazevedo/rinha-backend-2025/api/key"
)

var QueueName string = "payment"
var addLock sync.Mutex
var popLock sync.Mutex

func Add(item []byte) error {
	addLock.Lock()
	err := key.Push(QueueName, string(item))
	addLock.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func Get() string {
	popLock.Lock()
	payment, err := key.Pop(QueueName)
	popLock.Unlock()
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