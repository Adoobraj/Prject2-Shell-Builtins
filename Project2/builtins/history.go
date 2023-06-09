//agd0109
// Aavinash Doobraj
package builtins

import (
	"fmt"
	"os"
	//"strconv"
)

func history() {
	history := []string{}

	for {
		// Prompt for user input
		fmt.Print("> ")
		var input string
		fmt.Scanln(&input)

		// Add command to history
		history = append(history, input)

		// Handle built-in commands
		switch input {
		case "history":
			// Display history
			for i, cmd := range history {
				fmt.Printf("%d %s\n", i+1, cmd)
			}
		case "exit":
			// Exit the shell
			os.Exit(0)
		default:
			// Execute the command
			fmt.Printf("Executing command: %s\n", input)
		}
	}
}
