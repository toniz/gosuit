/*
 * Create By Xinwenjia 2020-05-05
 */

package mqtt_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMqtt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Matt Test Suite")
}
