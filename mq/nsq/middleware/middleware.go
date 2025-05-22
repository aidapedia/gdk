package middleware

import "github.com/nsqio/go-nsq"

type Middleware func(topic, channel string, next nsq.HandlerFunc) nsq.HandlerFunc
