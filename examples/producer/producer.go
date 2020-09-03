package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var (
	uri          = flag.String("uri", "amqp://user01:password01@192.168.0.46:5672/dev", "AMQP URI")
	exchangeName = flag.String("exchange", "xiaoyusan", "Durable AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	routingKey   = flag.String("key", "normal.topic", "AMQP routing key")
	body         = flag.String("body", "foobar", "Body of message")
	reliable     = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")
)

func init() {
	flag.Parse()
}

func main() {
	err := publish(*uri, *exchangeName, *exchangeType, *routingKey, *body, *reliable)
	if err != nil {
		log.Fatalf("%s", err)
	}
	time.Sleep(time.Millisecond * 500)
}

func publish(amqpURI, exchange, exchangeType, routingKey, body string, reliable bool) error {

	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Printf("Dial: %s ", err)
		return err
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s ", err)
	}
	defer channel.Close()

	log.Printf("got Channel, declaring %q Exchange (%q)", exchangeType, exchange)
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

	log.Printf("enabling publishing confirms.")
	err = channel.Confirm(false)
	if err != nil {
		return fmt.Errorf("Channel could not be put into confirm mode: %s ", err)
	}
	confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	defer confirmOne(confirms)

	log.Printf("declared Exchange, publishing %d B body (%q)", len(body), body)
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
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}

}
