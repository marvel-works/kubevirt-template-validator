package validation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"

	"github.com/fromanirh/kubevirt-template-validator/pkg/validation"
)

var _ = Describe("Path", func() {
	Context("The JSONPATH filter", func() {
		It("Should detect non-jsonpaths", func() {
			testStrings := []string{
				"string-literal",
				"$.spec.domain.resources.requests.memory",
			}
			for _, s := range testStrings {
				p, err := validation.NewJSONPathFromString(s)
				Expect(p).To(Equal(""))
				Expect(err).To(Equal(validation.ErrInvalidJSONPath))
			}
		})

		It("Should detect non-jsonpaths on creation", func() {
			testStrings := []string{
				"string-literal",
				"$.spec.domain.resources.requests.memory",
			}
			for _, s := range testStrings {
				p, err := validation.NewPath(s)
				Expect(p).To(BeNil())
				Expect(err).To(Equal(validation.ErrInvalidJSONPath))
			}
		})

		It("Should mangle valid JSONPaths", func() {
			expected := "{.spec.template.spec.domain.resources.requests.memory}"
			testStrings := []string{
				"jsonpath::$.spec.domain.resources.requests.memory",
				"jsonpath::.spec.domain.resources.requests.memory",
			}
			for _, s := range testStrings {
				p, err := validation.NewJSONPathFromString(s)
				Expect(p).To(Equal(expected))
				Expect(err).To(BeNil())
			}
		})
	})

	Context("With invalid path", func() {

		var (
			vmCirros *k6tv1.VirtualMachine
		)

		BeforeEach(func() {
			vmCirros = NewVMCirros()
		})

		It("Should return error", func() {
			p, err := validation.NewPath("jsonpath::.spec.this.path.does.not.exist")
			Expect(p).To(Not(BeNil()))
			Expect(err).To(BeNil())

			err = p.Find(vmCirros)
			Expect(err).To(Equal(validation.ErrInvalidJSONPath))
		})

		It("Should detect malformed path", func() {
			p, err := validation.NewPath("jsonpath::random56junk%(*$%&*()")
			Expect(p).To(BeNil())
			Expect(err).To(Not(BeNil()))
		})
	})

	Context("With valid paths", func() {

		var (
			vmCirros *k6tv1.VirtualMachine
		)

		BeforeEach(func() {
			vmCirros = NewVMCirros()
		})

		It("Should provide some integer results", func() {
			s := "jsonpath::.spec.domain.resources.requests.memory"
			p, err := validation.NewPath(s)
			Expect(p).To(Not(BeNil()))
			Expect(err).To(BeNil())

			err = p.Find(vmCirros)
			Expect(err).To(BeNil())
			Expect(p.Len()).To(BeNumerically(">=", 1))

			vals, err := p.AsInt64()
			Expect(err).To(BeNil())
			Expect(len(vals)).To(Equal(1))
			Expect(vals[0]).To(BeNumerically(">", 1024))
		})

		It("Should provide some string results", func() {
			s := "jsonpath::.spec.domain.machine.type"
			p, err := validation.NewPath(s)
			Expect(p).To(Not(BeNil()))
			Expect(err).To(BeNil())

			err = p.Find(vmCirros)
			Expect(err).To(BeNil())
			Expect(p.Len()).To(BeNumerically(">=", 1))

			vals, err := p.AsString()
			Expect(err).To(BeNil())
			Expect(len(vals)).To(Equal(1))
			Expect(vals[0]).To(Equal("q35"))
		})

	})
})
