package main

import "go-demo/topic/rabbitmq"

func main() {
	kutouone := rabbitmq.NewRabbitMQTopic("exKutouTopic", "#")
	kutouone.ReceiveTopic()
}
