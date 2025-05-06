package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var binaryName = "tflintenv"

func main() {
	if runtime.GOOS == "windows" {
		binaryName = fmt.Sprintf("%s.exe", binaryName)
	}
	// Get the command-line arguments
	args := os.Args[1:]

	// Store the output in the dst variable
	dst, err := currentBinaryPath()
	if err != nil {
		os.Exit(1)
	}

	if dst == "" || strings.Contains(dst, "no version") {
		cmd := exec.Command(binaryName, "use", defaultVersion())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
		dst, err = currentBinaryPath()
		if err != nil {
			os.Exit(1)
		}
	}

	// Create a new command with dst and the command-line arguments
	cmd := exec.Command(dst, args...)

	// Set the command's Stdin and Stdout to the main process's Stdin and Stdout
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and pass through exit code
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			os.Exit(ee.ExitCode())
		}
		_, _ = os.Stderr.WriteString(fmt.Sprintf("Error executing command but could not get exit code: %s\n", err))
		os.Exit(1)
	}
}

func currentBinaryPath() (string, error) {
	// Create a new command
	cmd := exec.Command(binaryName, "path")

	// Run the command and capture the output
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return "", err
	}
	return string(out), nil
}

func defaultVersion() string {
	v := os.Getenv("TFLINTENV_DEFAULT_VERSION")
	if v == "" {
		v = "latest"
	}
	return v
}
