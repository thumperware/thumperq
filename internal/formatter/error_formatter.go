package formatter

import "fmt"

func FormatErr(methodPath string, err error) error {
	return fmt.Errorf("%s:%v", methodPath, err)
}
