/*
 * Create By Xinwenjia 2018-05-05
 */

package mqtt_test

import (
    "log"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    . "hsb.com/queue"
    _ "hsb.com/queue/mqtt"
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

    Describe("Test Mqtt API", func() {
        Context("Test Worker", func() {
            queueName := "test/task_queue"
            It("Should Return No Error", func() {
                err := mq.Worker(queueName, func(s []byte) int{
                    log.Printf("Received a message: %s.", s)
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
            })
        })
    })
})
