package e2e

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/containers/podman/v4/pkg/machine"
	"github.com/containers/podman/v4/pkg/machine/qemu"
	. "github.com/onsi/ginkgo" //nolint:golint,stylecheck
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/gexec" //nolint:golint,stylecheck
)

var originalHomeDir = os.Getenv("HOME")

const (
	defaultTimeout time.Duration = 90 * time.Second
)

type machineCommand interface {
	buildCmd(Names []string) []string
}

type MachineTestBuilder interface {
	setName(string) MachineTestBuilder
	setCmd(mc machineCommand) MachineTestBuilder
	setTimeout(duration time.Duration) MachineTestBuilder
	run() (*machineSession, error)
}
type machineSession struct {
	*gexec.Session
}

type machineTestBuilder struct {
	cmd          []string
	imagePath    string
	names        []string
	podmanBinary string
	timeout      time.Duration
}
type qemuMachineInspectInfo struct {
	State machine.MachineStatus
	VM    qemu.MachineVM
}

// waitWithTimeout waits for a command to complete for a given
// number of seconds
func (ms *machineSession) waitWithTimeout(timeout time.Duration) {
	Eventually(ms, timeout).Should(Exit())
	os.Stdout.Sync()
	os.Stderr.Sync()
}

func (ms *machineSession) Bytes() []byte {
	return []byte(ms.outputToString())
}

// outputToString returns the output from a session in string form
func (ms *machineSession) outputToString() string {
	if ms == nil || ms.Out == nil || ms.Out.Contents() == nil {
		return ""
	}

	fields := strings.Fields(string(ms.Out.Contents()))
	return strings.Join(fields, " ")
}

// new constructor for machine test builders
func new() (*machineTestBuilder, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	podmanBinary := filepath.Join(cwd, "../../../bin/podman")
	if os.Getenv("PODMAN_BINARY") != "" {
		podmanBinary = os.Getenv("PODMAN_BINARY")
	}
	mb := machineTestBuilder{podmanBinary: podmanBinary, timeout: defaultTimeout}
	return &mb, nil
}

// setName adds a name to the array of names.  For commands that can take
// multiple names, this can be used more than once.
func (m *machineTestBuilder) setName(name string) MachineTestBuilder {
	m.names = append(m.names, name)
	return m
}

// setCmd takes a machineCommand struct and assembles a cmd line
// representation of the podman machine command
func (m *machineTestBuilder) setCmd(mc *machineCommand) MachineTestBuilder {
	// If no name for the machine exists, we set a random name.
	if len(m.names) < 1 {
		m.names = []string{randomString(12)}
	}
	m.cmd = mc.buildCmd(m.names)
	return m
}

func (m *machineTestBuilder) setTimeout(timeout time.Duration) MachineTestBuilder {
	m.timeout = timeout
	return m
}

// toQemuInspectInfo is only for inspecting qemu machines.  Other providers will need
// to make their own.
func (mb *machineTestBuilder) toQemuInspectInfo() ([]qemuMachineInspectInfo, int, error) {
	args := []string{"machine", "inspect"}
	args = append(args, mb.names...)
	session, err := runWrapper(mb.podmanBinary, args, defaultTimeout)
	if err != nil {
		return nil, -1, err
	}
	mii := []qemuMachineInspectInfo{}
	err = json.Unmarshal(session.Bytes(), &mii)
	return mii, session.ExitCode(), err
}

func (m *machineTestBuilder) run() (*machineSession, error) {
	return runWrapper(m.podmanBinary, m.cmd, m.timeout)
}

func runWrapper(podmanBinary string, cmdArgs []string, timeout time.Duration) (*machineSession, error) {
	fmt.Println(podmanBinary + " " + strings.Join(cmdArgs, " "))
	c := exec.Command(podmanBinary, cmdArgs...)
	session, err := Start(c, GinkgoWriter, GinkgoWriter)
	if err != nil {
		Fail(fmt.Sprintf("Unable to start session: %q", err))
		return nil, err
	}
	ms := machineSession{session}
	ms.waitWithTimeout(timeout)
	fmt.Println("output:", ms.outputToString())
	return &ms, nil
}

func (m *machineTestBuilder) init() {}

// randomString returns a string of given length composed of random characters
func randomString(n int) string {
	var randomLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = randomLetters[rand.Intn(len(randomLetters))]
	}
	return string(b)
}
