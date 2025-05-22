# NSQ

NSQ is a realtime distributed messaging platform.

## Implementation
This module contains the implementation of Consumer and Producer. You can check the example below.

### Consumer
Consumer is a process that consumes messages from a queue. You can check the example below.
Example:
```go
package main

import (
    "github.com/aidapedia/devkit/mq/nsq/middleware"
    "github.com/aidapedia/devkit/mq/nsq"
)

func main() {
    // new consumer
    consumer, err := nsq.NewConsumer(
        middleware.PanicRecover(),
    )
    if err != nil {
        panic(err)
    }
    err = consumer.AddConsumer(nsq.ConsumerConfig{
        Topic: "topic", // your topic name
        Channel: "channel", // your channel name
        NSQDAddress: []string{"localhost:4150"},
        NSQLookupAddress: []string{"localhost:4161"},
        Concurrent: 1,
        Handler: func(message *nsq.Message) error {
            log.Println("message: ", message.Body)
            return nil
        },
    })
    if err!= nil {
        // got error when adding consumer
        panic(err)
    }

    // start blocking function
    consumer.Start()

    // gracefully stop when done
    consumer.Stop()
}
```

#### Middleware
You can add middleware to the consumer. Middleware is a function that will be executed before the handler. Middleware can be used to implement logging, authentication, etc. You can check the example below.

```go
package main

import (
    "github.com/aidapedia/devkit/mq/nsq/middleware"
    "github.com/aidapedia/devkit/mq/nsq"
)
func main() {
    // ... existing code
    consumer, err := nsq.NewConsumer(
        middleware.PanicRecover(),
        // add your middleware here
    )
    if err!= nil {
        panic(err)
    }
    // ... existing code
}
```

Add your custom middleware like logging or something. 

> IMPORTANT: PanicRecover Middleware MUST BE on the first indices of your middleware initialization to make sure recover function is called when panic occurred.

Example:
```go
package main

import (
    "github.com/aidapedia/devkit/mq/nsq/middleware"
    "github.com/aidapedia/devkit/mq/nsq"
)
func main() {
    //... existing code
    consumer, err := nsq.NewConsumer(
        // add your middleware here
        middleware.PanicRecover(),
        func(topic, channel string, next nsq.HandlerFunc) nsq.HandlerFunc {
		return func(message *nsq.Message) error {
			log.Println("incoming message from topic: ", topic)
			return next(message)
		}
	    },
    )
    if err!= nil {
        panic(err)
    }
}
```

### Producer
Producer is a process that produces messages to a queue. You can check the example below.
Example:
```go
package main
import (
    "github.com/aidapedia/devkit/mq/nsq"
)
func main() {
    // new producer
    producer, err := nsq.NewProducer("localhost:4150")
    if err!= nil {
        panic(err)
    }
    // publish message
    err = producer.Publish("topic", []byte("message"))
    if err!= nil {
        panic(err)
    }
}
```