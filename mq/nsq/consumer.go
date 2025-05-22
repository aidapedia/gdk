package nsq

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aidapedia/gdk/mq/nsq/middleware"
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
	consumers  []*nsq.Consumer
	middleware []middleware.Middleware
}

// NewConsumer create consumer
func NewConsumer(middlewares ...middleware.Middleware) (*Consumer, error) {
	return &Consumer{
		consumers:  []*nsq.Consumer{},
		middleware: middlewares,
	}, nil
}

func (q *Consumer) AddConsumer(consumers []ConsumerConfig) error {
	for _, c := range consumers {
		consumer, err := nsq.NewConsumer(c.Topic, c.Channel, nsq.NewConfig())
		if err != nil {
			return err
		}
		if c.Concurrent == 0 {
			c.Concurrent = 1
		}
		// Initialize the middlewares
		for _, m := range q.middleware {
			c.Handler = m(c.Topic, c.Channel, c.Handler)
		}
		// Add the handler
		consumer.AddConcurrentHandlers(c.Handler, c.Concurrent)
		err = consumer.ConnectToNSQDs(c.NSQDAdress)
		if err != nil {
			return err
		}
		err = consumer.ConnectToNSQLookupds(c.NSQLookupAddress)
		if err != nil {
			return err
		}
		q.consumers = append(q.consumers, consumer)
	}
	return nil
}

func (q *Consumer) Start() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
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
