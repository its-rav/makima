package twitter

import "encoding/json"

type User struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	ProfileImageURL string `json:"profile_image_url"`
}

type Domain struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ContextAnnotationEntity struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ContextAnnotation struct {
	Domain Domain                  `json:"domain"`
	Entity ContextAnnotationEntity `json:"entity"`
}

// "display_url": "pic.twitter.com/WxbHcyftdY",
// "end": 199,
// "expanded_url": "https://twitter.com/unusual_whales/status/1649385241071476737/photo/1",
// "media_key": "3_1649238636800647168",
// "start": 176,
// "url": "https://t.co/WxbHcyftdY"

type URLImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type EntityURL struct {
	DisplayURL  string     `json:"display_url"`
	Description string     `json:"description"`
	ExpandedURL string     `json:"expanded_url"`
	Title       string     `json:"title"`
	UnwoundURL  string     `json:"unwound_url"`
	MediaKey    string     `json:"media_key"`
	Images      []URLImage `json:"images"`
	End         int        `json:"end"`   // end index of the URL in the Tweet text
	Start       int        `json:"start"` // start index of the URL in the Tweet text
	URL         string     `json:"url"`   // URL
}

type EntityAnnotation struct {
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Prob           float32 `json:"probability"`
	NormalizedText string  `json:"normalized_text"`
	Type           string  `json:"type"`
}

type CashTag struct {
	Start int    `json:"start"`
	End   int    `json:"end"`
	Tag   string `json:"tag"`
}

type Entity struct {
	Annotations []EntityAnnotation `json:"annotations"`
	Urls        []EntityURL        `json:"urls"`
	CashTags    []CashTag          `json:"cashtags"`
}

type Media struct {
	MediaKey string `json:"media_key"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

type TweetInclude struct {
	Users []User  `json:"users"`
	Media []Media `json:"media"`
}

type TweetData struct {
	PossiblySentitive   bool                `json:"possibly_sensitive"`
	Content             string              `json:"text"`
	TweetID             string              `json:"id"`
	CreatedAt           string              `json:"created_at"`
	EditHistoryTweetIDs []string            `json:"edit_history_tweet_ids"`
	ContextAnnotations  []ContextAnnotation `json:"context_annotations"`
	Entities            Entity              `json:"entities"`
}

type TweetResponse struct {
	Data     TweetData    `json:"data"`
	Includes TweetInclude `json:"includes"`
}

func (resp *TweetResponse) ParseFrom(raw []byte) {
	var r TweetResponse
	// Marshal the raw data into a Tweet struct
	err := json.Unmarshal(raw, &r)
	if err != nil {
		panic(err)
	}

	// Copy the Tweet struct into the receiver
	*resp = r
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
// media.fields: duration_ms, height, media_key, preview_image_url, type, url, width, public_metrics, non_public_metrics, organic_metrics, promoted_metricPrintjson
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
