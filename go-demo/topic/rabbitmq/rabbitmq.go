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

func NewRabbitMQTopic(exchageName string, routingKey string) *RabbitMQ {
	rabbitMQ := NewRabbitMQ("", exchageName, routingKey)
	var err error

	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.Mqurl)
	rabbitMQ.failOnErr(err, "failed to connect rabbitmq!")

	// 获取 channel
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnErr(err, "failed to open a channel")

	return rabbitMQ
}

func (r *RabbitMQ) PublishTopic(message string) {
	// 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")

	// 发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// "*"匹配一个单词，"#"匹配多个单词（可以是0个）
// kutou.hello "kutou.*"可以匹配到
// kutou.hello.world "kutou.#"可以匹配到
func (r *RabbitMQ) ReceiveTopic() {
	// 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")

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
		r.Key,
		r.Exchange,
		false,
		nil,
	)

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
			log.Printf("recieved a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C")
	<-forever
}
