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

// 发布消息
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	rabbitMq := NewRabbitMQ("", exchangeName, "")
	var err error

	// 获取 connection
	rabbitMq.conn, err = amqp.Dial(rabbitMq.Mqurl)
	rabbitMq.failOnErr(err, "failed to connect rabbitmq!")

	// 获取 channel
	rabbitMq.channel, err = rabbitMq.conn.Channel()
	rabbitMq.failOnErr(err, "failed to open a channel")

	return rabbitMq
}

func (r *RabbitMQ) PublishPub(message string) {
	// 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机类型
		"fanout",
		true,
		false,
		// true 表示这个 exchange 不可以被 client 用来推送信息，仅用来进行交换机之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to decleare an exchange")

	// 发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	r.failOnErr(err, "failed to publish message")
}

// 消费
func (r *RabbitMQ) ReceiveSub() {
	// 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机类型
		"fanout",
		true,
		false,
		// true 表示这个 exchange 不可以被 client 用来推送信息，仅用来进行交换机之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to decleare an exchange")

	// 尝试创建队列，注意这里队列名称不要写
	q, err := r.channel.QueueDeclare(
		// 随机生产队列名称
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
		// 在 pub/sub 模式下，这里的key要为空
		"",
		r.Exchange,
		false,
		nil,
	)
	r.failOnErr(err, "failed to bind queue")

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
	r.failOnErr(err, "failed to recieve message")

	forever := make(chan bool)

	go func() {
		for d := range messages {
			log.Printf("received a message: %s", d.Body)
		}
	}()
	fmt.Println("退出请按 CTRL+C")
	<-forever
}
