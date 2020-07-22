/*
 * Create By Xinwenjia 2018-04-25
 */

package kafka

import (
	"context"
	"log"

	. "github.com/segmentio/kafka-go"
	"github.com/toniz/gosuit/queue"
)

type Kafka struct {
	broker  string
	writers map[string]*Writer
	readers map[string]*Reader
}

func init() {
	queue.Register("kafka", func() queue.MessageQueuer {
		return NewKafka()
	})
}

// Create RabbitMQ Connection Handler, Declare Queue.
func NewKafka() *Kafka {
	return &Kafka{
		broker:  "",
		writers: make(map[string]*Writer),
		readers: make(map[string]*Reader),
	}
}

// Set Variables
func (c *Kafka) SetParameter(paramsMap map[string]interface{}) error {

	return nil
}

func (c *Kafka) Connect(endpoint string, user string, password string) error {
	c.broker = endpoint
	return nil
}

func (c *Kafka) AddWriter(topic string) {
	w := NewWriter(WriterConfig{
		Brokers:  []string{c.broker},
		Topic:    topic,
		Balancer: &Hash{},
	})

	c.writers[topic] = w
	return
}

func (c *Kafka) GetWriter(topic string) *Writer {
	if _, ok := c.writers[topic]; !ok {
		c.AddWriter(topic)
	}
	return c.writers[topic]
}

func (c *Kafka) SendTask(topic string, msg string) error {
	w := c.GetWriter(topic)
	err := w.WriteMessages(context.Background(),
		Message{
			Key:   []byte(msg[0:1]),
			Value: []byte(msg)})

	if err != nil {
		log.Printf("Failed to write message: %s", err)
	}

	return err
}

func (c *Kafka) AddReader(topic string) {
	r := NewReader(ReaderConfig{
		Brokers: []string{c.broker},
		Topic:   topic,
		GroupID: topic, // use topic name because task only need run one time.
		// Partition: 0,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	c.readers[topic] = r
	return
}

func (c *Kafka) GetReader(topic string) *Reader {
	if _, ok := c.readers[topic]; !ok {
		c.AddReader(topic)
	}
	return c.readers[topic]
}

func (c *Kafka) Worker(topic string, fn func([]byte) int) error {
	r := c.GetReader(topic)

	go func() {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Failed To Receive Message: %s", err)
			}

			log.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
			ret := fn(m.Value)
			if ret == 1 {
				break
			}
		}
	}()

	return nil
}

func (c *Kafka) Close() {

}

