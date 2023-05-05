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

func OnStreamReceived(bearerToken string, params GetStreamQueryParams, callback func(tweet TweetResponse)) {
	// init http client
	httpClient := &http.Client{}

	var queryParams string = convertStructToQueryParams(params)

	// call to https://api.twitter.com/2/tweets/search/stream which is a stream endpoint

	// init request
	url := fmt.Sprintf("https://api.twitter.com/2/tweets/search/stream?%s", queryParams)

	fmt.Printf("url: %s", url)
	fmt.Println()

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

		headers := getResponseHeaders(resp, []string{"x-rate-limit-limit", "x-rate-limit-remaining", "x-rate-limit-reset"})

		// read response body
		dec := json.NewDecoder(resp.Body)
		for {
			var tweetResponse TweetResponse
			err := dec.Decode(&tweetResponse)
			if err != nil {
				if err == io.EOF {
					break
				}
			}

			fmt.Println("-------------------------")
			fmt.Println(tweetResponse)
			fmt.Println("-------------------------")

			callback(tweetResponse)
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
			fmt.Printf("Remaining requests: %d, in time window: %d seconds", remaining, timeWindow)
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
