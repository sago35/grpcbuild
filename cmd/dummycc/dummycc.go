package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	output = flag.String("o", "", "Place the output into <file>")
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 || len(*output) == 0 {
		log.Fatalf("usage : dummycc -o output.o input.c")
	}

	input := flag.Arg(0)

	r, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)

	if scanner.Scan() {
		waitMs, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("compile %s (wait %dms)\n", input, waitMs)

		buf := []string{}
		for scanner.Scan() {
			buf = append(buf, scanner.Text())
		}

		if len(buf) == 0 {
			time.Sleep(time.Duration(waitMs) * time.Millisecond)
		} else {
			for _, s := range buf {
				time.Sleep(time.Duration(int(waitMs)/len(buf)) * time.Millisecond)
				fmt.Printf("%s\n", s)
			}
		}

		w2, err := os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer w2.Close()
	} else {
		fmt.Printf("link %s (wait %dms)\n", input, 1000)
		time.Sleep(1 * time.Second)
	}
}
