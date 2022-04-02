package e2e

import (
	"fmt"
	"io"
	"io/ioutil"
	url2 "net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/containers/podman/v4/pkg/machine"
	"github.com/containers/storage/pkg/reexec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	if reexec.Init() {
		return
	}
	os.Exit(m.Run())
}

const (
	defaultStream string = "podman-testing"
	tmpDir        string = "/tmp"
)

var (
	fqImageName    string
	suiteImageName string
)

// TestLibpod ginkgo master function
func TestMachine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Podman Machine tests")
}

var _ = BeforeSuite(func() {
	fcd, err := machine.GetFCOSDownload(defaultStream)
	if err != nil {
		Fail("unable to get virtual machine image")
	}
	suiteImageName = strings.TrimSuffix(path.Base(fcd.Location), ".xz")
	fqImageName = filepath.Join(tmpDir, suiteImageName)
	if _, err := os.Stat(fqImageName); err != nil {
		if os.IsNotExist(err) {
			getMe, err := url2.Parse(fcd.Location)
			if err != nil {
				Fail(fmt.Sprintf("unable to create url for download: %q", err))
			}
			now := time.Now()
			if err := machine.DownloadVMImage(getMe, fqImageName+".xz"); err != nil {
				Fail(fmt.Sprintf("unable to download machine image: %q", err))
			}
			fmt.Println("Download took: ", time.Since(now).String())
			if err := machine.Decompress(fqImageName+".xz", fqImageName); err != nil {
				Fail(fmt.Sprintf("unable to decompress image file: %q", err))
			}
		} else {
			Fail(fmt.Sprintf("unable to check for cache image: %q", err))
		}
	}
})

var _ = SynchronizedAfterSuite(func() {},
	func() {
		fmt.Println("After")
	})

func setup() (string, *machineTestBuilder) {
	homeDir, err := ioutil.TempDir("", "podman_test")
	if err != nil {
		Fail(fmt.Sprintf("failed to create home directory: %q", err))
	}
	if err := os.MkdirAll(filepath.Join(homeDir, ".ssh"), 0700); err != nil {
		Fail(fmt.Sprintf("failed to create ssh dir: %q", err))
	}
	if err := os.Setenv("HOME", homeDir); err != nil {
		Fail("failed to set home dir")
	}
	mb, err := new()
	if err != nil {
		Fail(fmt.Sprintf("failed to create machine test: %q", err))
	}
	f, err := os.Open(fqImageName)
	if err != nil {
		Fail(fmt.Sprintf("failed to open file %s: %q", fqImageName, err))
	}
	mb.imagePath = filepath.Join(homeDir, suiteImageName)
	n, err := os.Create(mb.imagePath)
	if err != nil {
		Fail(fmt.Sprintf("failed to create file %s: %q", mb.imagePath, err))
	}
	if _, err := io.Copy(n, f); err != nil {
		Fail(fmt.Sprintf("failed to copy %ss to %s: %q", fqImageName, mb.imagePath, err))
	}
	return homeDir, mb
}

func teardown(origHomeDir string, testDir string, mb *machineTestBuilder) {
	r := rmMachine{}
	if _, err := mb.setCmd(r.withForce()).run(); err != nil {
		fmt.Printf("error occured rm'ing machine: %q\n", err)
	}
	if err := os.RemoveAll(testDir); err != nil {
		Fail(fmt.Sprintf("failed to remove test dir: %q", err))
	}
	// this needs to be last in teardown
	if err := os.Setenv("HOME", origHomeDir); err != nil {
		Fail("failed to set home dir")
	}
}
