/*
 * Create By Xinwenjia 2020-02-09
 */

package cosclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMysql(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "COS Client Suite")
}
