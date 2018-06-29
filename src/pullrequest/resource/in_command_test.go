package resource_test

import (
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

	Context("when request is valid", func() {
		It("should return downloaded version", func() {
			fakeGithub := &fake.FGithub{
				Pulls: []*r.Pull{
					&r.Pull{Number: 1, SHA: "fake-sha1"},
					&r.Pull{Number: 2, SHA: "fake-sha2"},
				},
			}
			inCommand := r.NewInCommand(fakeGithub)
			inRequest := r.InRequest{
				Source:  r.Source{},
				Version: r.Version{Ref: "fake-sha1"},
			}

			inResponse, err := inCommand.Run(fakeDestDir, inRequest)
			Expect(err).ToNot(HaveOccurred())
			Expect(inResponse.Version.Ref).To(Equal("fake-sha1"))
		})
	})

	Context("when there is not pull request", func() {
		It("should return error", func() {
			fakeGithub := &fake.FGithub{Pulls: []*r.Pull{}}
			inCommand := r.NewInCommand(fakeGithub)
			inRequest := r.InRequest{
				Source:  r.Source{},
				Version: r.Version{Ref: "fake-sha1"},
			}

			_, err := inCommand.Run(fakeDestDir, inRequest)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("version fake-sha1 not found"))
		})
	})

	Context("when request is not valid", func() {
		It("should return error", func() {
			fakeGithub := &fake.FGithub{
				Pulls: []*r.Pull{
					&r.Pull{Number: 1, SHA: "fake-sha1"},
					&r.Pull{Number: 2, SHA: "fake-sha2"},
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
