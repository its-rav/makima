package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/its-rav/makima/pkg/model"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewClient(connString string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: connString,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return client
}

func Subscribe(client *redis.Client, channel string) <-chan *redis.Message {
	pubsub := client.Subscribe(ctx, channel)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	return pubsub.Channel()
}

func Publish[T any](client *redis.Client, channel string, message model.PublishMessage[T]) {
	// unmarshall to string
	json, marshalErr := json.Marshal(message)
	if marshalErr != nil {
		panic(marshalErr)
	}
	fmt.Print(string(json))
	err := client.Publish(ctx, channel, string(json)).Err()
	if err != nil {
		panic(err)
	}
}

func Unsubscribe(client *redis.Client, channel string) {
	pubsub := client.Subscribe(ctx, channel)
	err := pubsub.Unsubscribe(ctx, channel)
	if err != nil {
		panic(err)
	}
}

func Listen(channel <-chan *redis.Message, handler func(message string)) {
	for {
		msg := <-channel
		handler(msg.Payload)
	}
}
