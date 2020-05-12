/*
 * Create By Xinwenjia 2018-04-25
 */

package queue

import (
    "errors"
    "fmt"
)

// Worker Definition: fn function return: 0 continue; 1 exit; 
type MessageQueuer interface {
    SetParameter(paramsMap map[string]interface{}) error
    Connect(endpoint string, user string, password string) error
    Worker(qname string, fn func([]byte) int) error
    SendTask(qname string, msg string) error
}

var (
    mqers = make(map[string]func() MessageQueuer)
)

func Register(name string, l func() MessageQueuer) {
    mqers[name] = l
}


// 't' has three values to choose from: rabbitmq, kafka, matt
func NewMessageQueue(t string) (m MessageQueuer, err error) {
    if f, ok := mqers[t]; ok {
        m = f()
    } else {
        err = errors.New(fmt.Sprintf("Loader type %s not recognized! [rabbitmq, kafka]", t))
    }
    return
}

