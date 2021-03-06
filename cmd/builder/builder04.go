package main

import (
	"fmt"
	"os/exec"

	"github.com/sago35/ochan"
)

func build04(cmds []*exec.Cmd) {
	outCh := make(chan string, 10000)
	done := make(chan struct{})

	go func() {
		for ch := range outCh {
			fmt.Print(ch)
		}
		close(done)
	}()

	limit := make(chan struct{}, *threads)

	oc := ochan.NewOchan(outCh, 100)
	for _, cmd := range cmds {
		cmd := cmd

		if cmd.Path != dummyCc {
			// コンパイラではない時は、直前までのコンパイルが終わるのを待つ
			oc.Wait()
		}

		limit <- struct{}{}
		ch := oc.GetCh()
		go func() {
			defer func() { <-limit }()
			defer close(ch)

			buf, _ := cmd.CombinedOutput()
			if len(buf) > 0 {
				ch <- string(buf)
			}
		}()
	}
	oc.Wait()
	close(outCh)
	<-done
}
