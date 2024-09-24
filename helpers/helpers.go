package helpers

import (
	"fmt"
	"os"
	"sync"
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

func Assert(truthy bool, msg string) error {
	if !truthy {
		return fmt.Errorf("Assert failed: %s", msg)
	}

	return nil
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
			f()
		}()
	}

	wg.Wait()
}
