package main

import (
	"os"
	"os/exec"
	"sync"
)

func build02(cmds []*exec.Cmd) error {
	var wg sync.WaitGroup

	for _, cmd := range cmds {
		cmd := cmd

		if cmd.Path != dummyCc {
			// コンパイラではない時は、直前までのコンパイルが終わるのを待つ
			wg.Wait()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd.Stdout = os.Stdout
			cmd.Run()
		}()
	}
	wg.Wait()
	return nil
}
