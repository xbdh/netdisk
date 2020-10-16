package mq

import (
	"fmt"
	"log"
	"netdisk/setting"

	"github.com/streadway/amqp"
)

type TransferData struct {
	FileHash     string
	CurLocation  string
	DestLocation string
	//DestStoreType cmn.StoreType
}

var conn *amqp.Connection
var channel *amqp.Channel

// 如果异常关闭，会接收通知
var notifyClose chan *amqp.Error

func Init(cfg *setting.MqConfig) error {
	conn, err := amqp.Dial(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func Publish(exchange, routingKey string, msg []byte) bool {

	if nil == channel.Publish(
		exchange,
		routingKey,
		false, // 如果没有对应的queue, 就会丢弃这条小心
		false, //
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg}) {
		return true
	}
	return false
}

var done chan bool

// StartConsume : 接收消息
func StartConsume(qName, cName string, callback func(msg []byte) bool) {
	msgs, err := channel.Consume(
		qName,
		cName,
		true,  //自动应答
		false, // 非唯一的消费者
		false, // rabbitMQ只能设置为false
		false, // noWait, false表示会阻塞直到有消息过来
		nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	done = make(chan bool)

	go func() {
		// 循环读取channel的数据
		for d := range msgs {
			processErr := callback(d.Body)
			if processErr {
				// TODO: 将任务写入错误队列，待后续处理
			}
		}
	}()

	// 接收done的信号, 没有信息过来则会一直阻塞，避免该函数退出
	<-done

	// 关闭通道
	channel.Close()
}

// StopConsume : 停止监听队列
func StopConsume() {
	done <- true
}
