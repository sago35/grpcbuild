package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sago35/ochan"
)

func build04(cmds []*exec.Cmd) error {
	outCh := make(chan string, 10000)
	defer close(outCh)

	oc := ochan.NewOchan(outCh, 100)
	go func() {
		for ch := range outCh {
			fmt.Println(ch)
		}
	}()

	limit := make(chan struct{}, *threads)

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

			so := new(bytes.Buffer)
			cmd.Stdout = so
			cmd.Run()

			if so.Len() > 0 {
				ch <- strings.TrimSpace(so.String())
			}
		}()
	}
	oc.Wait()
	return nil
}
