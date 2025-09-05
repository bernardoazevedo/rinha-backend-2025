package key

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func getNewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
		Protocol: 3,
	})

	return client
}

func Set(key string, value string) error {
	client := getNewClient()
	ctx := context.Background()

	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func Get(key string) (string, error) {
	client := getNewClient()
	ctx := context.Background()

	value, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}
