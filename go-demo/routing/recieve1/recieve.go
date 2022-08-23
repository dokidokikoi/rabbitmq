package main

import "go-demo/routing/rabbitmq"

func main() {
	kutouone := rabbitmq.NewRabbitMQRouting("kutou", "kutou_one")
	kutouone.ReceiveRouting()
}
