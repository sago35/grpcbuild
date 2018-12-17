package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sago35/limichan"
	"github.com/sago35/ochan"
)

type remoteInfo struct {
	addr    string
	port    int
	threads int
}

func run(threads int, ris []remoteInfo) error {
	startTime := time.Now()

	l, ctx := limichan.New(context.Background())
	if false {
		fmt.Println(ctx)
	}

	for i := 0; i < threads; i++ {
		l.AddWorker(newLocalWorker(fmt.Sprintf("localworker(%d)", i)))
	}

	for _, ri := range ris {
		ri := ri

		go func() {
			w, err := newWorker(ri.addr, ri.port, ri.threads)
			if err != nil {
				log.Fatal(err)
			}

			for i := 0; i < w.threads; i++ {
				l.AddWorker(w.Clone(fmt.Sprintf("%s(%d)", w.name, i)))
			}
		}()
	}

	result := make(chan string, 10000)

	go func() {
		oc := ochan.NewOchan(result, 10000)

		jobs, err := makeJobs(`input.txt`)
		if err != nil {
			log.Fatal(err)
		}

		for _, job := range jobs {
			if job.Cmds[0].Path == "dummyld" {
				oc.Wait()
			}

			job := job
			job.ch = oc.GetCh()
			err = l.Do(job)
			if err != nil {
				close(job.ch)
				log.Print(err)
			}
		}
		oc.Wait()
		err = l.Wait()
		if err != nil {
			log.Fatal(err)
		}

		result <- fmt.Sprintf("%dms", time.Since(startTime).Nanoseconds()/1000000)
		close(result)
	}()

	for r := range result {
		fmt.Println(r)
	}

	return nil
}

func makeJobs(path string) ([]*job, error) {
	ret := []*job{}

	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		j := newJob()
		j.name = fields[2]
		j.Cmds = append(j.Cmds, exec.Cmd{
			Path: fields[0],
			Args: fields,
		})

		j.outFile = []string{fields[2]}

		if strings.HasPrefix(fields[0], "dummyld") {
			j.depFile = fields[3:]
		}

		ret = append(ret, j)
	}

	return ret, nil
}

const (
	dummyCc = "dummycc"
)

func makeFileList() ([]string, error) {
	ret := []string{}

	b, err := exec.Command(`git`, `ls-files`, `testdata`).Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	return ret, nil
}

var (
	threads = flag.Int("threads", 1, "threads")
)

func main() {
	flag.Parse()

	remoteInfos := []remoteInfo{}

	for _, arg := range flag.Args() {
		f := strings.SplitN(arg, ":", 2)
		ri := remoteInfo{
			addr: f[0],
		}

		_, err := fmt.Sscanf(f[1], "%d:%d", &ri.port, &ri.threads)
		if err != nil {
			log.Fatal(err)
		}
		remoteInfos = append(remoteInfos, ri)
	}

	err := run(*threads, remoteInfos)
	if err != nil {
		log.Fatal(err)
	}
}
