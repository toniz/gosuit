/*
Package qqsms implements a simple library to send short message using Tencent SMS Service.

The Tencent Doc:
    https://cloud.tencent.com/document/product/382/5976
The Query Example like:
{
    "ext": "",
    "extend": "",
    "params": [
        "验证码",
        "1234",
        "4"
    ],
    "sig": "ee80ad3d76bb1da68387428ca752eb885e52621a3129dcf4d9bc4fd4",
    "sign": "腾讯云",
    "tel": {
        "mobile": "13788888888",
        "nationcode": "86"
    },
    "time": 1457336869,
    "tpl_id": 19
}
*/

package qqsms

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"ibbwhat.com/util/randstr"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
    "errors"

    "github.com/toniz/gosuit/sms"
)

type Tel struct {
	Mobile     string `json:"mobile"`
	Nationcode string `json:"nationcode"`
}

type QQSmsRequest struct {
	Ext    string   `json:"ext"`
	Extend string   `json:"extend"`
	Params []string `json:"params"`
	Sig    string   `json:"sig"`
	Sign   string   `json:"sign"`
	Tel    Tel      `json:"tel"`
	Time   int64    `json:"time"`
	Tpl_id int      `json:"tpl_id"`
}

type QQSms struct {
	AppKey  string
	AppID   string
	request QQSmsRequest
}

func init() {
    sms.Register("qqsms", func() sms.SmsAgent {
        return &QQSms{}
    })
}

// Create Sms Client Handler
func (s *QQSms) Connect(accessKeyId string, accessKeySecret string, params map[string]string) error {

    if len(params) == 0 {
        return errors.New(fmt.Sprintf("Params Missing!"))
    }

    var sign string
    if len(params["sign"]) == 0 {
        sign = "[XXX]"
    }

    var nation string
    if len(params["nation"]) == 0 {
        nation = "86"
    }

    s.AppID = accessKeyId
    s.AppKey = accessKeySecret
	s.request.Sign = sign
	s.request.Tel.Nationcode = nation

    return OK
}

// SendSms: Send SMS Message Using Aliyun SMS Services.
func (s *QQSms) SendSms(phone string, content []string, tpl int) (string, error) {
	s.request.Tel.Mobile = phone
	s.request.Params = content
	s.request.Tpl_id = tpl
	s.request.Time = time.Now().Unix()

	rstr := randstr.RandNum(12)
	url := "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=" + s.AppID + "&random=" + rstr
	s.request.Sig = s.BuildSigStr(rstr, phone, strconv.FormatInt(s.request.Time, 10), s.AppKey)

	b, err := json.Marshal(s.request)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//statuscode := resp.StatusCode
	body, _ := ioutil.ReadAll(resp.Body)

	return fmt.Sprintf("%s", body), err
}

/*
BuildSigStr: Build Sha256 String
Tencent Doc: https://cloud.tencent.com/document/product/382/5976
The Example Like:
    string strMobile = "13788888888"; //tel 的 mobile 字段的内容
    string strAppKey = "5f03a35d00ee52a21327ab048186a2c4"; //sdkappid 对应的 appkey，需要业务方高度保密
    string strRand = "7226249334"; //url 中的 random 字段的值
    string strTime = "1457336869"; //UNIX 时间戳
    string sig = sha256(appkey=5f03a35d00ee52a21327ab048186a2c4&random=7226249334&time=1457336869&mobile=13788888888)
               = ecab4881ee80ad3d76bb1da68387428ca752eb885e52621a3129dcf4d9bc4fd4;
*/
func (s *QQSms) BuildSigStr(rstr, phone, time, appkey string) string {
	shastr := "appkey=" + appkey + "&random=" + rstr + "&time=" + time + "&mobile=" + phone
	h := sha256.New()
	h.Write([]byte(shastr))
	return fmt.Sprintf("%x", h.Sum(nil))
}
