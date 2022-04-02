package e2e

type startMachine struct {
	/*
		No command line args other than a machine vm name (also not required)
	*/
	cmd []string
}

func (s startMachine) buildCmd(names []string) []string {
	cmd := []string{"machine", "start"}
	if len(names) > 0 {
		cmd = append(cmd, names[0])
	}
	return cmd
}
