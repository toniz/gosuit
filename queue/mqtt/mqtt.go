/*
 * Create By Xinwenjia 2018-04-25
 */

package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	. "github.com/eclipse/paho.mqtt.golang"
	"github.com/toniz/gosuit/queue"
)

type Mqtt struct {
	clientID             string
	keepAliveTime        time.Duration
	pingTimeout          time.Duration
	connectRetry         bool
	connectRetryInterval time.Duration
	connectTimeout       time.Duration
	autoReconnect        bool
	maxReconnectInterval time.Duration

	opts               *ClientOptions
	conn               Client
	mattServiceQuality int
}

func init() {
	queue.Register("mqtt", func() queue.MessageQueuer {
		return NewMqtt()
	})
}

func NewMqtt() *Mqtt {
	timestamp := time.Now().Unix()
	return &Mqtt{
		connectRetry:         false,
		connectRetryInterval: 3 * time.Second,
		connectTimeout:       10 * time.Second,
		autoReconnect:        true,
		maxReconnectInterval: 10 * time.Minute,
		keepAliveTime:        10 * time.Second,
		clientID:             fmt.Sprintf("tl_server_%d", timestamp),
		pingTimeout:          1 * time.Second,

		opts:               nil,
		conn:               nil,
		mattServiceQuality: 1,
	}
}

// Set Mqtt Paramters
func (c *Mqtt) SetParameter(paramsMap map[string]interface{}) error {
	if val, ok := paramsMap["connectRetry"]; ok {
		if _, ok = val.(bool); ok {
			c.connectRetry = val.(bool)
		}
	}

	if val, ok := paramsMap["connectRetryInterval"]; ok {
		if _, ok = val.(time.Duration); ok {
			c.connectRetryInterval = val.(time.Duration)
		}
	}

	if val, ok := paramsMap["connectTimeout"]; ok {
		if _, ok = val.(time.Duration); ok {
			c.connectTimeout = val.(time.Duration)
		}
	}

	if val, ok := paramsMap["keepAliveTime"]; ok {
		if _, ok = val.(time.Duration); ok {
			c.keepAliveTime = val.(time.Duration)
		}
	}

	if val, ok := paramsMap["autoReconnect"]; ok {
		if _, ok = val.(bool); ok {
			c.autoReconnect = val.(bool)
		}
	}

	if val, ok := paramsMap["maxReconnectInterval"]; ok {
		if _, ok = val.(time.Duration); ok {
			c.maxReconnectInterval = val.(time.Duration)
		}
	}

	if val, ok := paramsMap["pingTimeout"]; ok {
		if _, ok = val.(time.Duration); ok {
			c.pingTimeout = val.(time.Duration)
		}
	}

	if val, ok := paramsMap["clientID"]; ok {
		if _, ok = val.(string); ok {
			c.clientID = val.(string)
		}
	}

	if val, ok := paramsMap["mattServiceQuality"]; ok {
		if _, ok = val.(int); ok {
			c.mattServiceQuality = val.(int)
		}
	}

	return nil
}

// Create Mqtt Connection Handler, Declare Queue.
func (c *Mqtt) Connect(endpoint string, user string, password string) error {
	c.opts = NewClientOptions().AddBroker(endpoint)
	c.opts.SetClientID(c.clientID)
	c.opts.SetKeepAlive(c.keepAliveTime)
	c.opts.SetPingTimeout(c.pingTimeout)
	c.opts.SetConnectTimeout(c.connectTimeout)
	c.opts.SetConnectRetry(c.connectRetry)
	c.opts.SetConnectRetryInterval(c.connectRetryInterval)
	c.opts.SetAutoReconnect(c.autoReconnect)
	c.opts.SetMaxReconnectInterval(c.connectRetryInterval)

	mq := NewClient(c.opts)
	if token := mq.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Failed to connect mqtt: %v", c)
		return token.Error()
	}
	c.conn = mq
	return nil
}

func (c *Mqtt) SendTask(topic string, msg string) error {
	token := c.conn.Publish(topic, byte(c.mattServiceQuality), false, msg)
	token.Wait()

	return token.Error()
}

func (c *Mqtt) Worker(topic string, fn func([]byte) int) error {
	choke := make(chan [2]string)
	h := func(client Client, msg Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	}

	if token := c.conn.Subscribe(topic, byte(c.mattServiceQuality), h); token.Wait() && token.Error() != nil {
		log.Printf("Failed to Subscribe Message: %v", token.Error())
		return token.Error()
	}

	go func() {
		log.Printf("Start Goroutine: Worker ")
		for {
			incoming := <-choke
			log.Printf("Received Topic: %s Message: %s", incoming[0], incoming[1])
			res, _ := json.Marshal(incoming)
			if ret := fn([]byte(res)); ret == 1 {
				if token := c.conn.Unsubscribe(topic); token.Wait() && token.Error() != nil {
					log.Println(token.Error())
				}
				log.Printf("Callback Function Return Error And Exist: %d", ret)
				break
			}
		}
	}()

	return nil
}
