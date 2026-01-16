package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("myshell> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		if input == "exit" {
			break
		}

		executeCommand(input)
	}
}

func executeCommand(input string) {

	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":

		if len(args) < 2 {

			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, "cd: %v\n", err)
				return
			}
			err = os.Chdir(homeDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cd: %v\n", err)
			}
			return
		}
		err := os.Chdir(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "cd: %v\n", err)
		}
		return

	case "pwd":

		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "pwd: %v\n", err)
		} else {
			fmt.Println(dir)
		}
		return

	case "clear":

		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
		return

	case "help":

		fmt.Println("\nAvailable commands:")
		fmt.Println("  cd [dir]   - Change directory")
		fmt.Println("  pwd        - Print working directory")
		fmt.Println("  clear      - Clear screen")
		fmt.Println("  help       - Show this help message")
		fmt.Println("  exit       - Exit the shell")
		fmt.Println("\nAll other commands are executed through cmd.exe")
		return
	}

	var cmd *exec.Cmd
	if len(args) == 1 {
		cmd = exec.Command("cmd", "/C", args[0])
	} else {
		cmd = exec.Command("cmd", "/C", input)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
