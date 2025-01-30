package prompt

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestAskForConfirmation(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Yes input", "y\n", true},
		{"No input", "n\n", false},
		{"Invalid input then Yes", "o\ny\n", true},
		{"Invalid input then No", "o\nn\n", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save original stdin and stdout
			oldStdin := os.Stdin
			oldStdout := os.Stdout

			// Create pipes for stdin and stdout
			inR, inW, _ := os.Pipe()
			outR, outW, _ := os.Pipe()

			// Set stdin and stdout to the pipes
			os.Stdin = inR
			os.Stdout = outW

			// Write test input
			go func() {
				inW.Write([]byte(tc.input))
				inW.Close()
			}()

			// Capture stdout
			outC := make(chan string)
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, outR)
				outC <- buf.String()
			}()

			// Defer cleanup
			defer func() {
				os.Stdin = oldStdin
				os.Stdout = oldStdout
				outW.Close()
			}()

			// Call the function
			result := Confirm("Test prompt")

			// Close the write end of the stdout pipe
			outW.Close()
			captured := <-outC

			// Check the result
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}

			// Check if invalid input message is printed for 'o' input
			if strings.Contains(tc.input, "o\n") {
				expectedOutput := "Invalid input. Please enter 'y' or 'n'.\n"
				if !strings.Contains(captured, expectedOutput) {
					t.Errorf("Expected output to contain %q, but it didn't", expectedOutput)
				}
			}
		})
	}
}
