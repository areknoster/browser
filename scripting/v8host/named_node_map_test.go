package v8host_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("V8 NamedNodeMap", func() {
	ctx := InitializeContextWithEmptyHtml()

	It("Should inherit directly from Object", func() {
		Expect(
			ctx.Eval("Object.getPrototypeOf(NamedNodeMap.prototype) === Object.prototype"),
		).To(BeTrue())
	})

	It("Should allow iterating attributes", func() {
		ctx := NewTestContext(LoadHTML(`<body><div id="foo" class="bar" hidden></div></body>`))
		Expect(ctx.Eval(`
const elm = document.getElementById("foo");
const attributes = elm.attributes;
let idAttribute;
for (let i = 0; i < attributes.length; i++) {
  const attr = attributes.item(i)
  if (attr.name === "id") { idAttribute = attr }
}`)).Error().ToNot(HaveOccurred())
		Expect(ctx.Eval("attributes.length")).To(BeEquivalentTo(3))
		Expect(ctx.Eval("idAttribute.value")).To(Equal("foo"))
		ctx.MustRunTestScript("idAttribute.value = 'bar'")
		Expect(ctx.Eval("idAttribute.value")).To(Equal("bar"))
	})

	It("Should allow indexing by number", func() {
		ctx := NewTestContext(LoadHTML(`<body><div id="foo" class="bar" hidden></div></body>`))
		Expect(ctx.Eval(`
const elm = document.getElementById("foo");
const attributes = elm.attributes;
attributes[0] instanceof Attr &&
attributes[1] instanceof Attr &&
attributes[2] instanceof Attr
`)).To(BeTrue())
	})

	It("Should return `null` when indexing outside the elements", func() {
		ctx := NewTestContext(LoadHTML(`<body><div id="foo" class="bar" hidden></div></body>`))
		Expect(ctx.Eval(`
const elm = document.getElementById("foo");
const attributes = elm.attributes;
attributes[3]
`)).To(BeNil())
	})

	It("Should have nodeType 2 on attributes", func() {
		ctx := NewTestContext(LoadHTML(`<body><div id="foo" class="bar" hidden></div></body>`))
		Expect(ctx.Eval(`
const elm = document.getElementById("foo");
const attribute = elm.attributes.item(0);
attribute.nodeType
`)).To(BeEquivalentTo(2))
	})

	Describe("Retrieve bad index", func() {
		It("Should return undefined when using attributes.[index]", func() {
			Expect(ctx.Eval("document.body.attributes[42] === undefined")).To(BeTrue())
		})

		It("Should return null when using attributes.item(index)", func() {
			Expect(ctx.Eval("document.body.attributes.item(42) === null")).To(BeTrue())
		})
	})
})
