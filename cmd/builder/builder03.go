package main

import (
	"fmt"
	"os/exec"
	"sync"
)

func build03(cmds []*exec.Cmd) {
	var wg sync.WaitGroup

	// *threads 分だけ cap を作っておく事で分散数を制御する
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
			buf, _ := cmd.CombinedOutput()
			fmt.Print(string(buf))
		}()
	}
	wg.Wait()
}
