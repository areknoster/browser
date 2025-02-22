package v8host_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	// . "github.com/gost-dom/browser/scripting"
)

var _ = Describe("V8 Element", func() {
	It("It should be a direct descendant of Node", func() {
		ctx := NewTestContext()
		Expect(
			ctx.Eval("Object.getPrototypeOf(Element.prototype) === Node.prototype"),
		).To(BeTrue())
	})

	It("Should have nodeType == 1", func() {
		ctx := NewTestContext(LoadHTML(`<div id="1" class="foo"></div>`))
		Expect(ctx.Eval("document.body.nodeType")).To(BeEquivalentTo(1))
	})

	It("Supports textContent property", func() {
		Skip("Add test?")
	})

	Describe("Attributes", func() {
		It("Should support set/getAtribute", func() {
			ctx := NewTestContext(LoadHTML(`<div id="1" class="foo"></div>`))
			Expect(ctx.Eval(
				`document.getElementById("1").getAttribute("class")`,
			)).To(Equal("foo"))
			Expect(ctx.Eval(`
				const e = document.getElementById("1")
				e.setAttribute("data-foo", "bar");
				e.getAttribute("data-foo")`,
			)).To(Equal("bar"))
		})

		It("Should return null when getting non-existing attribute", func() {
			ctx := NewTestContext()
			Expect(
				ctx.Eval(`document.createElement("div").getAttribute("dummy") === null`),
			).To(BeTrue())

		})

		It("Should support hasAtribute", func() {
			ctx := NewTestContext(LoadHTML(`<div id="1" class="foo"></div>`))
			Expect(ctx.Eval(
				`document.getElementById("1").hasAttribute("class")`,
			)).To(BeTrue())
			Expect(ctx.Eval(
				`document.getElementById("1").hasAttribute("foo")`,
			)).To(BeFalse())
		})
	})

	It("Should support insertAdjacentHTML", func() {
		ctx := NewTestContext(LoadHTML(`<div id="1" class="foo"></div>`))
		Expect(ctx.Eval(
			`document.getElementById("1").insertAdjacentHTML("beforebegin", "<p>foo</p>")`,
		)).Error().ToNot(HaveOccurred())
		Expect(
			ctx.Window().Document().Body().OuterHTML(),
		).To(Equal(`<body><p>foo</p><div id="1" class="foo"></div></body>`))
	})

	It("Should have a querySelector function", func() {
		ctx := NewTestContext(LoadHTML(`<div id="1" class="foo"></div>`))
		Expect(ctx.Eval(
			`typeof document.getElementById("1").querySelector`,
		)).To(Equal("function"))
	})

	It("Should have a querySelectorAll function", func() {
		ctx := NewTestContext(LoadHTML(`<div id="1" class="foo"></div>`))
		Expect(ctx.Eval(
			`typeof document.getElementById("1").querySelectorAll`,
		)).To(Equal("function"))
	})
})
