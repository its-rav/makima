package main

import (
	"encoding/json"
	"time"

	"github.com/its-rav/makima/pkg/config"
	"github.com/its-rav/makima/pkg/logger"
	"github.com/its-rav/makima/pkg/redis"
	"github.com/its-rav/makima/pkg/twitter"
)

func main() {
	logger.InitLogrusLogger()
	var log = logger.Log

	var config config.CollectorConfig
	config.Load()

	// overrideStreamRules(consumerKey, consumerSecret)

	var getStreamQueryParams twitter.GetStreamQueryParams = twitter.GetStreamQueryParams{
		TweetFields: []string{"created_at"},
		Expansions:  []string{"author_id"},
		UserFields:  []string{"name", "username", "profile_image_url"},
		MediaFields: []string{"url", "preview_image_url"},
	}

	bearerToken := twitter.GetBearerToken(config.Twitter.ConsumerKey, config.Twitter.ConsumerSecret)

	twitter.OnStreamReceived(bearerToken, getStreamQueryParams, func(tweet twitter.Tweet) {

		log.Infof("[%s] (%s) (%s) New tweet received: %+v", config.ChannelID, tweet.CreatedAt, time.Now().Format(time.RFC1123), tweet)

		redisClient := redis.NewClient(redis.RedisConfig{
			ConnString: config.Redis.ConnString,
			DB:         config.Redis.DB,
			Password:   config.Redis.Password,
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

		redis.Publish(redisClient, config.ChannelID, string(message))

	})
}
