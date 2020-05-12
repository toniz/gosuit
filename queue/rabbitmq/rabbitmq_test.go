/*
 * Create By Xinwenjia 2018-04-25
 */

package rabbitmq_test

import (
    "time"
    "bytes"
    "log"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    . "hsb.com/queue"
    _ "hsb.com/queue/rabbitmq"
)

var _ = Describe("RabbitMQ Test", func() {

    mq, _ := NewMessageQueue("rabbitmq")
    Context("Test Connection", func() {
        endpoint := "10.111.50.176:5672"
        user := "user"
        password := "6DuA9eBfLu"

        It("Should Return No Error", func() {
            err := mq.Connect(endpoint, user, password)
            Expect(err).NotTo(HaveOccurred())
        })
    })

    Describe("Test RabbitMQ API", func() {
        Context("Test New Task", func() {
            queueName := "task_queue"
            msg := "Test Rabbitmq Task Send"

            It("Should Return No Error", func() {
                err := mq.SendTask(queueName, msg)
                Expect(err).NotTo(HaveOccurred())
            })
        })

        Context("Test Worker", func() {
            queueName := "task_queue"
            It("Should Return No Error", func() {
                err := mq.Worker(queueName, func(s []byte) int{
                    log.Printf("Received a message: %s", s)
                    dot_count := bytes.Count(s, []byte("."))
                    t := time.Duration(dot_count)
                    time.Sleep(t * time.Second)
                    log.Printf("Done")
                    return 0
                })

                Expect(err).NotTo(HaveOccurred())
            })
        })
    })
})
