package paymentqueue

import "github.com/adjust/rmq/v5"

func GetNewConnection() (rmq.Connection, error) {
	connection, err := rmq.OpenConnection("queue", "tcp", "redis:6379", 1, nil)
	if err != nil {
		return connection, err
	}
	return connection, nil
}

func Add(item []byte) error {
	client, err := GetNewConnection()
	if err != nil {
		return err
	}

	queue, err := client.OpenQueue("payment")
	if err != nil {
		return err
	}

	err = queue.PublishBytes(item)
	if err != nil {
		return err
	}

	return nil
}