/*
 * Create By Xinwenjia 2020-02-09
 */

package ossclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMysql(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OSS Client Suite")
}
