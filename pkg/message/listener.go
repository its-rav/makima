package message

import (
	"github.com/its-rav/makima/pkg/cache"
	"github.com/its-rav/makima/pkg/config"
)

type ListenerConfig struct {
	Channel string
	Redis   config.RedisConfig
}

type Listener[TMessage any] interface {
	Listen()
}

type listener[TMessage any] struct {
	config  ListenerConfig
	handler MessageHandler[TMessage]
	parser  MessageParser[TMessage]
}

func (l *listener[TMessage]) Listen() {
	client := cache.NewClient(l.config.Redis.ConnString)
	channel := cache.Subscribe(client, l.config.Channel)

	cache.Listen(channel, func(rawMessage string) {

		parsed := l.parser.ParseMessage(rawMessage)
		l.handler.HandleMessage(parsed)
	})
}

func (c *ListenerConfig) validate() {
	if c.Channel == "" {
		panic("Channel cannot be empty")
	}
}

func NewListener[TMessage any](config ListenerConfig, handler MessageHandler[TMessage], parser MessageParser[TMessage]) Listener[TMessage] {
	config.validate()

	if parser == nil {
		panic("Parser cannot be null")
	}

	if handler == nil {
		panic("handler cannot be null")
	}

	return &listener[TMessage]{
		config:  config,
		handler: handler,
		parser:  parser,
	}
}
