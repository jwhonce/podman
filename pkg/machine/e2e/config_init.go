package e2e

import (
	"strconv"
)

type initMachine struct {
	/*
	      --cpus uint              Number of CPUs (default 1)
	      --disk-size uint         Disk size in GB (default 100)
	      --ignition-path string   Path to ignition file
	      --image-path string      Path to qcow image (default "testing")
	  -m, --memory uint            Memory in MB (default 2048)
	      --now                    Start machine now
	      --rootful                Whether this machine should prefer rootful container exectution
	      --timezone string        Set timezone (default "local")
	  -v, --volume stringArray     Volumes to mount, source:target
	      --volume-driver string   Optional volume driver

	*/
	cpus         *uint
	diskSize     *uint
	ignitionPath string
	imagePath    string
	memory       *uint
	now          bool
	timezone     string
	volumes      []string

	cmd []string
}

func (i initMachine) buildCmd(names []string) []string {
	cmd := []string{"machine", "init"}
	if i.cpus != nil {
		cmd = append(cmd, "--cpus", strconv.Itoa(int(*i.cpus)))
	}
	if i.diskSize != nil {
		cmd = append(cmd, "--disk-size", strconv.Itoa(int(*i.diskSize)))
	}
	if l := len(i.ignitionPath); l > 0 {
		cmd = append(cmd, "--ignition-path", i.ignitionPath)
	}
	if l := len(i.imagePath); l > 0 {
		cmd = append(cmd, "--image-path", i.imagePath)
	}
	if i.memory != nil {
		cmd = append(cmd, "--memory", strconv.Itoa(int(*i.memory)))
	}
	if l := len(i.timezone); l > 0 {
		cmd = append(cmd, "--timezone", i.timezone)
	}
	for _, v := range i.volumes {
		cmd = append(cmd, "--volume", v)
	}
	if i.now {
		cmd = append(cmd, "--now")
	}
	cmd = append(cmd, names...)
	i.cmd = cmd
	return cmd
}

func (i *initMachine) withCPUs(num uint) machineCommand {
	i.cpus = &num
	return i
}
func (i *initMachine) withDiskSize(size uint) machineCommand {
	i.diskSize = &size
	return i
}

func (i *initMachine) withIgnitionPath(path string) machineCommand {
	i.ignitionPath = path
	return i
}

func (i *initMachine) withImagePath(path string) machineCommand {
	i.imagePath = path
	return i
}

func (i *initMachine) withMemory(num uint) machineCommand {
	i.memory = &num
	return i
}

func (i *initMachine) withNow() machineCommand {
	i.now = true
	return i
}

func (i *initMachine) withTimezone(tz string) machineCommand {
	i.timezone = tz
	return i
}

func (i *initMachine) withVolume(v string) machineCommand {
	i.volumes = append(i.volumes, v)
	return i
}
