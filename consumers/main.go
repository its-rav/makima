package main

import (
	"encoding/json"
	"fmt"

	"github.com/its-rav/makima/pkg/discord"
	"github.com/its-rav/makima/pkg/redis"
	"github.com/its-rav/makima/pkg/twitter"
)

const (
	TwitterLogo = "https://abs.twimg.com/favicons/twitter.2.ico"
	AppLogo     = "https://static0.gamerantimages.com/wordpress/wp-content/uploads/2022/12/makima-focused-on-gesture.jpg?q=50&fit=contain&w=1140&h=&dpr=1.5"
	AppName     = "Makima"
	ChannelID   = "makima:twitter"
)

func main() {
	//
	config := redis.RedisConfig{
		ConnString: "pubsub-redis:6379",

		Password: "",
		DB:       0,
	}

	redisClient := redis.NewClient(config)

	channel := redis.Subscribe(redisClient, ChannelID)
	redis.Listen(channel, func(message string) {
		// fmt.Println("[%s] Message received: %s",twitterConsumerChannel, message)
		fmt.Printf("[%s] Message received: %s", ChannelID, message)
		fmt.Println()

		// TODO: parse message to twitter.Tweet
		var tweet twitter.Tweet
		_ = json.Unmarshal([]byte(message), &tweet)

		var webhookMessage = discord.DiscordWebhookMessage{
			Embeds: []discord.DiscordWebhookEmbed{
				{
					Title:       ":link:",
					URL:         fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.Username, tweet.TweetID),
					Description: tweet.Content,
					Color:       0x00ACEE,
					Author: discord.DiscordWebhookEmbedAuthor{
						Name:    fmt.Sprintf("%s (@%s)", tweet.Name, tweet.Username),
						URL:     fmt.Sprintf("https://twitter.com/%s", tweet.Username),
						IconURL: tweet.Avatar,
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

		fmt.Printf("[%s] Webhook message sent: %+v", ChannelID, webhookMessage)

		// send message to discord webhook
		discord.SendDiscordWebhookMessage(
			"https://discord.com/api/webhooks//",
			webhookMessage,
		)
	})

}
