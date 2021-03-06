package main

import (
	"bufio"
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var (
	mode    = flag.Int("m", 1, "mode")
	threads = flag.Int("threads", 1, "threads")
)

const (
	dummyCc = "dummycc"
)

func main() {
	flag.Parse()

	cmds := makeCmds()

	start := time.Now()

	switch *mode {
	case 1:
		build01(cmds)
	case 2:
		build02(cmds)
	case 3:
		build03(cmds)
	case 4:
		build04(cmds)
	case 5:
		build05(cmds)
	case 6:
		build06(cmds)
	case 7:
		build07(cmds)
	default:
		panic("error")
	}
	fmt.Printf("%d ms\n", time.Since(start).Nanoseconds()/1000000)
}

func makeCmds() []*exec.Cmd {
	xxx := `dummycc -o testdata/aa.o testdata/aa.c
dummycc -o testdata/ab.o testdata/ab.c
dummycc -o testdata/ac.o testdata/ac.c
dummycc -o testdata/ad.o testdata/ad.c
dummycc -o testdata/ae.o testdata/ae.c
dummycc -o testdata/af.o testdata/af.c
dummycc -o testdata/ba.o testdata/ba.c
dummycc -o testdata/bb.o testdata/bb.c
dummycc -o testdata/bc.o testdata/bc.c
dummycc -o testdata/bd.o testdata/bd.c
dummycc -o testdata/be.o testdata/be.c
dummycc -o testdata/bf.o testdata/bf.c
dummycc -o testdata/ca.o testdata/ca.c
dummycc -o testdata/cb.o testdata/cb.c
dummycc -o testdata/cc.o testdata/cc.c
dummycc -o testdata/cd.o testdata/cd.c
dummycc -o testdata/ce.o testdata/ce.c
dummycc -o testdata/cf.o testdata/cf.c
dummycc -o testdata/da.o testdata/da.c
dummycc -o testdata/db.o testdata/db.c
dummycc -o testdata/dc.o testdata/dc.c
dummycc -o testdata/dd.o testdata/dd.c
dummycc -o testdata/de.o testdata/de.c
dummycc -o testdata/df.o testdata/df.c
dummycc -o testdata/ea.o testdata/ea.c
dummycc -o testdata/eb.o testdata/eb.c
dummycc -o testdata/ec.o testdata/ec.c
dummycc -o testdata/ed.o testdata/ed.c
dummycc -o testdata/ee.o testdata/ee.c
dummycc -o testdata/ef.o testdata/ef.c
dummycc -o testdata/fa.o testdata/fa.c
dummycc -o testdata/fb.o testdata/fb.c
dummycc -o testdata/fc.o testdata/fc.c
dummycc -o testdata/fd.o testdata/fd.c
dummycc -o testdata/fe.o testdata/fe.c
dummycc -o testdata/ff.o testdata/ff.c
dummyld -o testdata/a.out testdata/aa.o testdata/ab.o testdata/ac.o testdata/ad.o testdata/ae.o testdata/af.o testdata/ba.o testdata/bb.o testdata/bc.o testdata/bd.o testdata/be.o testdata/bf.o testdata/ca.o testdata/cb.o testdata/cc.o testdata/cd.o testdata/ce.o testdata/cf.o testdata/da.o testdata/db.o testdata/dc.o testdata/dd.o testdata/de.o testdata/df.o testdata/ea.o testdata/eb.o testdata/ec.o testdata/ed.o testdata/ee.o testdata/ef.o testdata/fa.o testdata/fb.o testdata/fc.o testdata/fd.o testdata/fe.o testdata/ff.o
`

	cmds := []*exec.Cmd{}
	scanner := bufio.NewScanner(strings.NewReader(xxx))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		cmds = append(cmds, &exec.Cmd{
			Path: fields[0],
			Args: fields,
		})
	}

	return cmds
}
