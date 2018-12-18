package main

import (
	"fmt"
	"os/exec"
	"sync"
)

func build02(cmds []*exec.Cmd) {
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
			buf, _ := cmd.CombinedOutput()
			fmt.Print(string(buf))
		}()
	}
	wg.Wait()
}
