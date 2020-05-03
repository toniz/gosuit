/*
 * Create By Xinwenjia 2018-04-25
 */

package rabbitmq_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRabbitmq(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RabbitMQ Test Suite")
}
