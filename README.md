# GoShell 🐚

A custom shell built from scratch in Go - Learning project to understand how shells work.

## Features

### Currently Implemented 

- **Interactive Shell Loop** - Read and execute commands continuously
- **Execute External Commands** - Run any Windows command through cmd.exe
- **Built-in Commands:**
  - `cd [dir]` - Change directory (no argument goes to home)
  - `pwd` - Print current working directory
  - `clear` - Clear the screen
  - `help` - Show available commands
  - `exit` - Exit the shell
- **Environment Variable Management:**
  - `set VAR_NAME value` - Set environment variables
  - `get VAR_NAME` - Retrieve environment variable values
  - `unset VAR_NAME` - Remove environment variables
  - `list` - Display all environment variables
- **Command History:**
  - `history` - Show command history
  - Arrow key navigation (↑/↓) through previous commands
  - Persistent history across sessions (saved to `.myshell_history`)
  - Line editing with ←/→, Home/End keys
  - Ctrl+C to cancel current line

### Upcoming Features 🚀

- [ ] Piping between commands (|)
- [ ] Input/output redirection (>, <, >>)
- [ ] Background processes (&)
- [ ] Better argument parsing (handle quotes)
- [ ] Tab completion
- [ ] Custom prompt showing current directory
- [ ] Cross-platform support (Linux/Mac)
- [ ] Script execution mode
- [ ] Job control

## Installation

### Prerequisites
- Go 1.21 or higher
- Git

### Clone and Build
```bash
git clone https://github.com/YOUR_USERNAME/goshell.git
cd goshell
go mod tidy
go build -o goshell.exe
```

## Usage

### Run directly
```bash
go run main.go
```

### Or build and run
```bash
go build -o goshell.exe
goshell.exe
```

### Example Session
```
myshell> pwd
C:\Users\Username\goshell

myshell> cd ..
myshell> pwd
C:\Users\Username

myshell> set MY_VAR hello world
Set MY_VAR=hello world

myshell> get MY_VAR
MY_VAR=hello world

myshell> dir
[directory listing]

myshell> # Press ↑ to navigate to previous commands
myshell> dir

myshell> history

Command History:
   1  pwd
   2  cd ..
   3  pwd
   4  set MY_VAR hello world
   5  get MY_VAR
   6  dir
   7  history

myshell> help
[shows help message]

myshell> exit
```

## Project Structure
```
goshell/
├── main.go              # Main shell implementation
├── go.mod               # Go module file
├── go.sum               # Go dependencies checksum
├── .myshell_history     # Command history (created at runtime)
├── README.md            # This file
└── .gitignore           # Git ignore file
```

## Dependencies

- [github.com/chzyer/readline](https://github.com/chzyer/readline) - For command history and line editing

## Commands Reference

### Built-in Commands
| Command | Description | Example |
|---------|-------------|---------|
| `cd [dir]` | Change directory (no argument = home) | `cd Documents` |
| `pwd` | Print working directory | `pwd` |
| `clear` | Clear the screen | `clear` |
| `help` | Show available commands | `help` |
| `exit` | Exit the shell | `exit` |

### Environment Variables
| Command | Description | Example |
|---------|-------------|---------|
| `set VAR value` | Set environment variable | `set PATH C:\bin` |
| `get VAR` | Get environment variable | `get PATH` |
| `unset VAR` | Remove environment variable | `unset MY_VAR` |
| `list` | List all environment variables | `list` |

### History
| Command | Description | Example |
|---------|-------------|---------|
| `history` | Show command history | `history` |
| `↑/↓` | Navigate command history | Press arrow keys |

## Learning Resources

This project was built to learn:
- Go programming basics
- Process execution and management
- Working with standard input/output
- Building interactive CLI applications
- Terminal/readline handling
- Environment variable manipulation
- Command history implementation

## Acknowledgments

- Inspired by Unix shells (bash, zsh)
- Built with the [readline](https://github.com/chzyer/readline) library for enhanced input handling