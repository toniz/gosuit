package main

import (
    "encoding/json"
    "fmt"
    "log"

    rb "hsb.com/mq/rabbitmq"
)


func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}



const (
    // RabbitMQ configure
    endpointRB = "10.108.197.89:5672"
    userRB = "user"
    passwordRB = "AMeXcxWc4d"
    queueNameRB = "image_storage"

    // ceph s3 configure
    //endpointCephS3        = "master01:7480"
    //accessKeyIDCephS3     = "42II6092AF4I2OGA5TP9"
    //secretAccessKeyCephS3 = "9MP7VuzkMFVpzIDzL5ueXubdB254RXgezQm5hN3W"

   // aliyun oss configure

)

type Object struct {
    Owner  string
    Bucket string
    Prefix string
    Name   string
}

type Message struct {
    Source Object
    Target Object
}
/*
type Storage interface {
    GetObjectList(bucketName string, objectPrefix string) (names []string, err error)
    PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64) (n int, err error)
    GetObject(ctx context.Context, bucketName string, objectName string) (reader io.Reader, err error)
}

func NewStorage() (Storage, error) {
    switch m {
    case Cash:
        return new(CashPM), nil
    case DebitCard:
        return new(DebitCardPM), nil
    default:
        return nil, errors.New(fmt.Sprintf("Payment method %d not recognized!", m))
    }
}
*/
func MessageHandler(s []byte) error {
    var msg Message
    err := json.Unmarshal(s, &msg)
    failOnError(err, "Failed To Parse RabbitMQ Message")

    fmt.Println(msg)

    return err

}

func main() {
    mq, err := rb.NewRabbitMQ(endpointRB, userRB, passwordRB)
    failOnError(err, "Failed To Connect RabbitMQ")

    //str := `{"Source":{"Owner":"1", "Bucket":"image", "Prefix":"beimian/","Name":"1560431791881-2.jpeg"},"Target":{"Owner":"1", "Bucket":"image", "Prefix":"beimian/","Name":"1560431791881-storage.jpeg"}
    str := `{"Source":{"Owner":"1", "Bucket":"image", "Prefix":"beimian/","Name":"1560431791881-2.jpeg"},"Target":{"Owner":"2", "Bucket":"ibbwhat", "Prefix":"beimian/","Name":"1560431791881-storage.jpeg"}
}`

    err = mq.SendTask(queueNameRB, str)
    failOnError(err, "Failed To Read From RabbitMQ")

}
