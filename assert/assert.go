package assert 

import "fmt"

func Assert(truthy bool, msg string) error {
	if !truthy {
		return fmt.Errorf("Assert failed: %s", msg)
	}

	return nil
}
