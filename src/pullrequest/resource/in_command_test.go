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
	var fakeDestDir string

	BeforeEach(func() {
		fakeDestDir = path.Join(os.TempDir(), "fakedir")
	})

	AfterEach(func() {
		os.Remove(fakeDestDir)
	})

	Context("when version is valid", func() {
		It("should return downloaded version", func() {
			fakeGithub := &fake.FGithub{
				ListPRResult: []*r.Pull{
					&r.Pull{Number: 1, Ref: "fake-ref1"},
					&r.Pull{Number: 2, Ref: "fake-ref2"},
				},
			}
			inCommand := r.NewInCommand(fakeGithub)
			inRequest := r.InRequest{
				Source:  r.Source{},
				Version: r.Version{Ref: "fake-ref1"},
			}

			inResponse, err := inCommand.Run(fakeDestDir, inRequest)
			Expect(err).ToNot(HaveOccurred())
			Expect(inResponse.Version.Ref).To(Equal("fake-ref1"))
		})
	})

	Context("when creating a folder fails", func() {
		It("should return error", func() {
			fakeGithub := &fake.FGithub{
				ListPRResult: []*r.Pull{
					&r.Pull{Number: 1, Ref: "fake-ref1"},
					&r.Pull{Number: 2, Ref: "fake-ref2"},
				},
			}
			inCommand := r.NewInCommand(fakeGithub)
			inRequest := r.InRequest{
				Source:  r.Source{},
				Version: r.Version{Ref: "fake-ref1"},
			}

			_, err := inCommand.Run("/dir/not/exist", inRequest)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when download failed", func() {
		It("should return error", func() {
			fakeGithub := &fake.FGithub{
				ListPRResult: []*r.Pull{
					&r.Pull{Number: 1, Ref: "fake-ref1"},
					&r.Pull{Number: 2, Ref: "fake-ref2"},
				},
				DownloadPRError: errors.New("fake-error"),
			}
			inCommand := r.NewInCommand(fakeGithub)
			inRequest := r.InRequest{
				Source:  r.Source{},
				Version: r.Version{Ref: "fake-ref1"},
			}

			_, err := inCommand.Run(fakeDestDir, inRequest)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("fake-error"))
		})
	})

	Context("when List PRs failed", func() {
		It("should return error", func() {
			fakeGithub := &fake.FGithub{
				ListPRError: errors.New("fake-list-error"),
			}
			inCommand := r.NewInCommand(fakeGithub)
			inRequest := r.InRequest{
				Source:  r.Source{},
				Version: r.Version{Ref: "fake-ref1"},
			}

			_, err := inCommand.Run(fakeDestDir, inRequest)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("fake-list-error"))
		})
	})

	Context("when there is not pull request", func() {
		It("should return error", func() {
			fakeGithub := &fake.FGithub{ListPRResult: []*r.Pull{}}
			inCommand := r.NewInCommand(fakeGithub)
			inRequest := r.InRequest{
				Source:  r.Source{},
				Version: r.Version{Ref: "fake-ref1"},
			}

			_, err := inCommand.Run(fakeDestDir, inRequest)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("version fake-ref1 not found"))
		})
	})

	Context("when version is not valid anymore", func() {
		It("should return error", func() {
			fakeGithub := &fake.FGithub{
				ListPRResult: []*r.Pull{
					&r.Pull{Number: 1, Ref: "fake-ref1"},
					&r.Pull{Number: 2, Ref: "fake-ref2"},
				},
			}
			inCommand := r.NewInCommand(fakeGithub)
			inRequest := r.InRequest{
				Source:  r.Source{},
				Version: r.Version{Ref: "fake-sha3"},
			}

			_, err := inCommand.Run(fakeDestDir, inRequest)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("version fake-sha3 not found"))
		})
	})
})
