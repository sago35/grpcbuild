package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
)

func build01(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		cmd := exec.Command(fields[0], fields[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	return nil
}
