# Scriptify

A simple Go CLI tool that acts as a proxy for running other CLI commands with custom shortcuts.

## Features

- Add custom command shortcuts
- Store commands persistently in `~/.scriptify.json`
- Execute stored commands with simple names
- Built-in help system

## Installation

1. Build the binary:
```bash
go build -o scriptify
```

2. Copy to `/usr/bin` for global access:
```bash
sudo cp scriptify /usr/bin/
```

## Usage

### Add a new command
```bash
scriptify add <name> <command>
```

Example:
```bash
scriptify add start-docker "sudo systemctl start docker"
scriptify add stop-docker "sudo systemctl stop docker"
scriptify add ll "ls -la"
```

### Execute a stored command
```bash
scriptify <name>
```

Example:
```bash
scriptify start-docker  # Runs: sudo systemctl start docker
```

### Show help and available commands
```bash
scriptify help
```

## Configuration

Commands are stored in `~/.scriptify.json` in the following format:
```json
{
  "commands": [
    {
      "name": "start-docker",
      "command": "sudo systemctl start docker"
    }
  ]
}
```