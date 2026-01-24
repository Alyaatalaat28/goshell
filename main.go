package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
)

var history []string

func main() {
	// Configure readline with history
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "myshell> ",
		HistoryFile:     ".myshell_history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing readline: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	for {
		input, err := rl.Readline()
		if err != nil {
			break
		}

		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		if input == "exit" {
			break
		}

		// Add to history
		history = append(history, input)

		executeCommand(input)
	}
}

func executeCommand(input string) {
	// Check if input contains a pipe
	if strings.Contains(input, "|") {
		executePipedCommands(input)
		return
	}

	// Execute single command
	executeSingleCommand(input, nil, true)
}

func executePipedCommands(input string) {

	commands := strings.Split(input, "|")

	for i := range commands {
		commands[i] = strings.TrimSpace(commands[i])
	}

	if len(commands) < 2 {
		fmt.Fprintln(os.Stderr, "pipe: invalid pipe syntax")
		return
	}

	// Check if all commands are external (not built-in)
	// If so, let cmd.exe handle the piping natively
	allExternal := true
	for _, cmdStr := range commands {
		args := strings.Split(cmdStr, " ")
		if isBuiltinCommand(args[0]) {
			allExternal = false
			break
		}
	}

	if allExternal {
		// Let Windows handle the pipe natively
		cmd := exec.Command("cmd", "/C", input)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		return
	}

	// Manual pipe handling for commands involving built-ins
	var previousOutput *bytes.Buffer

	for i, cmdStr := range commands {
		isLast := i == len(commands)-1

		output := executeSingleCommand(cmdStr, previousOutput, isLast)

		if output == nil && !isLast {

			return
		}

		previousOutput = output
	}
}

func isBuiltinCommand(cmd string) bool {
	builtins := []string{"cd", "pwd", "clear", "set", "get", "unset", "list", "history", "help", "exit"}
	for _, builtin := range builtins {
		if cmd == builtin {
			return true
		}
	}
	return false
}

func executeSingleCommand(input string, pipeInput *bytes.Buffer, isLastInPipe bool) *bytes.Buffer {
	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, "cd: %v\n", err)
				return nil
			}
			err = os.Chdir(homeDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cd: %v\n", err)
			}
			return nil
		}
		err := os.Chdir(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "cd: %v\n", err)
		}
		return nil

	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "pwd: %v\n", err)
			return nil
		}

		output := bytes.NewBufferString(dir + "\n")

		if pipeInput == nil && isLastInPipe {
			fmt.Println(dir)
		}
		return output

	case "clear":

		if pipeInput != nil {
			fmt.Fprintln(os.Stderr, "clear: cannot be used in a pipe")
			return nil
		}
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
		return nil

	case "set":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "set: usage: set VAR_NAME value")
			return nil
		}
		varName := args[1]
		varValue := strings.Join(args[2:], " ")
		err := os.Setenv(varName, varValue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "set: %v\n", err)
			return nil
		}

		output := bytes.NewBufferString(fmt.Sprintf("Set %s=%s\n", varName, varValue))

		if pipeInput == nil && isLastInPipe {
			fmt.Printf("Set %s=%s\n", varName, varValue)
		}
		return output

	case "get":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "get: usage: get VAR_NAME")
			return nil
		}
		varName := args[1]
		value, exists := os.LookupEnv(varName)

		var output *bytes.Buffer
		if exists {
			output = bytes.NewBufferString(fmt.Sprintf("%s=%s\n", varName, value))

			if pipeInput == nil && isLastInPipe {
				fmt.Printf("%s=%s\n", varName, value)
			}
		} else {
			output = bytes.NewBufferString(fmt.Sprintf("%s is not set\n", varName))

			if pipeInput == nil && isLastInPipe {
				fmt.Printf("%s is not set\n", varName)
			}
		}
		return output

	case "unset":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "unset: usage: unset VAR_NAME")
			return nil
		}
		varName := args[1]
		err := os.Unsetenv(varName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unset: %v\n", err)
			return nil
		}

		output := bytes.NewBufferString(fmt.Sprintf("Unset %s\n", varName))

		if pipeInput == nil && isLastInPipe {
			fmt.Printf("Unset %s\n", varName)
		}
		return output

	case "list":
		envVars := os.Environ()
		if len(envVars) == 0 {
			fmt.Fprintln(os.Stderr, "No environment variables set")
			return nil
		}

		var output bytes.Buffer
		output.WriteString("\nEnvironment Variables:\n")
		for _, env := range envVars {
			output.WriteString(env + "\n")
		}

		if pipeInput == nil && isLastInPipe {
			fmt.Print(output.String())
		}
		return &output

	case "history":
		if len(history) == 0 {
			fmt.Fprintln(os.Stderr, "No commands in history")
			return nil
		}

		var output bytes.Buffer
		output.WriteString("\nCommand History:\n")
		for i, cmd := range history {
			output.WriteString(fmt.Sprintf("%4d  %s\n", i+1, cmd))
		}

		if pipeInput == nil && isLastInPipe {
			fmt.Print(output.String())
		}
		return &output

	case "help":
		var output bytes.Buffer
		output.WriteString("\n=== MyShell Help ===\n")
		output.WriteString("\n**Built-in Commands:**\n")
		output.WriteString("  cd [dir]   - Change directory (no argument goes to home)\n")
		output.WriteString("  pwd        - Print current working directory\n")
		output.WriteString("  clear      - Clear the screen\n")
		output.WriteString("  help       - Show available commands\n")
		output.WriteString("  exit       - Exit the shell\n")
		output.WriteString("\n**Environment Variable Commands:**\n")
		output.WriteString("  set        - Set environment variable (usage: set VAR_NAME value)\n")
		output.WriteString("  get        - Get environment variable (usage: get VAR_NAME)\n")
		output.WriteString("  unset      - Unset environment variable (usage: unset VAR_NAME)\n")
		output.WriteString("  list       - List all environment variables\n")
		output.WriteString("\n**History Commands:**\n")
		output.WriteString("  history    - Show command history\n")
		output.WriteString("\n**Navigation & Shortcuts:**\n")
		output.WriteString("  ↑/↓        - Navigate command history\n")
		output.WriteString("  ←/→        - Move cursor within line\n")
		output.WriteString("  Home/End   - Jump to start/end of line\n")
		output.WriteString("  Ctrl+C     - Cancel current line\n")
		output.WriteString("\n**Piping:**\n")
		output.WriteString("  |          - Pipe output between commands (e.g., history | findstr cd)\n")
		output.WriteString("\n**External Commands:**\n")
		output.WriteString("  All other commands are executed through cmd.exe\n")

		if pipeInput == nil && isLastInPipe {
			fmt.Print(output.String())
		}
		return &output
	}

	var cmd *exec.Cmd
	if len(args) == 1 {
		cmd = exec.Command("cmd", "/C", args[0])
	} else {
		cmd = exec.Command("cmd", "/C", input)
	}

	if pipeInput != nil {
		cmd.Stdin = strings.NewReader(pipeInput.String())
	}

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return nil
	}

	if isLastInPipe {
		fmt.Print(output.String())
	}

	return &output
}
