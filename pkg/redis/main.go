package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	ConnString string
	Password   string
	DB         int
}

var ctx = context.Background()

// create a new client from the redis connection string
// ping the server to make sure the connection works
// and return the client
func NewClient(config RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: config.ConnString,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return client
}

// subscribe to a channel and listen for messages on it
// return the channel to the caller
func Subscribe(client *redis.Client, channel string) <-chan *redis.Message {
	pubsub := client.Subscribe(ctx, channel)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	return pubsub.Channel()
}

// publish a message to a channel
func Publish(client *redis.Client, channel string, message string) {
	err := client.Publish(ctx, channel, message).Err()
	if err != nil {
		panic(err)
	}
}

// unsubscribe from a channel
func Unsubscribe(client *redis.Client, channel string) {
	pubsub := client.Subscribe(ctx, channel)
	err := pubsub.Unsubscribe(ctx, channel)
	if err != nil {
		panic(err)
	}
}

// while loop to listen for messages on a channel with handler function
func Listen(channel <-chan *redis.Message, handler func(message string)) {
	for {
		msg := <-channel
		handler(msg.Payload)
	}
}

func main() {
	// rdb := NewClient("localhost:6379")
	// channel := Subscribe(rdb, "testa")

	// fmt.Println("Listening for messages on channel 'test'")

	// Listen(channel, func(message string) {
	// 	fmt.Println(message)

	// 	// parse json message
	// 	parsedMessage := struct {
	// 		Channel string `json:"channel"`
	// 		Message string `json:"message"`
	// 	}{}
	// 	_ = json.Unmarshal([]byte(message), &parsedMessage)

	// 	fmt.Println(parsedMessage.Channel)

	// })

	// // unsubscribe from the channel on exit or unexpected error
	// defer Unsubscribe(rdb, "testa")

}
