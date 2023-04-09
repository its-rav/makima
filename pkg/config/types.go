package config

type RedisConfig struct {
	ConnString string `json:"connectionString, env:"REDIS_CONN_STRING"`
	Password   string `json:"password", env:"REDIS_PASSWORD"`
	DB         int    `json:"db", env:"REDIS_DB"`
}

type TwitterConfig struct {
	ConsumerKey    string `json:"consumerKey", env:"TWITTER_CONSUMER_KEY"`
	ConsumerSecret string `json:"consumerSecret", env:"TWITTER_CONSUMER_SECRET"`
}

type LoggerConfig struct {
	ApiToken string `json:"apiToken", env:"LOGGER_API_TOKEN"`
}

type BaseLoggerConfig struct {
	Logger LoggerConfig `json:"logger"`
}

type ConsumerConfig struct {
	Redis      RedisConfig  `json:"redis"`
	ChannelID  string       `json:"channelId", env:"CHANNEL_ID"`
	WebhookURL string       `json:"webhookUrl", env:"WEBHOOK_URL"`
	Logger     LoggerConfig `json:"logger"`
}

type CollectorConfig struct {
	Redis     RedisConfig   `json:"redis"`
	ChannelID string        `json:"channelId", env:"CHANNEL_ID"`
	Twitter   TwitterConfig `json:"twitter"`
	Logger    LoggerConfig  `json:"logger"`
}
