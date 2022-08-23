package main

import "go-demo/routing/rabbitmq"

func main() {
	kutoutwo := rabbitmq.NewRabbitMQRouting("kutou", "kutou_two")
	kutoutwo.ReceiveRouting()
}
