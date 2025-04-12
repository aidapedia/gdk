package nsq

import (
	nsq "github.com/nsqio/go-nsq"
)

type Producer struct {
	*nsq.Producer
}

func NewProducer(address string) (*Producer, error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(address, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		Producer: producer,
	}, nil
}
