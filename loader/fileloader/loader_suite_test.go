/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package fileloader_test

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFileLoader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "File Loader Suite")
}

