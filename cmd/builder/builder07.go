package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/sago35/limichan"
	"github.com/sago35/ochan"
)

func build07(cmds []*exec.Cmd) {
	outCh := make(chan string, 10000)
	done := make(chan struct{})

	go func() {
		for ch := range outCh {
			fmt.Print(ch)
		}
		close(done)
	}()

	l, _ := limichan.New(context.Background())
	w, _ := newLocalWorker()
	for i := 0; i < *threads; i++ {
		l.AddWorker(w)
	}
	go func() {
		w, _ := newWorker(`127.0.0.1`, 12345, *threads)
		time.Sleep(1 * time.Second)
		for i := 0; i < *threads; i++ {
			l.AddWorker(w)
		}
	}()

	oc := ochan.NewOchan(outCh, 100)
	for _, cmd := range cmds {
		cmd := cmd

		if cmd.Path != dummyCc {
			// コンパイラではない時は、直前までのコンパイルが終わるのを待つ
			oc.Wait()
		}

		j := &job{
			cmd:     cmd,
			ch:      oc.GetCh(),
			outFile: cmd.Args[2:3],
			depFile: cmd.Args[3:],
		}

		l.Do(j)
	}
	oc.Wait()
	l.Wait()
	close(outCh)
	<-done
}
