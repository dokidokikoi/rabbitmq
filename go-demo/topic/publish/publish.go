package main

import (
	"fmt"
	"go-demo/topic/rabbitmq"
	"strconv"
	"time"
)

func main() {
	kutouone := rabbitmq.NewRabbitMQTopic("exKutouTopic", "kutou.topic.one")
	kutoutwo := rabbitmq.NewRabbitMQTopic("exKutouTopic", "kutou.topic.two")

	for i := 0; i < 100; i++ {
		kutouone.PublishTopic("Hello hutou topic one!" + strconv.Itoa(i))
		kutoutwo.PublishTopic("Hello hutou topic two!" + strconv.Itoa(i))

		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
