/* 基本类型的设置:
 * SDK采用的是指针风格指定参数，即使对于基本类型你也需要用指针来对参数赋值。
 * SDK提供对基本类型的指针引用封装函数
 * 帮助链接：
 * 短信控制台: https://console.cloud.tencent.com/sms/smslist
 * sms helper: https://cloud.tencent.com/document/product/382/3773 */

package qqsms

import (
    "encoding/json"
    "fmt"
    "strings"
    "errors"

    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
    qqsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"

    "github.com/toniz/gosuit/sms"
)

type QQSms struct {
    client *qqsms.Client
}

func init() {
    sms.Register("qqsms", func() sms.SmsAgent {
        return &QQSms{}
    })
}

// Create Sms Client Handler
func (s *QQSms) Connect(accessKeyId string, accessKeySecret string, params map[string]string) error {
    var err error
    credential := common.NewCredential(
        accessKeyId,
        accessKeySecret,
    )
    cpf := profile.NewClientProfile()
    s.client, err = qqsms.NewClient(credential, "ap-guangzhou", cpf)
    return err
}

// SendSms: Send SMS Message Using Tencent SMS Services.
func (s *QQSms) SendSms(content map[string]string) (string, string, error) {
    request := qqsms.NewSendSmsRequest()
    /* 短信应用ID: 短信SdkAppid在 [短信控制台] 添加应用后生成的实际SdkAppid，示例如1400006666 */
    request.SmsSdkAppid = common.StringPtr(content["appid"])
    /* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名，签名信息可登录 [短信控制台] 查看 */
    request.Sign = common.StringPtr(content["sign"])
    /* 国际/港澳台短信 senderid: 国内短信填空，默认未开通，如需开通请联系 [sms helper] */
    request.SenderId = common.StringPtr("")
    /* 用户的 session 内容: 可以携带用户侧 ID 等上下文信息，server 会原样返回 */
    request.SessionContext = common.StringPtr(content["session"])
    /* 短信码号扩展号: 默认未开通，如需开通请联系 [sms helper] */
    request.ExtendCode = common.StringPtr("0")
    /* 模板 ID: 必须填写已审核通过的模板 ID。模板ID可登录 [短信控制台] 查看 */
    request.TemplateID = common.StringPtr(content["template"])
    /* 下发手机号码，采用 e.164 标准，+[国家或地区码][手机号]
     * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
    if len(content["phone"]) == 0 {
        return "", "", errors.New(fmt.Sprintf("Content Phone Not Exists"))
    }
    phoneList := strings.Split(content["phone"],",")
    for i:=0; i<len(phoneList); i++ {
        phoneList[i] = "+86" + phoneList[i]
    }
    request.PhoneNumberSet = common.StringPtrs(phoneList)

    /* 模板参数: 若无模板参数，则设置为空*/
    if len(content["param"]) == 0 {
        request.TemplateParamSet = common.StringPtrs([]string{"0"})
    } else {
        var paramSet []string
        if jsonErr := json.Unmarshal([]byte(content["param"]), &paramSet); jsonErr != nil {
            return "", "", errors.New(fmt.Sprintf("ParamSet Parse Json Failed"))
        }
        request.TemplateParamSet = common.StringPtrs(paramSet)
    }

    // 通过client对象调用想要访问的接口，需要传入请求对象
    strRequest := fmt.Sprintln(request)
    strResponse := ""
    response, err := s.client.SendSms(request)
    if err == nil {
        res, _ := json.Marshal(response.Response)
        strResponse = string(res)
    }

    return strRequest, strResponse, err
}

