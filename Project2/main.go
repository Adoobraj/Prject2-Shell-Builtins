package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/jh125486/CSCE4600/Project2/builtins"
	//"./builtins" //location of builtins folder
)

func main() {
	history := []string{}
	exit := make(chan struct{}, 2) // buffer this so there's no deadlock.
	runLoop(os.Stdin, os.Stdout, os.Stderr, exit, history)
}

func runLoop(r io.Reader, w, errW io.Writer, exit chan struct{}, history []string) {

	var (
		input    string
		err      error
		readLoop = bufio.NewReader(r)
	)
	for {
		select {
		case <-exit:
			_, _ = fmt.Fprintln(w, "exiting gracefully...")
			return
		default:
			if err := printPrompt(w); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if input, err = readLoop.ReadString('\n'); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if err = handleInput(w, input, exit, &history); err != nil {
				_, _ = fmt.Fprintln(errW, err)
			}
		}
	}
}

func printPrompt(w io.Writer) error {
	// Get current user.
	// Don't prematurely memoize this because it might change due to `su`?
	u, err := user.Current()
	if err != nil {
		return err
	}
	// Get current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// /home/User [Username] $
	_, err = fmt.Fprintf(w, "%v [%v] $ ", wd, u.Username)

	return err
}

func handleInput(w io.Writer, input string, exit chan<- struct{}, history *[]string) error {
	// Remove trailing spaces.
	input = strings.TrimSpace(input)

	// Save input to history
	*history = append(*history, input)

	// Split the input separate the command name and the command arguments.
	args := strings.Split(input, " ")
	name, args := args[0], args[1:]

	// Check for built-in commands.
	// New builtin commands should be added here. Eventually this should be refactored to its own func.
	switch name {
	case "cd":
		return builtins.ChangeDirectory(args...)
	case "env":
		return builtins.EnvironmentVariables(w, args...)
	case "pwd":
		return builtins.PrintWorkingDirectory(w)
	case "mkdir":
		return builtins.Mkdir(args...)
	case "ls":
		return builtins.Ls(args...)
	case "exit":
		exit <- struct{}{}
		return nil
	case "history":
		return showHistory(w, history)
	case "rm":
		return builtins.RemoveFile(args...)
	}

	return executeCommand(name, args...)
}

func executeCommand(name string, arg ...string) error {
	// Otherwise prep the command
	cmd := exec.Command(name, arg...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}

func showHistory(w io.Writer, history *[]string) error {
	_, err := fmt.Fprintf(w, "Command history:\n")
	if err != nil {
		return err
	}
	for i, cmd := range *history {
		_, err := fmt.Fprintf(w, "%d: %s\n", i+1, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
