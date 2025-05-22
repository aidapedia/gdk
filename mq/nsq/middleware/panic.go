package middleware

import (
	"log"

	"github.com/nsqio/go-nsq"
)

// PanicRecover panic recover middleware for nsq.
func PanicRecover() Middleware {
	return func(topic, channel string, next nsq.HandlerFunc) nsq.HandlerFunc {
		return func(message *nsq.Message) error {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic recover: %v", err)
				}
			}()
			return next(message)
		}
	}
}
