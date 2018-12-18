package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sago35/limichan"
	"github.com/sago35/ochan"
)

type job struct {
	cmd     *exec.Cmd
	ch      chan string
	outFile []string
	depFile []string
}

type lworker struct {
	name string
}

func (w *lworker) Do(ctx context.Context, lj limichan.Job) error {
	job, _ := lj.(*job)
	defer close(job.ch)

	cmd := job.cmd

	so := new(bytes.Buffer)
	cmd.Stdout = so
	cmd.Run()

	if so.Len() > 0 {
		job.ch <- strings.TrimSpace(so.String())
	}

	return nil
}

func build05(cmds []*exec.Cmd) error {
	outCh := make(chan string, 10000)
	defer close(outCh)

	oc := ochan.NewOchan(outCh, 100)
	go func() {
		for ch := range outCh {
			fmt.Println(ch)
		}
	}()

	l, _ := limichan.New(context.Background())
	for i := 0; i < *threads; i++ {
		l.AddWorker(&lworker{name: fmt.Sprintf("lworker(%d)", i)})
	}

	for _, cmd := range cmds {
		cmd := cmd

		if cmd.Path != dummyCc {
			// コンパイラではない時は、直前までのコンパイルが終わるのを待つ
			oc.Wait()
		}

		ch := oc.GetCh()
		j := &job{
			cmd: cmd,
			ch:  ch,
		}

		l.Do(j)
	}

	return nil
}
