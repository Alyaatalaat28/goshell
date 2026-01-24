package main

import (
	"bytes"
	"fmt"
	"io"
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

	redirectInfo := parseRedirections(input)

	// Check if input contains a pipe
	if strings.Contains(redirectInfo.command, "|") {
		executePipedCommands(redirectInfo)
		return
	}

	// Execute single command with redirections
	executeSingleCommand(redirectInfo.command, nil, true, redirectInfo)
}

type RedirectionInfo struct {
	command      string
	inputFile    string
	outputFile   string
	appendOutput bool
	hasInput     bool
	hasOutput    bool
}

func parseRedirections(input string) RedirectionInfo {
	info := RedirectionInfo{command: input}

	// Check for input redirection (<)
	if strings.Contains(input, "<") {
		parts := strings.SplitN(input, "<", 2)
		info.command = strings.TrimSpace(parts[0])
		remaining := strings.TrimSpace(parts[1])

		fields := strings.Fields(remaining)
		if len(fields) > 0 {
			info.inputFile = fields[0]
			info.hasInput = true
		}
	}

	// Check for output redirection (>> or >)
	if strings.Contains(info.command, ">>") {
		parts := strings.SplitN(info.command, ">>", 2)
		info.command = strings.TrimSpace(parts[0])
		info.outputFile = strings.TrimSpace(parts[1])
		info.appendOutput = true
		info.hasOutput = true
	} else if strings.Contains(info.command, ">") {
		parts := strings.SplitN(info.command, ">", 2)
		info.command = strings.TrimSpace(parts[0])
		info.outputFile = strings.TrimSpace(parts[1])
		info.appendOutput = false
		info.hasOutput = true
	}

	return info
}

