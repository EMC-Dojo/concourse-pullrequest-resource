package resource_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	r "pullrequest/resource"
	"pullrequest/resource/fake"
)

var _ = Describe("CheckCommand", func() {
	Context("when request is valid", func() {
		Context("when there is a pull", func() {
			It("should return latest versions", func() {
				fakeGithub := &fake.FGithub{
					ListPRResult: []*r.Pull{
						&r.Pull{Number: 1, Ref: "fake-ref1"},
					},
				}
				checkCommand := r.NewCheckCommand(fakeGithub)
				checkRequest := r.CheckRequest{
					Source:  r.Source{},
					Version: r.Version{},
				}

				versions, err := checkCommand.Run(checkRequest)
				Expect(err).ToNot(HaveOccurred())
				Expect(versions).To(HaveLen(1))
				Expect(versions[0].Ref).To(Equal("fake-ref1"))
			})
		})

		Context("when there is no pull", func() {
			It("should return empty array", func() {
				fakeGithub := &fake.FGithub{ListPRResult: []*r.Pull{}}
				checkCommand := r.NewCheckCommand(fakeGithub)
				checkRequest := r.CheckRequest{
					Source:  r.Source{},
					Version: r.Version{},
				}

				versions, err := checkCommand.Run(checkRequest)
				Expect(err).ToNot(HaveOccurred())
				Expect(versions).To(HaveLen(0))
			})
		})

		Context("when given version is already the latest", func() {
			It("should return only the latest version", func() {
				fakeGithub := &fake.FGithub{
					ListPRResult: []*r.Pull{
						&r.Pull{Number: 1, Ref: "fake-ref1"},
						&r.Pull{Number: 2, Ref: "fake-ref2"},
						&r.Pull{Number: 3, Ref: "fake-ref3"},
					},
				}
				checkCommand := r.NewCheckCommand(fakeGithub)
				checkRequest := r.CheckRequest{
					Source:  r.Source{},
					Version: r.Version{Ref: "fake-ref3"},
				}

				versions, err := checkCommand.Run(checkRequest)
				Expect(err).ToNot(HaveOccurred())
				Expect(versions).To(HaveLen(1))
				Expect(versions[0].Ref).To(Equal("fake-ref3"))
			})
		})

		Context("when given version is not the latest", func() {
			It("should return newer versions with given version", func() {
				fakeGithub := &fake.FGithub{
					ListPRResult: []*r.Pull{
						&r.Pull{Number: 1, Ref: "fake-ref1"},
						&r.Pull{Number: 2, Ref: "fake-ref2"},
						&r.Pull{Number: 3, Ref: "fake-ref3"},
					},
				}
				checkCommand := r.NewCheckCommand(fakeGithub)
				checkRequest := r.CheckRequest{
					Source:  r.Source{},
					Version: r.Version{Ref: "fake-ref2"},
				}

				versions, err := checkCommand.Run(checkRequest)
				Expect(err).ToNot(HaveOccurred())
				Expect(versions).To(HaveLen(2))
				Expect(versions[0].Ref).To(Equal("fake-ref2"))
				Expect(versions[1].Ref).To(Equal("fake-ref3"))
			})
		})
	})
})
