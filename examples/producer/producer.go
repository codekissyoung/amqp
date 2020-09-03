package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	log "gitlab.xinhulu.com/platform/GoPlatform/logger"
)

const (
	Uri          = "amqp://user01:password01@192.168.0.46:5672/dev"
	LocalUri     = "amqp://root:root@127.0.0.1:5672/"
	RouteKey     = "normal.topic"
	Queue        = "normal.queue"
	Exchange     = "xiaoyusan"
	ExchangeType = "direct"
)

var (
	uri          = flag.String("uri", LocalUri, "AMQP URI")
	exchangeName = flag.String("exchange", Exchange, "Durable AMQP exchange name")
	exchangeType = flag.String("exchange-type", ExchangeType, "Exchange type - direct|fanout|topic|x-custom")
	routingKey   = flag.String("key", RouteKey, "AMQP routing key")
	body         = flag.String("body", "foobar", "Body of message")
)

func init() {
	flag.Parse()
}

func main() {

	log.SetRequestIdOff()
	log.SetGoroutineIdOn()
	log.SetStdoutOn()
	log.InitLogger("producer")

	err := publish(*uri, *exchangeName, *exchangeType, *routingKey, *body)
	if err != nil {
		log.Info("%s", err)
	}
	time.Sleep(time.Millisecond * 500)
}

func publish(amqpURI, exchange, exchangeType, routingKey, body string) error {

	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Info("Dial: %s ", err)
		return err
	}
	defer connection.Close()
	return nil

	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s ", err)
	}
	defer channel.Close()

	if err := channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s ", err)
	}

	err = channel.Confirm(false)
	if err != nil {
		return fmt.Errorf("Channel could not be put into confirm mode: %s ", err)
	}
	confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	defer confirmOne(confirms)

	if err = channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
			Priority:        0,               // 0-9
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s ", err)
	}

	return nil
}

func confirmOne(confirms <-chan amqp.Confirmation) {

	confirmed := <-confirms

	if confirmed.Ack {
		log.Infof("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Infof("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}

}
