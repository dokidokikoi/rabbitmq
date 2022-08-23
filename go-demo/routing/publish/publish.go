package main

import (
	"fmt"
	"go-demo/routing/rabbitmq"
	"strconv"
	"time"
)

func main() {
	kutouone := rabbitmq.NewRabbitMQRouting("kutou", "kutou_one")
	kutoutwo := rabbitmq.NewRabbitMQRouting("kutou", "kutou_two")

	for i := 0; i < 100; i++ {
		kutouone.PublishRouting("Hello kutou one!" + strconv.Itoa(i))
		kutoutwo.PublishRouting("Hello kutou two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
