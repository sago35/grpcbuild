package main

import (
	"os"
	"os/exec"
	"sync"
)

func build03(cmds []*exec.Cmd) error {
	var wg sync.WaitGroup

	limit := make(chan struct{}, *threads)

	for _, cmd := range cmds {
		cmd := cmd

		if cmd.Path != dummyCc {
			// コンパイラではない時は、直前までのコンパイルが終わるのを待つ
			wg.Wait()
		}

		limit <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() { <-limit }()
			defer wg.Done()
			cmd.Stdout = os.Stdout
			cmd.Run()
		}()
	}
	wg.Wait()
	return nil
}
