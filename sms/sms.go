/*
 * Create By Xinwenjia 2020-02-09
 */

package sms

import (
    "errors"
    "fmt"
)

type SmsAgent interface {
    Connect(accessKeyID string, secretAccessKey string, other ...string) error
    SendSms(phone string, sign string, templateCode string, templateParam string) (string, string, error)
}

var (
    namedSmser = make(map[string]func() SmsAgent)
)

func Register(name string, driver func() SmsAgent) {
    namedSmser[name] = driver
}

func NewSmsAgent(name string) (s SmsAgent, err error) {
    if f, ok := namedSmser[name]; ok {
        s = f()
    } else {
        err = errors.New(fmt.Sprintf("Smser type %s not recognized!", name))
    }
    return
}

