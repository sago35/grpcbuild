package main

import (
	"context"
	"fmt"
	"os/exec"

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

func newLocalWorker() (*lworker, error) {
	return &lworker{}, nil
}

func (w *lworker) Do(ctx context.Context, lj limichan.Job) error {
	job, _ := lj.(*job)
	defer close(job.ch)

	cmd := job.cmd

	buf, _ := cmd.CombinedOutput()
	if len(buf) > 0 {
		job.ch <- string(buf)
	}

	return nil
}

func build05(cmds []*exec.Cmd) {
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

	oc := ochan.NewOchan(outCh, 100)
	for _, cmd := range cmds {
		cmd := cmd

		if cmd.Path != dummyCc {
			// コンパイラではない時は、直前までのコンパイルが終わるのを待つ
			oc.Wait()
		}

		j := &job{
			cmd: cmd,
			ch:  oc.GetCh(),
		}

		l.Do(j)
	}
	oc.Wait()
	l.Wait()
	close(outCh)
	<-done
}
