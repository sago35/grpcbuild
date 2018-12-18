package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	mode    = flag.Int("m", 1, "mode")
	threads = flag.Int("threads", 1, "threads")
)

const (
	dummyCc = "dummycc.exe"
)

func main() {
	flag.Parse()

	cmds, err := read(`input.txt`)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	switch *mode {
	case 1:
		err = build01(cmds)
	case 2:
		err = build02(cmds)
	case 3:
		err = build03(cmds)
	case 4:
		err = build04(cmds)
	case 5:
		err = build05(cmds)
	case 6:
		err = build06(cmds)
	default:
		panic("error")
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d ms\n", time.Since(start).Nanoseconds()/1000000)
}

func read(file string) ([]*exec.Cmd, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	cmds := []*exec.Cmd{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		cmds = append(cmds, exec.Command(fields[0], fields[1:]...))
	}

	return cmds, nil
}
