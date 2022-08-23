package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// amqp://用户名:密码@ip:port/Virtual Host
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

//创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	// 创建 RabbitMQ 实例
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error

	// 获取 conncetion
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")

	// 获取 channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open channel")

	return rabbitmq
}

//直接模式队列生产
func (r *RabbitMQ) PublishSimple(message string) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞处理
		false,
		// 额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	// 调用 channel 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞处理
		false,
		// 额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	// 接收消息
	msgs, err := r.channel.Consume(
		q.Name,
		// 用来区分多个消费者
		"", // consumer
		// 是否自动应答
		true, //auto-ack
		// 是否独有
		false, // exclusive
		// 设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		// 列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	// 启动协程处理消息
	go func() {
		for d := range msgs {
			// 消息逻辑处理，可以自行设计逻辑
			fmt.Printf("Received a message: %s\n", string(d.Body))
		}
	}()

	// log.Panicf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
