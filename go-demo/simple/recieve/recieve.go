package main

import mq "go-demo/simple/rabbitmq"

func main() {
	rabbitmq := mq.NewRabbitMQSimple("" + "kutou")
	rabbitmq.ConsumeSimple()
}
