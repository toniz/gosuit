/*
 * Create By Xinwenjia 2018-04-25
 */

package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
	"github.com/toniz/gosuit/queue"
)

type RabbitMQ struct {
	conn *amqp.Connection
	chs  map[string]*amqp.Channel
}

func init() {
	queue.Register("rabbitmq", func() queue.MessageQueuer {
		return NewRabbitMQ()
	})
}

func NewRabbitMQ() *RabbitMQ {
	return &RabbitMQ{
		conn: nil,
		chs:  make(map[string]*amqp.Channel),
	}
}

// Set Variables
func (c *RabbitMQ) SetParameter(paramsMap map[string]interface{}) error {

	return nil
}

// Create RabbitMQ Connection Handler, Declare Queue.
func (c *RabbitMQ) Connect(endpoint string, user string, password string) error {
	connString := "amqp://" + user + ":" + password + "@" + endpoint + "/"
	mq, err := amqp.Dial(connString)

	if err != nil {
		log.Printf("Failed to connect amq: %s", connString)
	}

	c.conn = mq
	return err
}

func (c *RabbitMQ) NewChannel(qname string) (*amqp.Channel, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %s", err)
		return ch, err
	}

	_, err = ch.QueueDeclare(
		qname, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Printf("Failed to declare a queue: %s", err)
	}

	c.chs[qname] = ch
	return ch, err
}

func (c *RabbitMQ) GetChannel(qname string) (*amqp.Channel, error) {
	if ch, ok := c.chs[qname]; ok {
		return ch, nil
	} else {
		return c.NewChannel(qname)
	}
}

func (c *RabbitMQ) SendTask(qname string, msg string) error {
	if ch, err := c.GetChannel(qname); err != nil {
		return nil
	} else {
		err = ch.Publish(
			"",    // exchange
			qname, // routing key
			false, // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(msg),
			})
		return err
	}
}

func (c *RabbitMQ) Worker(qname string, fn func([]byte) int) error {
	ch, err := c.conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %s", err)
		return err
	}

	_, err = ch.QueueDeclare(
		qname, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Printf("Failed to declare a queue: %s", err)
		return err
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		log.Printf("Failed to set QoS: %s", err)
		return err
	}

	msgs, err := ch.Consume(
		qname, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %s", err)
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			ret := fn(d.Body)
			if ret == 0 {
				d.Ack(false)
			}

			if ret == 1 {
				break
			}
		}
	}()

	return nil
}
