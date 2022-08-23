package main

import "go-demo/topic/rabbitmq"

func main() {
	kutoutwo := rabbitmq.NewRabbitMQTopic("exKutouTopic", "kutou.*.two")
	kutoutwo.ReceiveTopic()
}
