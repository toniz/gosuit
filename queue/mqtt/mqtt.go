/*
 * Create By Xinwenjia 2018-04-25
 */

package mqtt

import (
	"encoding/json"
	"fmt"
	"time"
    //"github.com/golang/glog"

	. "github.com/eclipse/paho.mqtt.golang"
	"github.com/toniz/gosuit/queue"
	"github.com/toniz/gosuit/glog"
)

type Msg struct {
    Duplicate bool
    Qos       byte
    Retained  bool
    Topic     string
    MessageID uint16
    Payload   string
}

type Mqtt struct {
	clientID             string
	keepAliveTime        time.Duration
	pingTimeout          time.Duration
	connectRetry         bool
	cleanSession         bool
	connectRetryInterval time.Duration
	connectTimeout       time.Duration
	autoReconnect        bool
	maxReconnectInterval time.Duration
	subscribeTimeout     time.Duration

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
		cleanSession:         true,
		connectRetryInterval: 3 * time.Second,
		connectTimeout:       10 * time.Second,
		autoReconnect:        true,
		maxReconnectInterval: 10 * time.Minute,
		keepAliveTime:        10 * time.Second,
		clientID:             fmt.Sprintf("gosuit_server_%d", timestamp),
		pingTimeout:          1 * time.Second,
        subscribeTimeout:     180 * time.Second,

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

	if val, ok := paramsMap["cleanSession"]; ok {
		if _, ok = val.(bool); ok {
			c.cleanSession = val.(bool)
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

	if val, ok := paramsMap["subscribeTimeout"]; ok {
		if _, ok = val.(int); ok {
			c.subscribeTimeout = val.(time.Duration)
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
	//c.opts.SetConnectRetry(c.connectRetry)
	//c.opts.SetConnectRetryInterval(c.connectRetryInterval)
	c.opts.SetAutoReconnect(c.autoReconnect)
	c.opts.SetMaxReconnectInterval(c.connectRetryInterval)
	c.opts.SetCleanSession(c.cleanSession)

	mq := NewClient(c.opts)
	if token := mq.Connect(); token.Wait() && token.Error() != nil {
		glog.Infof("Failed to connect mqtt: %v", c)
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
	choke := make(chan Message)
	h := func(client Client, msg Message) {
		choke <- msg
	}

	if token := c.conn.Subscribe(topic, byte(c.mattServiceQuality), h); token.Wait() && token.Error() != nil {
		glog.Infof("Failed to Subscribe Message: %v", token.Error())
		return token.Error()
	} else {

	go func() {
		glog.Infof("Start Goroutine: Worker ")
		for {
            select {
                case incoming := <-choke: {
                    glog.Infof("Received Message: %v", incoming)
                    m := Msg{
                        Duplicate: incoming.Duplicate(),
                        Qos: incoming.Qos(),
                        Retained: incoming.Retained(),
                        Topic: incoming.Topic(),
                        MessageID: incoming.MessageID(),
                        Payload: string(incoming.Payload()),
                    }
                    res, _ := json.Marshal(m)
                    if ret := fn([]byte(res)); ret == 1 {
                        if token := c.conn.Unsubscribe(topic); token.Wait() && token.Error() != nil {
                            glog.Infoln(token.Error())
                        }
                        glog.Infof("Callback Function Return Error And Exist: %d", ret)
                        break
                    }
                }
                case <-time.After(c.subscribeTimeout): {
                    glog.Warningf("Subscribe Timeout[%v]: %v", c.subscribeTimeout, token.Error())
                    fn([]byte("Timeout"))
                    break
                }
            }
		}
	}()

    }
	return nil
}

func (c *Mqtt) Close() {
    c.conn.Disconnect(250)
}


