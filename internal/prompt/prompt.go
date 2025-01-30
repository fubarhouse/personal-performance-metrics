package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm will accept input for a prompt.
func Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", prompt)

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}
}
