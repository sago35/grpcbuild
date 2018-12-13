package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/sago35/ochan"
)

func build04(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	outCh := make(chan string, 10000)
	oc := ochan.NewOchan(outCh, 100)
	go func() {
		for ch := range outCh {
			fmt.Println(ch)
		}
	}()

	limit := make(chan struct{}, *threads)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		if fields[0] == dummyCc {
			limit <- struct{}{}
			ch := oc.GetCh()
			go func() {
				defer func() { <-limit }()
				defer close(ch)

				so := new(bytes.Buffer)
				se := new(bytes.Buffer)

				cmd := exec.Command(fields[0], fields[1:]...)
				cmd.Stdout = so
				cmd.Run()

				if so.Len() > 0 {
					ch <- strings.TrimSuffix(so.String(), "\n")
				}
				if se.Len() > 0 {
					ch <- strings.TrimSuffix(se.String(), "\n")
				}

			}()
		} else {
			oc.Wait()

			limit <- struct{}{}
			ch := oc.GetCh()
			go func() {
				defer func() { <-limit }()
				defer close(ch)

				so := new(bytes.Buffer)
				se := new(bytes.Buffer)

				cmd := exec.Command(fields[0], fields[1:]...)
				cmd.Stdout = so
				cmd.Run()

				if so.Len() > 0 {
					ch <- strings.TrimSuffix(so.String(), "\n")
				}
				if se.Len() > 0 {
					ch <- strings.TrimSuffix(se.String(), "\n")
				}

			}()
		}
	}
	oc.Wait()
	return nil
}
