package main

import (
	"fmt"
	"go-demo/publish/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rabbiMQ := rabbitmq.NewRabbitMQPubSub("" + "newProduct")
	for i := 0; i < 100; i++ {
		rabbiMQ.PublishPub("订阅模式生产第" + strconv.Itoa(i) + "条数据")
		fmt.Println("订阅模式生产第" + strconv.Itoa(i) + "条数据")
		time.Sleep(1 * time.Second)
	}
}
