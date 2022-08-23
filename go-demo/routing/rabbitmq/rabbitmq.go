package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const MQURL = "amqp://harukaze:123456@192.168.111.132:5672/kutou"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机名称
	Exchange string
	// bind Key 名称
	Key string
	// 连接信息
	Mqurl string
}

func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

// 断开channel 和 connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	// 创建 RabbitMQ 实例
	rabbitMQ := NewRabbitMQ("", exchangeName, routingKey)
	var err error

	// 获取 connection
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.Mqurl)
	rabbitMQ.failOnErr(err, "failed to connect rabbitmq!")

	// 获取 channel
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnErr(err, "failed to open a channel")

	return rabbitMQ
}

func (r *RabbitMQ) PublishRouting(messgae string) {
	// 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机类型要改成direct
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchage")

	// 发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(messgae),
		},
	)
}

func (r *RabbitMQ) ReceiveRouting() {
	// 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failde to declare an exchage")

	// 尝试创建队列
	q, err := r.channel.QueueDeclare(
		// 随机队列名称
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare a queue")

	// 绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)

	// 消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messages {
			log.Printf("recieve a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C")
	<-forever
}
