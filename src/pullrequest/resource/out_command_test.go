package resource_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	r "pullrequest/resource"
	"pullrequest/resource/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckCommand", func() {
	var fakeSrcDir string
	var err error

	Context("when pr_number is there", func() {
		BeforeEach(func() {
			fakeSrcDir = path.Join(os.TempDir(), "fakedir")
			err = os.Mkdir(fakeSrcDir, 0777)
			Expect(err).ToNot(HaveOccurred())

			err = ioutil.WriteFile(path.Join(fakeSrcDir, "pr_number"), []byte("1"), 0777)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err = os.RemoveAll(fakeSrcDir)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when update succeed", func() {
			It("should return correct version", func() {
				fakeGithub := &fake.FGithub{
					UpdatePRResult: "fake-ref1",
					UpdatePRError:  nil,
				}
				outCommand := r.NewOutCommand(fakeGithub)

				outResponse, err := outCommand.Run(fakeSrcDir, r.OutRequest{})
				Expect(err).ToNot(HaveOccurred())
				Expect(outResponse.Version).To(Equal(r.Version{Ref: "fake-ref1", PR: "1"}))
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

	Context("when pr_number is not there", func() {
		BeforeEach(func() {
			fakeSrcDir = path.Join(os.TempDir(), "fakedir")
			err = os.Mkdir(fakeSrcDir, 0777)
			Expect(err).ToNot(HaveOccurred())

		})

		AfterEach(func() {
			err = os.RemoveAll(fakeSrcDir)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when trying to update without pr_number", func() {
			It("should return error", func() {
				fakeGithub := &fake.FGithub{
					UpdatePRResult: "fake-ref1",
					UpdatePRError:  nil,
				}
				outCommand := r.NewOutCommand(fakeGithub)

				_, err := outCommand.Run(fakeSrcDir, r.OutRequest{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("no such file or directory"))
			})
		})
	})
})
