package builtins

/*
 _				___  _ _ _____
| |_ _ __ __ _ / _ \/ / |___  |
| __| '__/ _` | | | | | |  / /
| |_| | | (_| | |_| | | | / /
 \__|_|  \__, |\___/|_|_|/_/
		 |___/
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Mkdir(args ...string) error {
	var (
		pflag  bool        // create parent directories if they don't exist
		vflag  bool        // verbose mode
		mode   os.FileMode // file permission mode
		target string      // target directory to create
	)

	// Parse command-line flags
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			switch arg {
			case "-p":
				pflag = true
			case "-v":
				vflag = true
			case "-m":
				if i+1 >= len(args) {
					return fmt.Errorf("missing mode argument for -m flag")
				}
				m, err := parseMode(args[i+1])
				if err != nil {
					return err
				}
				mode = m
				i++
			default:
				return fmt.Errorf("invalid flag: %s", arg)
			}
		} else {
			// Target directory is the first non-flag argument
			target = arg
			break
		}
	}

	if target == "" {
		return fmt.Errorf("missing operand")
	}

	if mode == 0 {
		// If -m flag is not present, set default mode to 0755
		umask := os.FileMode(022)
		umaskStr, ok := os.LookupEnv("UMASK")
		if ok {
			umaskInt, err := strconv.ParseInt(umaskStr, 8, 32)
			if err == nil {
				umask = os.FileMode(umaskInt)
			}
		}
		switch umask {
		case 077, 027:
			mode = 0700
		default:
			mode = 0755
		}
	}

	if pflag {
		if err := os.MkdirAll(target, mode); err != nil {
			return err
		}
		if vflag {
			fmt.Printf("mkdir: created directory '%s'\n", target)
		}
	} else {
		// Create the parent directories if they don't exist
		parentDir := filepath.Dir(target)
		if parentDir != "." && parentDir != "/" {
			if err := os.MkdirAll(parentDir, mode); err != nil {
				return err
			}
		}
		if err := os.Mkdir(target, mode); err != nil {
			return err
		}
		if vflag {
			fmt.Printf("mkdir: created directory '%s'\n", target)
		}
	}

	return nil
}

// parseMode parses the file permission mode from a string in octal or symbolic format.
func parseMode(s string) (os.FileMode, error) {
	if strings.HasPrefix(s, "0") {
		// Parse octal notation
		mode, err := strconv.ParseInt(s, 8, 32)
		if err != nil {
			return 0, err
		}
		return os.FileMode(mode), nil
	} else {
		// Parse symbolic notation
		mode, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			return 0, err
		}
		return os.FileMode(mode), nil
	}
}
