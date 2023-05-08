// erb0149
package builtins

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ls(args ...string) error {
	var dir string
	if len(args) == 0 {
		// No arguments, list contents of current directory.
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	} else {
		// List contents of specified directory.
		dir = args[0]
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}

	return nil
}
