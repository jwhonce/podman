package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("podman machine start", func() {
	var (
		mb      *machineTestBuilder
		testDir string
	)
	BeforeEach(func() {
		testDir, mb = setup()
	})
	AfterEach(func() {
		teardown(originalHomeDir, testDir, mb)
	})

	It("unknown machine name", func() {
		s := startMachine{}
		session, err := mb.setName("abc123").setCmd(&s).run()
		Expect(err).To(BeNil())
		Expect(session.ExitCode()).To(Equal(125))
	})
})
