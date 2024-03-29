package queue

import (
	"github.com/Lawliet-Chan/yu/config"
	"github.com/Lawliet-Chan/yu/storage"
	. "github.com/Lawliet-Chan/yu/yerror"
)

type Queue interface {
	storage.StorageType
	Push(topic string, msg []byte) error
	Pop(topic string) ([]byte, error)

	// The type of msgChan must be chan!
	PushAsync(topic string, msgChan interface{}) error
	PopAsync(topic string, msgChan interface{}) error
}

func NewQueue(cfg *config.QueueConf) (Queue, error) {
	switch cfg.QueueType {
	case "nats":
		return NewNatsQueue(cfg.Url, cfg.Encoder)

	default:
		return nil, NoQueueType
	}
}
