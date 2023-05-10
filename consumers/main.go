package main

import (
	"encoding/json"
	"fmt"
	"time"

	conf "github.com/its-rav/makima/pkg/config"
	"github.com/its-rav/makima/pkg/discord"
	logger "github.com/its-rav/makima/pkg/logger"
	"github.com/its-rav/makima/pkg/message"
	"github.com/its-rav/makima/pkg/model"
	"github.com/its-rav/makima/pkg/twitter"
)

const (
	TwitterLogo = "https://abs.twimg.com/favicons/twitter.2.ico"
	AppLogo     = "https://static0.gamerantimages.com/wordpress/wp-content/uploads/2022/12/makima-focused-on-gesture.jpg?q=50&fit=contain&w=1140&h=&dpr=1.5"
	AppName     = "Makima"
)

var log logger.Logger
var config conf.ConsumerConfig

type TweetHandler[TMessage twitter.TweetResponse] struct{}

type TweetParser[TMessage twitter.TweetResponse] struct{}

func (h *TweetHandler[TMessage]) HandleMessage(message model.PublishMessage[twitter.TweetResponse]) {
	tweetResponse := message.Message
	data := tweetResponse.Data

	//  check if tweetResponse.Includes not null
	var users string
	if tweetResponse.Includes.Users != nil && len(tweetResponse.Includes.Users) > 1 {
		for _, user := range tweetResponse.Includes.Users[1:] {
			users += fmt.Sprintf("[%s](https://twitter.com/%s)\n ", user.Name, user.Username)
		}
	}

	var medias []string
	for _, media := range tweetResponse.Includes.Media {
		medias = append(medias, media.URL)
	}

	// entities string seperated by comma
	var entityUrls string
	// only urls from 1 index

	for _, entity := range data.Entities.Urls {
		val := fmt.Sprintf("[%s](%s)\n ", entity.DisplayURL, entity.ExpandedURL)
		entityUrls += val
		fmt.Println(val)
	}
	// remove last comma
	if len(entityUrls) > 0 {
		entityUrls = entityUrls[:len(entityUrls)-2]
	}

	// context annotations string seperated by comma
	var contextAnnotationsDict map[string][]string
	contextAnnotationsDict = make(map[string][]string)
	for _, contextAnnotation := range data.ContextAnnotations {
		fmt.Println(contextAnnotation)
		val := contextAnnotation.Entity.Name
		if contextAnnotationsDict[contextAnnotation.Domain.Name] == nil {
			contextAnnotationsDict[contextAnnotation.Domain.Name] = []string{
				val,
			}
		} else {
			contextAnnotationsDict[contextAnnotation.Domain.Name] = append(contextAnnotationsDict[contextAnnotation.Domain.Name], val)
		}
	}

	contextAnnotationsStr := ""
	// check null
	if contextAnnotationsDict != nil && len(contextAnnotationsDict) > 0 {
		for domain, rawItems := range contextAnnotationsDict {
			// join items by comma replace last comma with \n
			itemsStr := ""
			for _, item := range rawItems {
				itemsStr += fmt.Sprintf("%s, ", item)
			}
			itemsStr = itemsStr[:len(itemsStr)-2]

			contextAnnotationsStr += fmt.Sprintf("**%s**: %s\n", domain, itemsStr)
		}
	}

	// group entity annotations by type then format it into entityAnnotation.Type: entityAnnotation.NormalizedText (entityAnnotation.Probability) entityAnnotation2.NormalizedText (entityAnnotation2.Probability)
	var entityAnnotationsDict map[string][]string
	entityAnnotationsDict = make(map[string][]string)
	for _, entityAnnotation := range data.Entities.Annotations {
		fmt.Println(entityAnnotation)
		percent := int(entityAnnotation.Prob * 100)
		val := fmt.Sprintf("%s (%d%%)", entityAnnotation.NormalizedText, percent)
		if entityAnnotationsDict[entityAnnotation.Type] == nil {
			entityAnnotationsDict[entityAnnotation.Type] = []string{
				val,
			}
		} else {

			entityAnnotationsDict[entityAnnotation.Type] = append(entityAnnotationsDict[entityAnnotation.Type], val)
		}
	}

	entityAnnotationsStr := ""
	if entityAnnotationsDict != nil && len(entityAnnotationsDict) > 0 {
		for t, rawItems := range entityAnnotationsDict {
			// join items by comma replace last comma with \n
			itemsStr := ""
			for _, item := range rawItems {
				itemsStr += fmt.Sprintf("%s, ", item)
			}
			itemsStr = itemsStr[:len(itemsStr)-2]

			entityAnnotationsStr += fmt.Sprintf("**%s**: %s\n", t, itemsStr)
		}
	}

	author := tweetResponse.Includes.Users[0]
	var fields []discord.DiscordWebhookEmbedField = []discord.DiscordWebhookEmbedField{
		{
			Name:   "Source",
			Value:  fmt.Sprintf("[tweet](https://twitter.com/%s/status/%s)", author.Username, data.TweetID),
			Inline: true,
		},
	}

	if users != "" {
		fields = append(fields, discord.DiscordWebhookEmbedField{
			Name:   "Mentions",
			Value:  users,
			Inline: true,
		})
	}

	if entityUrls != "" {
		fields = append(fields, discord.DiscordWebhookEmbedField{
			Name:   "Urls",
			Value:  entityUrls,
			Inline: true,
		})
	}

	if contextAnnotationsStr != "" {
		fields = append(fields, discord.DiscordWebhookEmbedField{
			Name:   "Context",
			Value:  contextAnnotationsStr,
			Inline: true,
		})
	}

	if entityAnnotationsStr != "" {
		fields = append(fields, discord.DiscordWebhookEmbedField{
			Name:   "Annotations",
			Value:  entityAnnotationsStr,
			Inline: true,
		})
	}

	var webhookMessage = discord.DiscordWebhookMessage{
		Embeds: []discord.DiscordWebhookEmbed{
			{
				Description: data.Content,
				Color:       0x00ACEE,
				Author: discord.DiscordWebhookEmbedAuthor{
					Name:    fmt.Sprintf("%s (@%s)", author.Name, author.Username),
					URL:     fmt.Sprintf("https://twitter.com/%s", author.Username),
					IconURL: author.ProfileImageURL,
				},
				Fields: fields,
				Footer: discord.DiscordWebhookEmbedFooter{
					Text:    "Twitter",
					IconURL: TwitterLogo,
				},
				Timestamp: data.CreatedAt,
			},
		},

		Username:  AppName,
		AvatarURL: AppLogo,
	}

	if medias != nil && len(medias) > 0 {

		webhookMessage.Embeds[0].Thumbnail = discord.DiscordWebhookEmbedImage{
			URL: medias[0],
		}

		if len(medias) > 1 {
			webhookMessage.Embeds[0].Image = discord.DiscordWebhookEmbedImage{
				URL: medias[1],
			}
		}
	}

	log.Infof("[%s] (%s) (%s) Webhook message sent: %+v", config.ChannelID, data.CreatedAt, time.Now().Format(time.RFC1123), webhookMessage)

	// send message to discord webhook
	discord.SendDiscordWebhookMessage(
		config.WebhookURL,
		webhookMessage,
	)
}

func (p *TweetParser[TMessage]) ParseMessage(raw string) model.PublishMessage[twitter.TweetResponse] {
	var message model.PublishMessage[twitter.TweetResponse]

	err := json.Unmarshal([]byte(raw), &message)
	if err != nil {
		log.Errorf(err, "[%s] Error while parsing message.", config.ChannelID)
	}

	return message
}

func main() {

	logger.InitLogrusLogger()
	log = logger.Log

	config.Load()

	fmt.Printf("[%s] Starting consumer...", config.ChannelID)

	l := message.NewListener[twitter.TweetResponse](
		message.ListenerConfig{
			Redis:   config.Redis,
			Channel: config.ChannelID,
		},
		&TweetHandler[twitter.TweetResponse]{},
		&TweetParser[twitter.TweetResponse]{},
	)

	l.Listen()
}
