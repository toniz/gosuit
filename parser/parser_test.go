/*
 * Create By Xinwenjia 2018-04-15
 * Modify From-https://github.com/toniz/gudp
 */

package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/toniz/gosuit/parser"
)

var _ = Describe("Parser", func() {
	Describe("Test Parser.", func() {
		Context("Test Json Decode", func() {
			type testJson struct {
				Name  string   `json:"name"`
				Test  int      `json:"test"`
				Array []string `json:"array"`
			}

			bytes := []byte(`{"name":"boouoo","test":3,"array":["1","3","4"]}`)
			var j testJson
			err := Decode("json", bytes, &j)

			It("Should Return Json Value", func() {
				Expect(j.Name).To(Equal("boouoo"))
				Expect(j.Test).To(Equal(3))
				Expect(len(j.Array)).To(Equal(3))
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Test XML Decode", func() {
			bytes := []byte(`<animals>
                                <animal>gopher</animal>
                                <animal>armadillo</animal>
                                <animal>zebra</animal>
                                <animal>unknown</animal>
                                <animal>gopher</animal>
                                <animal>bee</animal>
                                <animal>gopher</animal>
                                <animal>zebra</animal>
                            </animals>`)

			var zoo struct {
				Animals []string `xml:"animal"`
			}
			err := Decode("xml", bytes, &zoo)
			It("Should Return Json Value", func() {
				Expect(len(zoo.Animals)).To(Equal(8))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Test Yaml Decode", func() {
			type StructA struct {
				A string `yaml:"a"`
			}

			type StructB struct {
				StructA `yaml:",inline"`
				B       string `yaml:"b"`
			}

			bytes := []byte("a: a string from struct A\nb: a string from struct B")

			var b StructB
			err := Decode("yaml", bytes, &b)

			It("Should Return Json Value", func() {
				Expect(b.A).To(Equal("a string from struct A"))
				Expect(b.B).To(Equal("a string from struct B"))
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
