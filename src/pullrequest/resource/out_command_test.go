package resource_test

import (
	"errors"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	r "pullrequest/resource"
	"pullrequest/resource/fake"
)

var _ = Describe("CheckCommand", func() {
	var fakeSrcDir string

	BeforeEach(func() {
		fakeSrcDir = path.Join(os.TempDir(), "fakedir")
	})

	AfterEach(func() {
		os.Remove(fakeSrcDir)
	})

	Context("when update succeed", func() {
		It("should return correct version", func() {
			fakeGithub := &fake.FGithub{
				UpdatePRResult: "fake-sha1",
				UpdatePRError:  nil,
			}
			outCommand := r.NewOutCommand(fakeGithub)

			outResponse, err := outCommand.Run(fakeSrcDir, r.OutRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(outResponse.Version.Ref).To(Equal("fake-sha1"))
		})
	})

	Context("when update failed", func() {
		It("should return error", func() {
			fakeGithub := &fake.FGithub{
				UpdatePRResult: "",
				UpdatePRError:  errors.New("fake-error"),
			}
			outCommand := r.NewOutCommand(fakeGithub)

			_, err := outCommand.Run(fakeSrcDir, r.OutRequest{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("updating pr: fake-error"))
		})
	})
})