func executePipedCommands(redirectInfo RedirectionInfo) {

	commands := strings.Split(redirectInfo.command, "|")

	for i := range commands {
		commands[i] = strings.TrimSpace(commands[i])
	}

	if len(commands) < 2 {
		fmt.Fprintln(os.Stderr, "pipe: invalid pipe syntax")
		return
	}

	allExternal := true
	for _, cmdStr := range commands {
		args := strings.Split(cmdStr, " ")
		if isBuiltinCommand(args[0]) {
			allExternal = false
			break
		}
	}

	if allExternal && !redirectInfo.hasInput && !redirectInfo.hasOutput {
		// Let Windows handle the pipe natively
		cmd := exec.Command("cmd", "/C", redirectInfo.command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		return
	}

	var previousOutput *bytes.Buffer

	for i, cmdStr := range commands {
		isLast := i == len(commands)-1

		var cmdRedirect RedirectionInfo
		if isLast {
			cmdRedirect = redirectInfo
			cmdRedirect.command = cmdStr
		} else {
			cmdRedirect = RedirectionInfo{command: cmdStr}
		}

		output := executeSingleCommand(cmdStr, previousOutput, isLast, cmdRedirect)

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

func executeSingleCommand(input string, pipeInput *bytes.Buffer, isLastInPipe bool, redirectInfo RedirectionInfo) *bytes.Buffer {

	args := strings.Split(input, " ")

	// Handle input redirection
	var inputReader io.Reader
	if redirectInfo.hasInput {
		file, err := os.Open(redirectInfo.inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
			return nil
		}
		defer file.Close()
		inputReader = file
	} else if pipeInput != nil {
		inputReader = strings.NewReader(pipeInput.String())
	}

	// Prepare output writer
	var outputWriter io.Writer
	var outputBuffer bytes.Buffer
	var outputFile *os.File

	if redirectInfo.hasOutput && isLastInPipe {
		var err error
		if redirectInfo.appendOutput {
			outputFile, err = os.OpenFile(redirectInfo.outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			outputFile, err = os.Create(redirectInfo.outputFile)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening output file: %v\n", err)
			return nil
		}
		defer outputFile.Close()
		outputWriter = io.MultiWriter(&outputBuffer, outputFile)
	} else {
		outputWriter = &outputBuffer
	}

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

		fmt.Fprintf(outputWriter, "%s\n", dir)

		if !redirectInfo.hasOutput && (pipeInput == nil && isLastInPipe) {
			fmt.Println(dir)
		}
		return &outputBuffer

	case "clear":
		// Clear doesn't work in pipes or with redirections
		if pipeInput != nil || redirectInfo.hasOutput {
			fmt.Fprintln(os.Stderr, "clear: cannot be used in a pipe or with redirections")
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

		fmt.Fprintf(outputWriter, "Set %s=%s\n", varName, varValue)

		// Only print to stdout if not redirected and not in middle of pipe
		if !redirectInfo.hasOutput && (pipeInput == nil && isLastInPipe) {
			fmt.Printf("Set %s=%s\n", varName, varValue)
		}
		return &outputBuffer

	case "get":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "get: usage: get VAR_NAME")
			return nil
		}
		varName := args[1]
		value, exists := os.LookupEnv(varName)

		if exists {
			fmt.Fprintf(outputWriter, "%s=%s\n", varName, value)
			// Only print to stdout if not redirected and not in middle of pipe
			if !redirectInfo.hasOutput && (pipeInput == nil && isLastInPipe) {
				fmt.Printf("%s=%s\n", varName, value)
			}
		} else {
			fmt.Fprintf(outputWriter, "%s is not set\n", varName)
			// Only print to stdout if not redirected and not in middle of pipe
			if !redirectInfo.hasOutput && (pipeInput == nil && isLastInPipe) {
				fmt.Printf("%s is not set\n", varName)
			}
		}
		return &outputBuffer

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

		fmt.Fprintf(outputWriter, "Unset %s\n", varName)

		// Only print to stdout if not redirected and not in middle of pipe
		if !redirectInfo.hasOutput && (pipeInput == nil && isLastInPipe) {
			fmt.Printf("Unset %s\n", varName)
		}
		return &outputBuffer

	case "list":
		envVars := os.Environ()
		if len(envVars) == 0 {
			fmt.Fprintln(os.Stderr, "No environment variables set")
			return nil
		}

		fmt.Fprintf(outputWriter, "\nEnvironment Variables:\n")
		for _, env := range envVars {
			fmt.Fprintf(outputWriter, "%s\n", env)
		}

		// Only print to stdout if not redirected and not in middle of pipe
		if !redirectInfo.hasOutput && (pipeInput == nil && isLastInPipe) {
			fmt.Print(outputBuffer.String())
		}
		return &outputBuffer

	case "history":
		if len(history) == 0 {
			fmt.Fprintln(os.Stderr, "No commands in history")
			return nil
		}

		fmt.Fprintf(outputWriter, "\nCommand History:\n")
		for i, cmd := range history {
			fmt.Fprintf(outputWriter, "%4d  %s\n", i+1, cmd)
		}

		// Only print to stdout if not redirected and not in middle of pipe
		if !redirectInfo.hasOutput && (pipeInput == nil && isLastInPipe) {
			fmt.Print(outputBuffer.String())
		}
		return &outputBuffer

	case "help":
		fmt.Fprintf(outputWriter, "\n=== MyShell Help ===\n")
		fmt.Fprintf(outputWriter, "\n**Built-in Commands:**\n")
		fmt.Fprintf(outputWriter, "  cd [dir]   - Change directory (no argument goes to home)\n")
		fmt.Fprintf(outputWriter, "  pwd        - Print current working directory\n")
		fmt.Fprintf(outputWriter, "  clear      - Clear the screen\n")
		fmt.Fprintf(outputWriter, "  help       - Show available commands\n")
		fmt.Fprintf(outputWriter, "  exit       - Exit the shell\n")
		fmt.Fprintf(outputWriter, "\n**Environment Variable Commands:**\n")
		fmt.Fprintf(outputWriter, "  set        - Set environment variable (usage: set VAR_NAME value)\n")
		fmt.Fprintf(outputWriter, "  get        - Get environment variable (usage: get VAR_NAME)\n")
		fmt.Fprintf(outputWriter, "  unset      - Unset environment variable (usage: unset VAR_NAME)\n")
		fmt.Fprintf(outputWriter, "  list       - List all environment variables\n")
		fmt.Fprintf(outputWriter, "\n**History Commands:**\n")
		fmt.Fprintf(outputWriter, "  history    - Show command history\n")
		fmt.Fprintf(outputWriter, "\n**Navigation & Shortcuts:**\n")
		fmt.Fprintf(outputWriter, "  ↑/↓        - Navigate command history\n")
		fmt.Fprintf(outputWriter, "  ←/→        - Move cursor within line\n")
		fmt.Fprintf(outputWriter, "  Home/End   - Jump to start/end of line\n")
		fmt.Fprintf(outputWriter, "  Ctrl+C     - Cancel current line\n")
		fmt.Fprintf(outputWriter, "\n**Piping:**\n")
		fmt.Fprintf(outputWriter, "  |          - Pipe output between commands (e.g., history | findstr cd)\n")
		fmt.Fprintf(outputWriter, "\n**Redirection:**\n")
		fmt.Fprintf(outputWriter, "  >          - Redirect output to file (e.g., history > output.txt)\n")
		fmt.Fprintf(outputWriter, "  >>         - Append output to file (e.g., pwd >> log.txt)\n")
		fmt.Fprintf(outputWriter, "  <          - Read input from file (e.g., sort < input.txt)\n")
		fmt.Fprintf(outputWriter, "\n**External Commands:**\n")
		fmt.Fprintf(outputWriter, "  All other commands are executed through cmd.exe\n")

		// Only print to stdout if not redirected and not in middle of pipe
		if !redirectInfo.hasOutput && (pipeInput == nil && isLastInPipe) {
			fmt.Print(outputBuffer.String())
		}
		return &outputBuffer
	}

	// Execute external command
	var cmd *exec.Cmd
	if len(args) == 1 {
		cmd = exec.Command("cmd", "/C", args[0])
	} else {
		cmd = exec.Command("cmd", "/C", input)
	}

	// Set up input
	if inputReader != nil {
		cmd.Stdin = inputReader
	}

	// Set up output
	cmd.Stdout = outputWriter
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return nil
	}

	// If this is the last command in pipe and not redirected, print the output
	if isLastInPipe && !redirectInfo.hasOutput {
		fmt.Print(outputBuffer.String())
	}

	return &outputBuffer
}
