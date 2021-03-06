/*
 * Create By Xinwenjia 2020-02-09
 */

package ailsms

import (
    "fmt"
    "errors"

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
func (c *AliSms) Connect(accessKeyId string, accessKeySecret string, params map[string]string) error {
	var errSms error
    if len(params) == 0 || len(params["regionId"]) == 0 {
        return errors.New(fmt.Sprintf("Params regionId Missing!"))
    }

    c.client, errSms = dysmsapi.NewClientWithAccessKey(params["regionId"], accessKeyId, accessKeySecret)
	return errSms
}

// Send Sms 
func (c *AliSms) SendSms(content map[string]string) (string, string, error) {
    request := dysmsapi.CreateSendSmsRequest()
    request.Scheme = "https"
    request.PhoneNumbers = content["phone"]
    request.SignName = content["sign"]
    request.TemplateCode = content["template"]
    request.TemplateParam = content["param"]

    strResponse := ""
    strRequest := fmt.Sprintln(request)
    response, errSms := c.client.SendSms(request)
    if errSms == nil {
        strResponse = fmt.Sprintln(response)
        if response.Code != "OK" {
            errSms = errors.New(fmt.Sprintf("Gateway Return Failed!"))
        }
    }

    return strRequest, strResponse, errSms
}

