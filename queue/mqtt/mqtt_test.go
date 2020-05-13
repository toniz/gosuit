/*
 * Create By Xinwenjia 2018-05-05
 */

package mqtt_test

import (
    "log"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "fmt"

    . "github.com/toniz/gosuit/queue"
    _ "github.com/toniz/gosuit/queue/mqtt"
)

var _ = Describe("Matt Test", func() {

    mq, _ := NewMessageQueue("mqtt")
    Context("Test Connection", func() {
        endpoint := "tcp://127.0.0.1:1883"

        It("Should Return No Error", func() {
            err := mq.Connect(endpoint, "", "")
            Expect(err).NotTo(HaveOccurred())
        })
    })

    choke := make(chan string)
    Describe("Test Mqtt API", func() {
        Context("Test Worker", func() {
            queueName := "test/task_queue"
            It("Should Return No Error", func() {
                err := mq.Worker(queueName, func(s []byte) int{
                    log.Printf("Received a message: %s.", s)
                    choke <- fmt.Sprintln(s)
                    return 0
                })

                Expect(err).NotTo(HaveOccurred())
            })
        })

        Context("Test New Task", func() {
            queueName := "test/task_queue"
            msg := "Test Mqtt Task Send"

            It("Should Return No Error", func() {
                log.Println("Send a message")
                err := mq.SendTask(queueName, msg)
                Expect(err).NotTo(HaveOccurred())
                err = mq.SendTask(queueName, msg)
                Expect(err).NotTo(HaveOccurred())
                receiveCount := 0
                for receiveCount < 2 {
                    incoming := <- choke
                    fmt.Printf("RECEIVED MESSAGE: %s\n", incoming)
                    receiveCount++
                }
            })
        })
    })


})
