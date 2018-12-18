package main

import (
	"fmt"
	"os/exec"
)

func build01(cmds []*exec.Cmd) {
	for _, cmd := range cmds {
		buf, _ := cmd.CombinedOutput()
		fmt.Print(string(buf))
	}
}
