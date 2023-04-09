package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// {
// 	"content": null,
// 	"embeds": [
// 	  {
// 		"title": "unusual_whales",
// 		"description": "Discohook has a bot as well, it's not strictly required to send messages it may be helpful to have it ready.\n\nBelow is a small but incomplete overview of what the bot can do for you.",
// 		"color": 5814783,
// 		"fields": [
// 		  {
// 			"name": "Source",
// 			"value": "[:link:](https://google.com)"
// 		  }
// 		],
// 		"author": {
// 		  "name": "@unusual_whales",
// 		  "url": "https://twitter.com/unusual_whales",
// 		  "icon_url": "https://pbs.twimg.com/profile_images/1642939373955035136/pDS3hgcq_400x400.jpg"
// 		},
// 		"footer": {
// 		  "text": "Twitter",
// 		  "icon_url": "https://abs.twimg.com/responsive-web/client-web/icon-ios.b1fc727a.png"
// 		},
// 		"timestamp": "2023-04-05T12:22:00.000Z",
// 		"thumbnail": {
// 		  "url": "https://pbs.twimg.com/profile_images/1642939373955035136/pDS3hgcq_400x400.jpg"
// 		}
// 	  }
// 	],
// 	"attachments": []
//   }

type DiscordWebhookEmbedField struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type DiscordWebhookEmbedAuthor struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type DiscordWebhookEmbedImage struct {
	URL string `json:"url,omitempty"`
}

type DiscordWebhookEmbedFooter struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type DiscordWebhookEmbed struct {
	Title       string                     `json:"title,omitempty"`
	URL         string                     `json:"url,omitempty"`
	Description string                     `json:"description,omitempty"`
	Color       int                        `json:"color,omitempty"`
	Fields      []DiscordWebhookEmbedField `json:"fields,omitempty"`
	Author      DiscordWebhookEmbedAuthor  `json:"author,omitempty"`
	Image       DiscordWebhookEmbedImage   `json:"image,omitempty"`
	Thumbnail   DiscordWebhookEmbedImage   `json:"thumbnail,omitempty"`
	Timestamp   string                     `json:"timestamp,omitempty"`
	Footer      DiscordWebhookEmbedFooter  `json:"footer,omitempty"`
}

type DiscordWebhookMessage struct {
	Content    string                `json:"content,omitempty"`
	Username   string                `json:"username,omitempty"`
	AvatarURL  string                `json:"avatar_url,omitempty"`
	Embeds     []DiscordWebhookEmbed `json:"embeds,omitempty"`
	Flags      int                   `json:"flags,omitempty"`
	ThreadName string                `json:"thread_name,omitempty"`
}

func SendDiscordWebhookMessage(webhookUrl string, message DiscordWebhookMessage) {
	// init http client
	httpClient := &http.Client{}

	// init request
	req, err := http.NewRequest("POST", webhookUrl, nil)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	// set headers
	req.Header.Set("Content-type", "application/json")

	// set body
	data := message

	bodyBytes, err := json.Marshal(data)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	resp, err := httpClient.Do(req)

	if err != nil {
		// print error
		fmt.Println(err)
	}

	defer resp.Body.Close()
}
