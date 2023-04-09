package main

import (
	"encoding/json"
	"fmt"

	"github.com/its-rav/makima/pkg/redis"
	"github.com/its-rav/makima/pkg/twitter"
)

const (
	CollectorID = "makima:twitter:collector"
	ChannelID   = "makima:twitter"
)

func main() {
	consumerKey := ""
	consumerSecret := ""
	// overrideStreamRules(consumerKey, consumerSecret)

	var getStreamQueryParams twitter.GetStreamQueryParams = twitter.GetStreamQueryParams{
		TweetFields: []string{"created_at"},
		Expansions:  []string{"author_id"},
		UserFields:  []string{"name", "username", "profile_image_url"},
	}

	bearerToken := twitter.GetBearerToken(consumerKey, consumerSecret)

	fmt.Printf("[%s] Starting collector...", CollectorID)

	twitter.OnStreamReceived(bearerToken, getStreamQueryParams, func(tweet twitter.Tweet) {
		redisClient := redis.NewClient(redis.RedisConfig{
			ConnString: "pubsub-redis:6379",
			DB:         0,
			Password:   "",
		})

		// var publishMessage types.PublishMessage = types.PublishMessage{
		// 	Channel: "makima:twitter:consumer",
		// 	Message: json.Marshal(tweet ),
		// 	Timestamp: time.Now(),
		// }

		message, err := json.Marshal(tweet)
		if err != nil {
			panic(err)
		}

		redis.Publish(redisClient, ChannelID, string(message))

	})
}
