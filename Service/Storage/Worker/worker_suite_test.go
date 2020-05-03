package main_test

import (
	"testing"
        "log"
        "os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


func TestMain(m *testing.M) {
    log.Println("Do stuff BEFORE the tests!")
    exitVal := m.Run()
    log.Println("Do stuff AFTER the tests!")
    os.Exit(exitVal)
}

func TestWorker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Worker Service Suite")
}
