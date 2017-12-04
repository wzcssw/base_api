package mq

import (
	"fmt"
	"log"
	"tx_base_api/configs"

	"github.com/streadway/amqp"
)

// RabbitMQ 最高级通道调用
var RabbitMQ RabbitEntity

// RabbitEntity Rabbit结构体实例
type RabbitEntity struct {
	Connection     *amqp.Connection
	Channel        *amqp.Channel
	SendQueueMap   map[string]amqp.Queue
	ReciveQueueMap map[string]amqp.Queue
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		fmt.Printf("%s: %s", msg, err)
	}
}

// InitRabbit 初始化
func InitRabbit() RabbitEntity {

	conn, err := amqp.Dial(configs.AppConfig.Rabbit.URL)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	er := RabbitEntity{
		conn, ch, map[string]amqp.Queue{}, map[string]amqp.Queue{},
	}
	return er
}

// AddSendQueue 增加发送管道
func (er *RabbitEntity) AddSendQueue(queueName string) {
	_, ok := er.SendQueueMap[queueName]
	if !ok {
		item, err := er.Channel.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err == nil {
			er.SendQueueMap[queueName] = item
		}
		failOnError(err, "Failed to add a queue")
	}
}

// GetSendQueue 获得发送管道
func (er *RabbitEntity) GetSendQueue(queueName string) interface{} {
	item, ok := er.SendQueueMap[queueName]
	if ok {
		return item
	}
	return nil
}

// AddExchangeQueue 增加交换机管道（接收消息专用）
func (er *RabbitEntity) AddExchangeQueue(exchangeName string) {
	// 1声明交换路由
	err := er.Channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a exchange")
	// 2声明匿名管道（匿名管道为随时管道名，可以有效防止管道重复）
	queue, err := er.Channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to add a queue")
	// 绑定路由声明
	err = er.Channel.QueueBind(
		queue.Name,
		"",
		exchangeName,
		true,
		nil,
	)
	failOnError(err, "Failed to bind a queue")
	er.ReciveQueueMap[exchangeName] = queue
}

// Close 关闭
func (er *RabbitEntity) Close() {
	fmt.Printf("Rabbit MQ will Close!!!!!!")
	er.Channel.Close()
	er.Connection.Close()
}

// SendMsgToQueue 向管道发送消息
func (er *RabbitEntity) SendMsgToQueue(msg string, queueName string) {
	queue, ok := er.SendQueueMap[queueName]
	if !ok {
		fmt.Printf("查询的管道“%s”不存在", queueName)
		return
	}
	err := er.Channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	log.Println("send ok")
	fmt.Printf("send ok")
	failOnError(err, "Fail to publish a message")
}

// SendMsgToExchange 向交换机发送消息
func (er *RabbitEntity) SendMsgToExchange(msg string, exchangeName string) {
	err := er.Channel.Publish(
		exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	log.Println("send ok")
	fmt.Printf("send ok")
	failOnError(err, "Fail to publish a message")
}

// RecivieMsg 从管道接收消息
func (er *RabbitEntity) RecivieMsg(queueName string, callBack func([]byte)) {
	msgs, err := er.Channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "创建监听管道发送测试消息失败")
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received Queue %s message: %s", queueName, d.Body)
			callBack(d.Body)
		}
	}()
	log.Printf("Waiting for message.To exit press CTRL+C")
	<-forever
}

// RecivieMsgFromExchangeQueue 从交换机管道接收消息
func (er *RabbitEntity) RecivieMsgFromExchangeQueue(exchangeName string, callBack func([]byte)) {
	queue, ok := er.ReciveQueueMap[exchangeName]
	if !ok {
		fmt.Printf("查询的交换机管道“%s”不存在", exchangeName)
		return
	}
	er.RecivieMsg(queue.Name, callBack)
}
