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

### Upcoming Features 🚀

- [ ] Piping between commands (|)
- [ ] Input/output redirection (>, <, >>)
- [ ] Background processes (&)
- [ ] Better argument parsing (handle quotes)
- [ ] Tab completion
- [ ] Custom prompt showing current directory

## Installation

### Prerequisites
- Go 1.19 or higher

### Clone and Build
```bash
git clone https://github.com/YOUR_USERNAME/goshell.git
cd goshell
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
C:\Users\Password\goshell
myshell> cd ..
myshell> pwd
C:\Users\Password
myshell> dir
[directory listing]
myshell> help
[shows help message]
myshell> exit
```

## Project Structure
```
goshell/
├── main.go          # Main shell implementation
├── go.mod           # Go module file
├── README.md        # This file
└── .gitignore       # Git ignore file
```

## Learning Resources

This project was built to learn:
- Go programming basics
- Process execution and management
- Working with standard input/output
- Building interactive CLI applications
