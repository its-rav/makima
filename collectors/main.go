package main

import (
	"time"

	"github.com/its-rav/makima/pkg/cache"
	"github.com/its-rav/makima/pkg/config"
	"github.com/its-rav/makima/pkg/logger"
	"github.com/its-rav/makima/pkg/model"
	"github.com/its-rav/makima/pkg/twitter"
)

func main() {
	logger.InitLogrusLogger()
	var log = logger.Log

	var config config.CollectorConfig
	config.Load()
	// overrideStreamRules(consumerKey, consumerSecret)

	var getStreamQueryParams twitter.GetStreamQueryParams = twitter.GetStreamQueryParams{
		TweetFields: []string{"created_at", "attachments", "context_annotations", "entities", "public_metrics", "possibly_sensitive", "referenced_tweets", "source", "withheld"},
		Expansions:  []string{"author_id", "attachments.media_keys"},
		UserFields:  []string{"name", "username", "profile_image_url"},
		MediaFields: []string{"url", "preview_image_url", "public_metrics", "alt_text", "variants"},
	}

	bearerToken := twitter.GetBearerToken(config.Twitter.ConsumerKey, config.Twitter.ConsumerSecret)

	twitter.OnStreamReceived(bearerToken, getStreamQueryParams, func(response twitter.TweetResponse) {
		data := response.Data
		log.Infof("[%s] (%s) (%s) New tweet received: %+v", config.ChannelID, data.CreatedAt, time.Now().Format(time.RFC1123), response)

		redisClient := cache.NewClient(config.Redis.ConnString)

		var publishMessage model.PublishMessage[twitter.TweetResponse] = model.PublishMessage[twitter.TweetResponse]{
			Source:      "twitter",
			Destination: "makima:twitter:consumer",
			Message:     response,
			Timestamp:   time.Now(),
		}

		cache.Publish(redisClient, config.ChannelID, publishMessage)

	})
}
