package main

import "go-demo/publish/rabbitmq"

func main() {
	rabbitMQ := rabbitmq.NewRabbitMQPubSub("" + "newProduct")
	rabbitMQ.ReceiveSub()
}
