package qqsms_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQqsms(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Qqsms Suite")
}
