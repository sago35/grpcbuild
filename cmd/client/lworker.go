package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sago35/limichan"
)

type lworker struct {
	name string
}

func newLocalWorker(name string) *lworker {
	w := &lworker{
		name: name,
	}
	return w
}

func (w *lworker) Do(ctx context.Context, lj limichan.Job) error {
	job, ok := lj.(*job)
	if !ok {
		log.Fatal("type error")
	}
	defer close(job.ch)

	start := time.Now()
	job.ch <- fmt.Sprintf("# %-16s %-16s %d", w.name, job.name, time.Now().UnixNano()/1000000)

	for _, cmd := range job.Cmds {
		so := new(bytes.Buffer)
		se := new(bytes.Buffer)

		cmd.Stdout = so
		cmd.Stderr = se

		err := cmd.Run()
		if err != nil {
			return err
		}

		if so.Len() > 0 {
			job.ch <- strings.TrimSpace(so.String())
		}
		if se.Len() > 0 {
			job.ch <- strings.TrimSpace(se.String())
		}
	}

	job.ch <- fmt.Sprintf("# %-16s %-16s %d (%dms)", w.name, job.name, time.Now().UnixNano()/1000000, time.Since(start).Nanoseconds()/1000000)

	return nil
}
