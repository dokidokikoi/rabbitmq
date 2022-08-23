package main

import (
	"fmt"
	mq "go-demo/simple/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rabbitmq := mq.NewRabbitMQSimple("" + "kutou")
	for i := 1; i < 30; i++ {
		rabbitmq.PublishSimple("Hello kutou!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println("发送成功！" + strconv.Itoa(i))
	}

}
