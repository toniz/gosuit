package qqsms

import (
	"fmt"
	"testing"
)

var s *QQSms
var err error

func Test_NewSms(t *testing.T) {
	key := "11111"
	id := "XXXXX"
	s = NewSms(key, id, "", "")

	if s.appKey != key && s.appID != id {
		t.Error("Tencent SDK NewSms Key[", s.appKey, "] Appid [", s.appID, "] Init Failed!")
	}
}

func Test_SendSms(t *testing.T) {
	phone := "15900000000"
	content := []string{"hello", "123", "world!"}
	tpl := 1
	r, err := s.SendSms(phone, content, tpl)
	if err != nil {
		t.Error("SendSms Failed: ", err)
	}
	fmt.Println("QQyun Return: ", r)
}

func Test_BuildSigStr(t *testing.T) {
	rstr := "1111"
	phone := "15900000000"
	time := "1548921526"
	appkey := "XXXXX"
	sha256 := "92e3478639c79fe810d47e37d916017ad33be73f8a69d619e72ae61850a4b064"

	sig := s.BuildSigStr(rstr, phone, time, appkey)
	fmt.Println(sig)
	if sig != sha256 {
		t.Error("BuildSigStr Not The Same!")
	}
}
