/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package dbproxy_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

)

func TestDBProxy(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "DB Proxy Suite")
}

