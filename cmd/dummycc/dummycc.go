package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

	if strings.HasPrefix(os.Args[0], "dummycc") {
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
			w2.Close()
		} else {
			log.Fatal("error")
		}
	} else {

		for _, f := range flag.Args() {
			fi, err := os.Stat(f)
			if err != nil {
				log.Fatal(err)
			}
			if false {
				fmt.Println(fi)
			}
		}

		w2, err := os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		w2.Close()

		fmt.Printf("link %s (wait %dms)\n", input, 1000)
		time.Sleep(1 * time.Second)
	}
}
