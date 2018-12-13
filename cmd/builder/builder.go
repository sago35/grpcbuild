package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	r, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	start := time.Now()

	switch *mode {
	case 1:
		err = build01(r)
	case 2:
		err = build02(r)
	case 3:
		err = build03(r)
	case 4:
		err = build04(r)
	default:
		panic("error")
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d ms\n", time.Since(start).Nanoseconds()/1000000)
}
