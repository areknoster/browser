package v8host_test

import (
	"github.com/gost-dom/browser/html"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("window.location", func() {
	It("Should have the location of the document", func() {
		window := html.NewWindow(html.WindowOptionLocation("http://example.com/foo"))
		ctx := host.NewContext(window)
		DeferCleanup(func() {
			ctx.Close()
		})
		Expect(ctx.Eval("location.href")).To(Equal("http://example.com/foo"))
	})
})
