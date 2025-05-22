# Message Queue

A queue is a linear data structure that follows the First In First Out (FIFO) principle. It is a collection of elements that are inserted and removed according to the order in which they were inserted.

There's some popular message queues like RabbitMQ, Kafka, NSQ, etc.

This module contains the implementation of Consumer and Producer.

## Consumer
The consumer is a process that consumes messages from a queue. Feature of consumer manager is: 
- Implement Middlewares.
- Manage multiple consumers in single client.
 
## Producer
The producer is a process that produces messages to a queue.