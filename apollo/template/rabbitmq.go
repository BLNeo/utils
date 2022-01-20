package template

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"sync"
	"time"
)

type RabbitMq struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
}

const RetryTimeInterval = 3

type RabbitChannel struct {
	up            int64
	mu            *sync.Mutex
	conn          *amqp.Connection
	channel       *amqp.Channel
	conf          *RabbitMq
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
}

func (o *RabbitMq) Engine() (*RabbitChannel, error) {
	url := fmt.Sprintf("amqp://%v:%v@%v/", o.User, o.Password, o.Host)
	Log.Info("RabbitMQ Connection : " +url,zap.String("中间件","rabbitmq") )
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	ch := &RabbitChannel{
		conn:          conn,
		up:            time.Now().Unix(),
		mu:            &sync.Mutex{},
		conf:          o,
		channel:       channel,
		connNotify:    conn.NotifyClose(make(chan *amqp.Error)),
		channelNotify: channel.NotifyClose(make(chan *amqp.Error)),
	}
	ch.monitor()
	return ch, nil
}

// 发送
func (o *RabbitChannel) Publish(exchange, message string) error {
	// 推送
	err := o.channel.ExchangeDeclare(
		exchange,
		//订阅模式下为广播类型
		"fanout",
		true,
		false,
		//true表示这个exchange不可以被client用来推送消息,仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//
	err = o.channel.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return err
	}
	return nil
}

// 声明交换机 ,如果交换机存在则创建交换机
func (o *RabbitChannel) declareExchange(exchangeName string) error {
	err := o.channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		Log.Error(fmt.Sprintf("template/rabbitmq: %v",err),zap.String("中间件","rabbitmq"))
		// 交换机不存在的时候
		err := o.channel.ExchangeDeclare(
			exchangeName,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}
		//return err
	}
	return nil
}

// 声明交换机 ,如果交换机存在则创建交换机
func (o *RabbitChannel) declareQueue(queueName string) error {
	_, err := o.channel.QueueDeclare(
		queueName,
		false, // 消息持久化
		false, // 是否自动删除
		false, // 是否独享数据，排他
		false, //
		nil,
	)
	if err != nil {
		Log.Error(fmt.Sprintf("template/rabbitmq: %v",err),zap.String("中间件","rabbitmq"))
		// 交换机不存在的时候
		_, errw := o.channel.QueueDeclare(
			queueName,
			false, // 消息持久化
			false, // 是否自动删除
			false, // 是否独享数据，排他
			false, //
			nil,
		)
		return errw
	}
	return nil
}

// 订阅来自作为消费者订阅来自channel的信息
func (r *RabbitChannel) ReceiveSub(exchangeName, queueName string) (chan []byte, chan struct{}, error) {
	ch := make(chan []byte, 50)
	cl := make(chan struct{})

	//创建交换机 和 队列
	err := r.declareExchange(exchangeName)
	if err != nil {
		return nil, nil, err
	}
	err = r.declareQueue(queueName)
	if err != nil {
		return nil, nil, err
	}

	//3.绑定队列到exchange中，在pub/sub模式下,这里的key要为空
	err = r.channel.QueueBind(queueName, "", exchangeName, false, nil)
	if err != nil {
		return nil, nil, err
	}

	//4.消费消息 监控这个queueName
	messages, err := r.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		for d := range messages {
			ch <- d.Body
		}
		cl <- struct{}{}
	}()

	return ch, cl, nil
}

// 重连机制
func (r *RabbitChannel) monitor() {
	go func() {
		for {
			select {
			case err := <-r.connNotify:
				if err != nil {
					Log.Error(fmt.Sprintf("template/rabbitmq connectionNotify:    %v", err),zap.String("中间件","rabbitmq"))
					r.retryAttachConnection()
				}
			case err := <-r.channelNotify:
				if err != nil {
					Log.Error(fmt.Sprintf("template/rabbitmq channelNotify:  %v", err),zap.String("中间件","rabbitmq"))
					r.retryAttachChannel()
				}

			}
		}
	}()
}

func (r *RabbitChannel) retryAttachChannel() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if time.Now().Unix()-r.up < 60 {
		return
	}

	var (
		err       error
		connected = false
		counter   int64
	)
	for !connected {
		counter++
		Log.Error(fmt.Sprintf("template/rabbitmq: retry create channel times %v ", counter),zap.String("中间件","rabbitmq"))
		r.channel, err = r.conn.Channel()
		if err != nil {
			Log.Error(fmt.Sprintf("template/rabbitmq: create channle is failed  %v", err),zap.String("中间件","rabbitmq"))
			time.Sleep(RetryTimeInterval * time.Second)
			continue
		}
		// success way
		r.channelNotify = r.channel.NotifyClose(make(chan *amqp.Error))
		connected = true
		r.up = time.Now().Unix()
	}
}

func (r *RabbitChannel) retryAttachConnection() {
	var (
		err       error
		connected = false
		counter   int64
	)
	for !connected {
		counter++
		Log.Error(fmt.Sprintf("rabbitmq: retry create connection times %v ", counter),zap.String("中间件","rabbitmq"))
		url := fmt.Sprintf("amqp://%v:%v@%v/", r.conf.User, r.conf.Password, r.conf.Host)
		r.conn, err = amqp.Dial(url)
		if err != nil {
			Log.Error(fmt.Sprintf("rabbitmq: create connection is failed  %v\n", err),zap.String("中间件","rabbitmq"))
			time.Sleep(RetryTimeInterval * time.Second)
			continue
		}
		// success way
		r.channel, err = r.conn.Channel()
		if err != nil {
			r.retryAttachChannel()
		}
		r.connNotify = r.conn.NotifyClose(make(chan *amqp.Error))
		connected = true
	}
}
