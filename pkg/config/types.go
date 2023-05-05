package config

type RedisConfig struct {
	ConnString string `json:"connectionString" env:"CONN_STRING"`
	Password   string `json:"password" env:"PASSWORD"`
	DB         int    `json:"db" env:"DB"`
}

type TwitterConfig struct {
	ConsumerKey    string `json:"consumerKey" env:"CONSUMER_KEY"`
	ConsumerSecret string `json:"consumerSecret" env:"CONSUMER_SECRET"`
}

type LoggerConfig struct {
	ApiToken string `json:"apiToken" env:"API_TOKEN"`
}

type BaseLoggerConfig struct {
	Logger LoggerConfig `json:"logger" envPrefix:"LOGGER_"`
}

type ConsumerConfig struct {
	Redis      RedisConfig  `json:"redis" envPrefix:"REDIS_"`
	ChannelID  string       `json:"channelId" env:"CHANNEL_ID"`
	WebhookURL string       `json:"webhookUrl" env:"WEBHOOK_URL"`
	Logger     LoggerConfig `json:"logger" envPrefix:"LOGGER_"`
}

type CollectorConfig struct {
	Redis     RedisConfig   `json:"redis" envPrefix:"REDIS_"`
	ChannelID string        `json:"channelId" env:"CHANNEL_ID"`
	Twitter   TwitterConfig `json:"twitter" envPrefix:"TWITTER_"`
	Logger    LoggerConfig  `json:"logger" envPrefix:"LOGGER_"`
}
