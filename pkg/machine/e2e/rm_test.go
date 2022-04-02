package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("podman machine rm", func() {
	var (
		mb      *machineTestBuilder
		testDir string
	)

	BeforeEach(func() {
		_, mb = setup()
	})
	AfterEach(func() {
		teardown(originalHomeDir, testDir, mb)
	})

	It("bad init name", func() {
		i := initMachine{}
		reallyLongName := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		session, err := mb.setName(reallyLongName).setCmd(&i).run()
		Expect(err).To(BeNil())
		Expect(session.ExitCode()).To(Equal(125))
	})
	It("Remove machine", func() {
		i := initMachine{}
		i.withNow().
		session, err := mb.setCmd(i.withImagePath(mb.imagePath)).run()
		Expect(err).To(BeNil())
		Expect(session.ExitCode()).To(Equal(0))
		rm := rmMachine{}
		_, err = mb.setCmd(rm.withForce()).run()
		Expect(err).To(BeNil())

		// Inspecting a non-existent machine should fail
		_, ec, err := mb.toQemuInspectInfo()
		Expect(err).To(BeNil())
		Expect(ec).To(Equal(125))
	})

	It("Remove machine", func() {
		i := initMachine{}
		session, err := mb.setCmd(i.withNow()).run()
		Expect(err).To(BeNil())
		Expect(session.ExitCode()).To(Equal(0))
		rm := rmMachine{}
		_, err = mb.setCmd(rm.withForce()).run()
		Expect(err).To(BeNil())

		// Inspecting a non-existent machine should fail
		_, ec, err := mb.toQemuInspectInfo()
		Expect(err).To(BeNil())
		Expect(ec).To(Equal(125))
	})
})
