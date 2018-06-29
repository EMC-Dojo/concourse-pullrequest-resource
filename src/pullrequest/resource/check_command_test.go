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
					Pulls: []*r.Pull{
						&r.Pull{Number: 1, SHA: "fake-sha1"},
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
				Expect(versions[0].Ref).To(Equal("fake-sha1"))
			})
		})

		Context("when there is no pull", func() {
			It("should return empty array", func() {
				fakeGithub := &fake.FGithub{Pulls: []*r.Pull{}}
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
					Pulls: []*r.Pull{
						&r.Pull{Number: 1, SHA: "fake-sha1"},
						&r.Pull{Number: 2, SHA: "fake-sha2"},
						&r.Pull{Number: 3, SHA: "fake-sha3"},
					},
				}
				checkCommand := r.NewCheckCommand(fakeGithub)
				checkRequest := r.CheckRequest{
					Source:  r.Source{},
					Version: r.Version{Ref: "fake-sha3"},
				}

				versions, err := checkCommand.Run(checkRequest)
				Expect(err).ToNot(HaveOccurred())
				Expect(versions).To(HaveLen(1))
				Expect(versions[0].Ref).To(Equal("fake-sha3"))
			})
		})

		Context("when given version is not the latest", func() {
			It("should return newer versions with given version", func() {
				fakeGithub := &fake.FGithub{
					Pulls: []*r.Pull{
						&r.Pull{Number: 1, SHA: "fake-sha1"},
						&r.Pull{Number: 2, SHA: "fake-sha2"},
						&r.Pull{Number: 3, SHA: "fake-sha3"},
					},
				}
				checkCommand := r.NewCheckCommand(fakeGithub)
				checkRequest := r.CheckRequest{
					Source:  r.Source{},
					Version: r.Version{Ref: "fake-sha2"},
				}

				versions, err := checkCommand.Run(checkRequest)
				Expect(err).ToNot(HaveOccurred())
				Expect(versions).To(HaveLen(2))
				Expect(versions[0].Ref).To(Equal("fake-sha2"))
				Expect(versions[1].Ref).To(Equal("fake-sha3"))
			})
		})
	})
})
