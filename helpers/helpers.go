package helpers

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"testing"
)

type Color string

const (
	Red     Color = "\033[31m"
	Green   Color = "\033[32m"
	Blue    Color = "\033[34m"
	Yellow  Color = "\033[33m"
	Magenta Color = "\033[35m"
	Cyan    Color = "\033[36m"
	White   Color = "\033[37m"
	Reset   Color = "\033[0m"
)

func Assert(truthy bool, msg string, t *testing.T) {
	if !truthy {
		t.Errorf("Assert failed: %s", msg)
	}
}

// Colored printf, if a is empty acts like println
func CPrintf(c Color, format string, a ...any) {
	if a == nil {
		fmt.Fprintln(os.Stdout, string(c)+format+string(Reset))
	} else {
		fmt.Fprintf(os.Stdout, string(c)+format+string(Reset), a...)
	}
}

func Wait(wg *sync.WaitGroup, funcs ...func()) {
	wg.Add(len(funcs))

	for _, f := range funcs {
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic in goroutine: %v", r)
				}
			}()

			f()
		}()
	}

	wg.Wait()
}

func GetDigitsFromString(s string) []string {
	return regexp.MustCompile(`\d+`).FindAllString(s, 1)
}

func CaptureStdout(f func()) ([]byte, error) {
	temp, err := os.CreateTemp("/tmp", "output-*.txt")
	defer os.Remove(temp.Name())
	if err != nil {
		return []byte{}, fmt.Errorf("Err while creating temp file: %v", err)
	}

	stdout := os.Stdout
	os.Stdout = temp

	wg := sync.WaitGroup{}

	Wait(&wg,
		func() {
			f()
		},
	)

	os.Stdout = stdout

	temp.Close()

	out, err := os.ReadFile(temp.Name())
	if err != nil {
		return []byte{}, fmt.Errorf("cant read file: %v", err)
	}

	return out, nil
}

func AwaitSIGINT() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)

	go func() {
		<-s
		fmt.Println("\nExiting")
		os.Exit(0)
	}()
}
