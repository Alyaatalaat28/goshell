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

	case "set":

		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "set: usage: set VAR_NAME value")
			return
		}
		varName := args[1]
		varValue := strings.Join(args[2:], " ")
		err := os.Setenv(varName, varValue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "set: %v\n", err)
		} else {
			fmt.Printf("Set %s=%s\n", varName, varValue)
		}
		return

	case "get":

		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "get: usage: get VAR_NAME")
			return
		}
		varName := args[1]
		value, exists := os.LookupEnv(varName)
		if exists {
			fmt.Printf("%s=%s\n", varName, value)
		} else {
			fmt.Printf("%s is not set\n", varName)
		}
		return

	case "unset":

		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "unset: usage: unset VAR_NAME")
			return
		}
		varName := args[1]
		err := os.Unsetenv(varName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unset: %v\n", err)
		} else {
			fmt.Printf("Unset %s\n", varName)
		}
		return

	case "list":

		envVars := os.Environ()
		if len(envVars) == 0 {
			fmt.Println("No environment variables set")
			return
		}
		fmt.Println("\nEnvironment Variables:")
		for _, env := range envVars {
			fmt.Println(env)
		}
		return

	case "help":

		fmt.Println("\nAvailable commands:")
		fmt.Println("  cd [dir]   - Change directory")
		fmt.Println("  pwd        - Print working directory")
		fmt.Println("  clear      - Clear screen")
		fmt.Println("  set        - Set environment variable (usage: set VAR_NAME value)")
		fmt.Println("  get        - Get environment variable (usage: get VAR_NAME)")
		fmt.Println("  unset      - Unset environment variable (usage: unset VAR_NAME)")
		fmt.Println("  list       - List all environment variables")
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
