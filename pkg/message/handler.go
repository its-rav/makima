package message

import (
	"github.com/its-rav/makima/pkg/model"
)

type MessageHandler[TMessage any] interface {
	HandleMessage(message model.PublishMessage[TMessage])
}

type MessageParser[TMessage any] interface {
	ParseMessage(raw string) model.PublishMessage[TMessage]
}
