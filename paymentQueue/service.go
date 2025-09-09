package paymentqueue

import (
	"os"

	"github.com/adjust/rmq/v5"
)

func GetNewConnection() (rmq.Connection, error) {
	// errChan  := make(chan error, 10)
	connection, err := rmq.OpenConnection("queue", "tcp", os.Getenv("REDIS_URL"), 1, nil)
	if err != nil {
		return connection, err
	}
	return connection, nil
}

func Add(item string) error {
	client, err := GetNewConnection()
	if err != nil {
		return err
	}

	queue, err := client.OpenQueue("payment")
	if err != nil {
		return err
	}

	err = queue.Publish(item)
	if err != nil {
		return err
	}

	return nil
}

