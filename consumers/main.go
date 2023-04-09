package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/its-rav/makima/pkg/config"
	"github.com/its-rav/makima/pkg/discord"
	"github.com/its-rav/makima/pkg/logger"
	"github.com/its-rav/makima/pkg/redis"
	"github.com/its-rav/makima/pkg/twitter"
)

const (
	TwitterLogo = "https://abs.twimg.com/favicons/twitter.2.ico"
	AppLogo     = "https://static0.gamerantimages.com/wordpress/wp-content/uploads/2022/12/makima-focused-on-gesture.jpg?q=50&fit=contain&w=1140&h=&dpr=1.5"
	AppName     = "Makima"
)

func main() {
	logger.InitLogrusLogger()
	var log = logger.Log

	var config config.ConsumerConfig
	config.Load()

	redisClient := redis.NewClient(redis.RedisConfig{
		ConnString: config.Redis.ConnString,

		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	channel := redis.Subscribe(redisClient, config.ChannelID)

	fmt.Printf("[%s] Starting consumer...", config.ChannelID)
	fmt.Println()

	redis.Listen(channel, func(message string) {
		// fmt.Println("[%s] Message received: %s",twitterConsumerChannel, message)
		log.Infof("[%s] Message received: %s", config.ChannelID, message)

		// TODO: parse message to twitter.Tweet
		var tweet twitter.Tweet
		_ = json.Unmarshal([]byte(message), &tweet)

		var webhookMessage = discord.DiscordWebhookMessage{
			Embeds: []discord.DiscordWebhookEmbed{
				{
					Description: tweet.Content,
					Color:       0x00ACEE,
					Author: discord.DiscordWebhookEmbedAuthor{
						Name:    fmt.Sprintf("%s (@%s)", tweet.Name, tweet.Username),
						URL:     fmt.Sprintf("https://twitter.com/%s", tweet.Username),
						IconURL: tweet.Avatar,
					},
					Fields: []discord.DiscordWebhookEmbedField{
						{
							Name:  "Source",
							Value: fmt.Sprintf("[:link:](https://twitter.com/%s/status/%s)", tweet.Username, tweet.TweetID),
						},
					},
					Footer: discord.DiscordWebhookEmbedFooter{
						Text:    "Twitter",
						IconURL: TwitterLogo,
					},
					Timestamp: tweet.CreatedAt,
				},
			},
			Username:  AppName,
			AvatarURL: AppLogo,
		}

		log.Infof("[%s] (%s) (%s) Webhook message sent: %+v", config.ChannelID, tweet.CreatedAt, time.Now().Format(time.RFC1123), webhookMessage)

		// send message to discord webhook
		discord.SendDiscordWebhookMessage(
			config.WebhookURL,
			webhookMessage,
		)
	})

}
