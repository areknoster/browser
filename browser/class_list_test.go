package browser_test

import (
	. "github.com/stroiman/go-dom/browser"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClassList", func() {
	var (
		el        Element
		classList DOMTokenList
	)
	BeforeEach(func() {
		doc := NewDocument(NewWindow(nil))
		el = doc.CreateElement("div")
		classList = el.ClassList()
	})

	Describe("Add", func() {
		It("Should add a new class", func() {
			Expect(classList.Add("c1", "c2")).To(Succeed())
			Expect(el.GetAttribute("class")).To(Equal("c1 c2"))
		})

		It("Should ignore existing classes", func() {
			el.SetAttribute("class", "c1 c2")
			Expect(classList.Add("c2", "c3")).To(Succeed())
			Expect(el.GetAttribute("class")).To(Equal("c1 c2 c3"))

		})

		It("Should generate a syntax error on empty string", func() {
			err := classList.Add("")
			Expect(IsSyntaxError(err)).To(BeTrue())
		})

		It("Should generate a invalidCharacterError error on string with space", func() {
			err := classList.Add("a b")
			Expect(IsInvalidCharacterError(err)).To(BeTrue())
		})

		It("Should generate an invalidCharacterError on empty string", func() {
			err := classList.Add("a b")
			Expect(IsInvalidCharacterError(err)).To(BeTrue())
		})
	})

	Describe("Contains", func() {
		It("Should return true for an existing class", func() {
			el.SetAttribute("class", "a b c")
			Expect(classList.Contains("a")).To(BeTrue())
		})

		It("Should return false for an non-existing class", func() {
			el.SetAttribute("class", "a b c")
			Expect(classList.Contains("x")).To(BeFalse())
		})
	})

	Describe("Length", func() {
		It("Should return the number of classes", func() {
			el.SetAttribute("class", "a b c")
			Expect(classList.Length()).To(Equal(3))
		})
	})

	Describe("Get/Set Value", func() {
		It("Should read the class attribute", func() {
			el.SetAttribute("class", "a b c")
			Expect(classList.GetValue()).To(Equal("a b c"))
		})

		It("Should write the class attribute", func() {
			classList.SetValue("x y  z")
			Expect(el.GetAttribute("class")).To(Equal("x y  z"))
		})
	})

	Describe("Item", func() {
		It("Should return the item with the specified value", func() {
			el.SetAttribute("class", "a b c")
			Expect(classList.Item(1)).To(HaveValue(Equal("b")))
		})

		It("Should return nil when the index is too large", func() {
			el.SetAttribute("class", "a b c")
			Expect(classList.Item(3)).To(BeNil())
		})
	})
})
