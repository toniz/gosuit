package qqsms_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "ibbwhat.com/sms/qqsms"
)

var _ = Describe("Qqsms", func() {
	var (
		s       *QQSms
		err     error
		key     string
		id      string
		phone   string
		content []string
		tpl     int
	)

	BeforeEach(func() {
		key = "1111"
		id = "XXXXX"
		phone = "15989207960"
		content = []string{"hello", "123", "world!"}
		tpl = 1
	})

	Describe("Test TencentYun Sms Request", func() {
		Context("When Get New Sms Client", func() {
			It("Should Return *QQsms", func() {
				s = NewSms(key, id, "", "")
				Expect(s.AppKey).To(Equal(key))
				Expect(s.AppID).To(Equal(id))
			})

			It("Should Return Response", func() {
				_, err = s.SendSms(phone, content, tpl)
				Expect(err).NotTo(HaveOccurred())
			})

			It("Should Return SIG String", func() {
				time := "1548921526"
				sha256 := "92e3478639c79fe810d47e37d916017ad33be73f8a69d619e72ae61850a4b064"
				sig := s.BuildSigStr(key, phone, time, id)
				Expect(sig).To(Equal(sha256))
			})

		})
	})

})
