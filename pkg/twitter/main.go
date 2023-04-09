package twitter

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Tweet struct {
	Content   string `json:"content"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	TweetID   string `json:"tweetId"`
	CreatedAt string `json:"createdAt"`
}

type StreamRule struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type GetStreamRulesResponse struct {
	Data []StreamRule `json:"data"`
	Meta struct {
		Sent        string `json:"sent"`
		ResultCount int    `json:"result_count"`
	} `json:"meta"`
}

type AddStreamRule struct {
	Value string `json:"value"`
	Tag   string `json:"tag,omitempty"`
}

type DeleteStreamRulesAction struct {
	Ids []string `json:"ids"`
}

type AddStreamRulesRequest struct {
	Add []AddStreamRule `json:"add"`
}

type DeleteStreamRulesRequest struct {
	Delete DeleteStreamRulesAction `json:"delete"`
}

type CommandStreamRulesResponse struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Meta  struct {
		Sent    string                 `json:"sent"`
		Summary map[string]interface{} `json:"summary"`
	} `json:"meta"`
}

// tweet.fields: attachments, author_id, context_annotations, conversation_id, created_at, entities, geo, id, in_reply_to_user_id, lang, non_public_metrics, organic_metrics, possibly_sensitive, promoted_metrics, public_metrics, referenced_tweets, reply_settings, source, text, withheld
// user.fields: created_at, description, entities, id, location, name, pinned_tweet_id, profile_image_url, protected, public_metrics, url, username, verified, withheld
// media.fields: duration_ms, height, media_key, preview_image_url, type, url, width, public_metrics, non_public_metrics, organic_metrics, promoted_metrics
// poll.fields: duration_minutes, end_datetime, id, options, voting_status
// place.fields: contained_within, country, country_code, full_name, geo, id, name, place_type
// expansions: author_id, entities.mentions.username, geo.place_id, in_reply_to_user_id, referenced_tweets.id, referenced_tweets.id.author_id

type GetStreamQueryParams struct {
	Expansions      []string `paramName:"expansions"`
	TweetFields     []string `paramName:"tweet.fields"`
	UserFields      []string `paramName:"user.fields"`
	BackfillMinutes int      `paramName:"backfill_minutes"`
	EndTime         string   `paramName:"end_time"`
	StartTime       string   `paramName:"start_time"`
	MediaFields     []string `paramName:"media.fields"`
	PollFields      []string `paramName:"poll.fields"`
	PlaceFields     []string `paramName:"place.fields"`
}

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

func convertStructToQueryParams(params interface{}) string {
	var queryParams []string
	value := reflect.ValueOf(params)
	for i := 0; i < value.NumField(); i++ {
		var paramName string = value.Type().Field(i).Tag.Get("paramName")
		if paramName == "" {
			paramName = value.Type().Field(i).Name
		}

		var paramValue string
		// if string[] then join with comma
		// else just convert to string
		if value.Field(i).Kind() == reflect.Slice {
			paramValue = strings.Join(value.Field(i).Interface().([]string), ",")
		} else {
			paramValue = fmt.Sprintf("%v", value.Field(i).Interface())
		}

		// if paramValue is int and 0 then skip
		if value.Field(i).Kind() == reflect.Int && paramValue == "0" {
			continue
		}

		// if paramValue is empty then skip
		if paramValue == "" {
			continue
		}

		queryParams = append(queryParams, fmt.Sprintf("%s=%s", paramName, paramValue))
	}

	return strings.Join(queryParams, "&")
}

// function to get the bearer token
func GetBearerToken(consumerKey string, consumerSecret string) string {
	keySecretConcat := fmt.Sprintf("%s:%s", consumerKey, consumerSecret)
	b64Encoded := base64.StdEncoding.EncodeToString([]byte(keySecretConcat))

	authURL := "https://api.twitter.com/oauth2/token"
	authHeader := fmt.Sprintf("Basic %s", b64Encoded)

	reqBody := []byte("grant_type=client_credentials")
	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(reqBody))
	if err != nil {
		// handle error
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	authResp := struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
	}{}
	_ = json.Unmarshal(bodyBytes, &authResp)

	return authResp.AccessToken
}

// get stream rules with httpClient and bearerToken
func getStreamRules(httpClient *http.Client, bearerToken string) GetStreamRulesResponse {
	url := "https://api.twitter.com/2/tweets/search/stream/rules"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	req.Header.Set("Content-type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()

	getResponseHeaders(resp, []string{"x-rate-limit-limit", "x-rate-limit-remaining", "x-rate-limit-reset"})

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	var getStreamRulesResponse GetStreamRulesResponse
	_ = json.Unmarshal(bodyBytes, &getStreamRulesResponse)

	return getStreamRulesResponse
}

// get response headers with given keys
func getResponseHeaders(resp *http.Response, keys []string) map[string]string {
	headers := make(map[string]string)
	for _, key := range keys {
		headers[key] = resp.Header.Get(key)
	}

	// print headers
	for key, value := range headers {
		fmt.Printf("%s: %s", key, value)
		fmt.Println()
	}

	return headers
}

func doPostAddStreamRules(httpClient *http.Client, bearerToken string, data AddStreamRulesRequest) CommandStreamRulesResponse {
	url := "https://api.twitter.com/2/tweets/search/stream/rules"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	req.Header.Set("Content-type", "application/json")

	bodyBytes, err := json.Marshal(data)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	fmt.Println(string(bodyBytes))

	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	resp, err := httpClient.Do(req)

	if err != nil {
		// print error
		fmt.Println(err)
	}

	defer resp.Body.Close()

	getResponseHeaders(resp, []string{"x-rate-limit-limit", "x-rate-limit-remaining", "x-rate-limit-reset"})

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	// print response body, status code
	fmt.Println(string(bodyBytes))
	fmt.Println(resp.StatusCode)

	var commandStreamRulesResponse CommandStreamRulesResponse
	_ = json.Unmarshal(bodyBytes, &commandStreamRulesResponse)

	return commandStreamRulesResponse
}

func doPostDeleteStreamRules(httpClient *http.Client, bearerToken string, data DeleteStreamRulesRequest) CommandStreamRulesResponse {
	url := "https://api.twitter.com/2/tweets/search/stream/rules"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	req.Header.Set("Content-type", "application/json")

	bodyBytes, err := json.Marshal(data)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	fmt.Println(string(bodyBytes))

	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	resp, err := httpClient.Do(req)

	if err != nil {
		// print error
		fmt.Println(err)
	}

	defer resp.Body.Close()

	getResponseHeaders(resp, []string{"x-rate-limit-limit", "x-rate-limit-remaining", "x-rate-limit-reset"})

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	// print response body, status code
	fmt.Println(string(bodyBytes))
	fmt.Println(resp.StatusCode)

	var commandStreamRulesResponse CommandStreamRulesResponse
	_ = json.Unmarshal(bodyBytes, &commandStreamRulesResponse)

	return commandStreamRulesResponse
}

// add stream rules
func addStreamRules(httpClient *http.Client, bearerToken string, rules []AddStreamRule) CommandStreamRulesResponse {

	// create CommandStreamRulesRequest
	data := AddStreamRulesRequest{
		Add: rules,
	}

	commandStreamRulesResponse := doPostAddStreamRules(httpClient, bearerToken, data)

	return commandStreamRulesResponse
}

// remove stream rules
func deleteStreamRules(httpClient *http.Client, bearerToken string, ids []string) CommandStreamRulesResponse {
	data := DeleteStreamRulesRequest{
		Delete: DeleteStreamRulesAction{
			Ids: ids,
		},
	}

	commandStreamRulesResponse := doPostDeleteStreamRules(httpClient, bearerToken, data)

	return commandStreamRulesResponse
}

func OnStreamReceived(bearerToken string, params GetStreamQueryParams, callback func(tweet Tweet)) {
	// init http client
	httpClient := &http.Client{}

	var queryParams string = convertStructToQueryParams(params)

	// call to https://api.twitter.com/2/tweets/search/stream which is a stream endpoint

	// init request
	url := fmt.Sprintf("https://api.twitter.com/2/tweets/search/stream?%s", queryParams)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	// set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	req.Header.Set("Content-type", "application/json")

	// loop stream while avoid rate limit
	// perform a number of ( x-rate-limit-remaining ) requests
	// in a given time window ( now - x-rate-limit-reset ) seconds
	// or wait until the time window is over
	for {
		resp, err := httpClient.Do(req)

		if err != nil {
			// print error
			fmt.Println("Err send request")
			fmt.Println(err)
		}
		fmt.Println(resp.StatusCode)

		headers := getResponseHeaders(resp, []string{"x-rate-limit-limit", "x-rate-limit-remaining", "x-rate-limit-reset"})

		// read response body
		dec := json.NewDecoder(resp.Body)
		for {
			var m map[string]interface{}
			err := dec.Decode(&m)
			if err != nil {
				if err == io.EOF {
					break
				}
			}

			fmt.Println(m)

			// print content, username, user, avatar and generate the url to the tweet
			content := m["data"].(map[string]interface{})["text"].(string)
			username := m["includes"].(map[string]interface{})["users"].([]interface{})[0].(map[string]interface{})["username"].(string)
			user := m["includes"].(map[string]interface{})["users"].([]interface{})[0].(map[string]interface{})["name"].(string)
			avatar := m["includes"].(map[string]interface{})["users"].([]interface{})[0].(map[string]interface{})["profile_image_url"].(string)
			tweetId := m["data"].(map[string]interface{})["id"].(string)
			createdAt := m["data"].(map[string]interface{})["created_at"].(string)

			callback(Tweet{
				Content:   content,
				Username:  username,
				Name:      user,
				Avatar:    avatar,
				TweetID:   tweetId,
				CreatedAt: createdAt,
			})
		}

		// print response body, status code
		// fmt.Println(string(bodyBytes))

		//parse remaining requests

		defer resp.Body.Close()
		remaining, _ := strconv.ParseInt(headers["x-rate-limit-remaining"], 10, 64)
		reset, _ := strconv.ParseInt(headers["x-rate-limit-reset"], 10, 64)
		timeWindow := reset - time.Now().Unix()

		// check if rate limit is reached
		if remaining > 0 {
			// display remaining requests in time window
			fmt.Printf("Remaining requests: %d, in time window: %sd seconds", remaining, timeWindow)
			fmt.Println()

			timeWait := time.Duration(timeWindow)*time.Second/time.Duration(remaining) + 500*time.Millisecond
			fmt.Printf("Waiting %s seconds", timeWait)
			fmt.Println()

			// time sleep miliseconds

		} else {
			// get current time
			now := time.Now().Unix()

			// get reset time

			// get time to wait
			wait := reset - now

			// wait until time window is over
			time.Sleep(time.Duration(wait)*time.Second + 1*time.Second)
		}

		// first time only
		break
	}
}

func overrideStreamRules(consumerKey string, consumerSecret string) {

	bearerToken := GetBearerToken(consumerKey, consumerSecret)

	fmt.Println(bearerToken)
	// init http client
	httpClient := &http.Client{}

	// getStreamRules(bearerToken)
	fmt.Println("Getting stream rules")
	rules := getStreamRules(httpClient, bearerToken)
	fmt.Println(rules)
	// print rules line by line
	for _, rule := range rules.Data {
		fmt.Println(rule)
	}

	// get ids of rules to delete
	var ids []string
	for _, rule := range rules.Data {
		ids = append(ids, rule.ID)
	}

	// check if there are rules to delete
	if len(ids) > 0 {
		// delete rules
		fmt.Println("Delete stream rules")
		deleteStreamRules(httpClient, bearerToken, ids)
	}

	// add new rules
	// accept list of string and append to string to string split with space
	var newRulesAdding []string = []string{
		"from:VitalikButerin",
		"from:cz_binance",
		"from:WatcherGuru",
		"from:0xfoobar",
		"from:0xQuit",
		"from:0xCygaar",
		"from:tier10k",
		"from:whale_alert",
		"from:brian_armstrong",
		"from:unusual_whales",
		"from:Tree_of_Alpha",
		"from:elonmusk",
		"from:const_phoenixed",
		"from:hyuktrades",
		"from:News_Of_Alpha",
		"from:GCRClassic",
		"from:CryptoCapo_",
		"from:HsakaTrades",
		"from:AlgodTrading",
	}

	value := strings.Join(newRulesAdding, " OR ")

	var addStreamRulesResponse CommandStreamRulesResponse

	fmt.Println("Add stream rules")
	addStreamRulesResponse = addStreamRules(httpClient, bearerToken, []AddStreamRule{
		{
			Value: value,
		},
	})
	// addStreamRulesResponse to string and print
	fmt.Println(addStreamRulesResponse)

	fmt.Println("Get stream rules")
	newRules := getStreamRules(httpClient, bearerToken)
	fmt.Println(newRules)
	// print rules line by line
	for _, rule := range newRules.Data {
		fmt.Println(rule)
	}

	// exit program
	return
}

// func main() {
// 	consumerKey := "hCVutwtQ7ktNSauVlq1PFAZJc"
// 	consumerSecret := "SdK6qfDt8yTfaNyorOJiluaZzfTYI84R5KoDVt102WDBn9dnY2"
// 	// overrideStreamRules(consumerKey, consumerSecret)

// 	var getStreamQueryParams GetStreamQueryParams = GetStreamQueryParams{
// 		TweetFields: []string{"created_at"},
// 		Expansions:  []string{"author_id"},
// 		UserFields:  []string{"name", "username", "profile_image_url"},
// 	}

// 	bearerToken := GetBearerToken(consumerKey, consumerSecret)

// 	fmt.Println("Loop get stream")
// 	fmt.Println(bearerToken)
// 	loopGetStream(bearerToken, getStreamQueryParams)

// }
