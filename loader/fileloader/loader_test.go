/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package fileloader_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/toniz/gosuit/loader"
	_ "github.com/toniz/gosuit/loader/fileloader"
)

var _ = Describe("Loader", func() {
	Describe("Test Json loader.", func() {

		Context("Test Json File List", func() {
			p := "example/"
			It("Should Return Two Json File", func() {
				l, err := loader.NewLoader("file")
				filelist, err := l.GetList(p, ".json", "t*")
				Expect(len(filelist)).To(Equal(2))
				Expect(err).NotTo(HaveOccurred())
			})

		})

		Context("Test Yaml From File", func() {
			type StructA struct {
				A string `yaml:"a"`
			}

			type StructB struct {
				StructA `yaml:",inline"`
				B       string `yaml:"b"`
			}

			p := "example/test.yaml"
			var result StructB
			l, err := loader.NewLoader("file")
			err = l.Load(p, &result)

			It("Should Return Not Errar", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Test Json File List", func() {
			type testJson struct {
				Name  string   `json:"name"`
				Test  int      `json:"test"`
				Array []string `json:"array"`
			}

			p := "example/test.json"
			var result testJson
			l, err := loader.NewLoader("file")
			err = l.Load(p, &result)
			It("Should Return Not Errar", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Test XML From File", func() {
			type Zoo struct {
				Animals []string `xml:"animal"`
			}

			p := "example/xml/test.xml"
			l, err := loader.NewLoader("file")
			var result []Zoo
			err = l.Load(p, &result)

			It("Should Return Not Errar", func() {
				Expect(err).NotTo(HaveOccurred())
			})

		})

	})
})
