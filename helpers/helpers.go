package helpers

import (
	"fmt"
	"os"
)

// NOTE: maybe make a helper module for colored prints

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

func CPrintf(c Color, format string, a interface{}) {
	fmt.Fprintf(os.Stdout,string(c)+format+string(Reset), a)
}
