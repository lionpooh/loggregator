package httpsetup_test

import (
	"crypto/tls"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"trafficcontroller/httpsetup"
)

var _ = Describe("DefaultClient", func() {
	BeforeEach(func() {
		httpsetup.Setup([]string{
			"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
			"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		})
	})

	It("has a 20 second timeout", func() {
		Expect(http.DefaultClient.Timeout).To(Equal(20 * time.Second))
	})

	Describe("DefaultClient.Transport", func() {
		var transport *http.Transport

		BeforeEach(func() {
			var ok bool
			transport, ok = http.DefaultClient.Transport.(*http.Transport)
			Expect(ok).To(BeTrue(), "Expected http.DefaultClient.Transport to be a *http.Transport")
		})

		It("has a 10 second handshake timeout", func() {
			Expect(transport.TLSHandshakeTimeout).To(Equal(10 * time.Second))
		})

		It("enables DisableKeepAlives", func() {
			Expect(transport.DisableKeepAlives).To(BeTrue())
		})

		It("requires TLS Version 1.2", func() {
			Expect(transport.TLSClientConfig.MinVersion).To(BeEquivalentTo(tls.VersionTLS12))
		})

		It("requires certain cipher suites", func() {
			Expect(transport.TLSClientConfig.CipherSuites).To(ContainElement(tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256))
			Expect(transport.TLSClientConfig.CipherSuites).To(ContainElement(tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384))
		})

		Describe("SetInsecureSkipVerify", func() {
			It("sets InsecureSkipVerify on TLSClientConfig", func() {
				httpsetup.SetInsecureSkipVerify(true)
				Expect(transport.TLSClientConfig.InsecureSkipVerify).To(BeTrue())
			})
		})
	})
})
