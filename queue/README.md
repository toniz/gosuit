## 简介
目前只做了任务队列的设计。
目前支持了kafka和rabbitmq  

## 创建对象：
* 创建rabbitmq对象：  
``` go
import (
    . "hsb.com/queue"
    _ "hsb.com/queue/rabbitmq"
)

mq, _ := NewMessageQueue("rabbitmq")
```

* 创建kafka对象：  
``` go
import (
    . "hsb.com/queue"
    _ "hsb.com/queue/kafka"
)
mq, _ := NewMessageQueue("kafka")
```

## 建立链接:
* 链接kafka (kafka可以没有用户名和密码)  
``` go
endpoint := "10.96.16.9:9092"
user := ""
password := ""
err := mq.Connect(endpoint, user, password)
```

*  链接rabbitmq  
``` go
endpoint := "10.111.50.176:5672"
user := "user"
password := "6DuA9eBfLu"
err := mq.Connect(endpoint, user, password)
```

## 发送消息：
* 发送消息 kafka 和rabbitmq 用法一样：  
``` go
queueName := "topic-test"
msg := "Test Send Message"
err := mq.SendTask(queueName, msg)
```
## 接受消息：
* 启用协程接受消息, kafka 和 rabbitmq 用法一样：  
``` go
err := mq.Worker(queueName, func (s []byte) error{
    log.Printf("Received a message: %s", s)
    dot_count := bytes.Count(s, []byte("."))
    t := time.Duration(dot_count)
    time.Sleep(t * time.Second)
    log.Printf("Done")
    return nil
})
```

### 使用例子可以参考
[rabbitmq 测试用例](rabbitmq/rabbitmq_test.go)  
[kafka测试用例](kafka/kafka_test.go)  


