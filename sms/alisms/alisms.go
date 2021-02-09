/*
 * Create By Xinwenjia 2020-02-09
 */

package ailsms

import (
	"context"
	"io"

    "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/toniz/gosuit/sms"
)

type AliSms struct {
    client *dysmsapi.Client
}

func init() {
	sms.Register("alisms", func() sms.SmsAgent {
		return new(AliSms)
	})
}

// Create Sms Client Handler
func (c *AliSms) Connect(regionId string, accessKeyID string, secretAccessKey string) error {
	var err error
    c.client, errSms := dysmsapi.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	return errSms
}

// Send Sms 
func (c *AliSms) SendSms(phone string, sign string, templateCode string, templateParam string) (string, string, error) {
    request := dysmsapi.CreateSendSmsRequest()
    request.Scheme = "https"
    request.PhoneNumbers = phone
    request.SignName = sign
    request.TemplateCode = templateCode
    request.TemplateParam = templatePara

    strResponse := ""
    strRequest := fmt.Sprintln(request)
    response, errSms := client.SendSms(request)
    if errSms == nil {
        strResponse = fmt.Sprintln(response)
        if response.Code != "OK" {
            errSms = errors.New(fmt.Sprintf("Gateway Return Failed!"))
        }
    }

    return strRequest, strResponse, errSms
}


