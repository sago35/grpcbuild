package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func build03(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	var wg sync.WaitGroup

	limit := make(chan struct{}, *threads)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		if fields[0] == dummyCc {
			limit <- struct{}{}
			wg.Add(1)
			go func() {
				defer func() { <-limit }()
				defer wg.Done()
				cmd := exec.Command(fields[0], fields[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Run()
			}()
		} else {
			wg.Wait()

			limit <- struct{}{}
			wg.Add(1)
			go func() {
				defer func() { <-limit }()
				defer wg.Done()
				cmd := exec.Command(fields[0], fields[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Run()
			}()
		}
	}
	wg.Wait()
	return nil
}
