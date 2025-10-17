package key

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func GetNewClient() *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
		Protocol: 3,
	})

	return RedisClient
}

func GetClient() *redis.Client {
	ctx := context.Background()
	err := RedisClient.Ping(ctx).Err()
	if err != nil {
		newClient := GetNewClient()
		return newClient
	}
	return RedisClient
}

func Set(key string, value string) error {
	client := GetClient()
	ctx := context.Background()

	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func Get(key string) (string, error) {
	client := GetClient()
	ctx := context.Background()

	value, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}


func Delete(key string) (error) {
	client := GetClient()
	ctx := context.Background()

	_, err := client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func Push(queue string, value string) error {
	client := GetClient()
	ctx := context.Background()

	err := client.RPush(ctx, queue, value).Err()
	if err != nil {
		return err
	}

	return nil
}

func Pop(queue string) (string, error) {
	client := GetClient()
	ctx := context.Background()

	value, err := client.LPop(ctx, queue).Result()
	if err != nil {
		return value, err
	}

	return value, nil
}

func Publish(channel string, value string) error {
	client := GetClient()
	ctx := context.Background()

	err := client.Publish(ctx, channel, value).Err()
	if err != nil {
		return err
	}

	return nil
}