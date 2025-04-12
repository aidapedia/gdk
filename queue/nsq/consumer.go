package nsq

import (
	"context"

	nsq "github.com/nsqio/go-nsq"
)

// Config queue config
type ConsumerConfig struct {
	NSQDAdress       []string `json:"nsqd_address"`
	NSQLookupAddress []string `json:"nsqlookup_address"`
	Topic            string   `json:"topic"`
	Channel          string   `json:"channel"`
	Concurrent       int      `json:"concurrentcy"`
	Handler          nsq.HandlerFunc
}

type Consumer struct {
	consumers []*nsq.Consumer
}

// NewConsumer create consumer
func NewConsumer(ctx context.Context, consumers []ConsumerConfig) (*Consumer, error) {
	queue := &Consumer{}
	for _, c := range consumers {
		consumer, err := nsq.NewConsumer(c.Topic, c.Channel, nsq.NewConfig())
		if err != nil {
			return nil, err
		}
		if c.Concurrent == 0 {
			c.Concurrent = 1
		}
		consumer.AddConcurrentHandlers(c.Handler, c.Concurrent)
		err = consumer.ConnectToNSQDs(c.NSQDAdress)
		if err != nil {
			return nil, err
		}
		err = consumer.ConnectToNSQLookupds(c.NSQLookupAddress)
		if err != nil {
			return nil, err
		}
		queue.consumers = append(queue.consumers, consumer)
	}
	return queue, nil
}

func (q *Consumer) Stop() {
	// Stop all consumers gracefully
	for _, c := range q.consumers {
		c.Stop()
	}
	for _, c := range q.consumers {
		<-c.StopChan
	}
}
