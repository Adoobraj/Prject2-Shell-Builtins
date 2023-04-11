package builtins

import (
	"fmt"
	"io"
	"os"
)

func PrintWorkingDirectory(w io.Writer) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Current working directory (made by Ted)--> %s\n", pwd)
	return err
}