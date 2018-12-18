package main

import (
	"os"
	"os/exec"
)

func build01(cmds []*exec.Cmd) error {
	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	return nil
}
