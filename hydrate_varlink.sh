#!/bin/bash

# Assuming you're in the project root directory

# 1. Update internal packages
cat << EOF > internal/env/env.go
package env

import (
    "fmt"
    "os"
    "strings"
)

func Deactivate() {
    for _, env := range os.Environ() {
        if strings.HasPrefix(env, "TF_VAR_") {
            parts := strings.SplitN(env, "=", 2)
            os.Unsetenv(parts[0])
        }
    }
    fmt.Println("Terraform environment variables have been deactivated.")
}

func SetVars(vars map[string]string) {
    for k, v := range vars {
        os.Setenv(fmt.Sprintf("TF_VAR_%s", k), v)
    }
}
EOF

cat << EOF > internal/tfvars/parser.go
package tfvars

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func ParseFiles(files []string) (map[string]string, error) {
    vars := make(map[string]string)

    for _, file := range files {
        fileVars, err := parseFile(file)
        if err != nil {
            return nil, fmt.Errorf("error parsing file %s: %w", file, err)
        }
        for k, v := range fileVars {
            vars[k] = v
        }
    }

    return vars, nil
}

func parseFile(file string) (map[string]string, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    vars := make(map[string]string)
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        value = strings.Trim(value, "\"")

        vars[key] = value
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return vars, nil
}
EOF

# 2. Update cmd/varlink/main.go
# (Copy the full implementation we discussed earlier)

# 3. Update internal/config/config.go
cat << EOF > internal/config/config.go
package config

import (
    "flag"
    "os"
    "strconv"
)

type Config struct {
    LevelsAbove    int
    MaxSearchDepth int
    Deactivate     bool
    DryRun         bool
}

func Parse() Config {
    cfg := Config{
        LevelsAbove:    1,
        MaxSearchDepth: 10,
    }

    flag.IntVar(&cfg.LevelsAbove, "levels", cfg.LevelsAbove, "Levels above 'environments' to search for tfvars")
    flag.IntVar(&cfg.MaxSearchDepth, "max-depth", cfg.MaxSearchDepth, "Maximum directory depth to search")
    flag.BoolVar(&cfg.Deactivate, "deactivate", false, "Deactivate Terraform environment variables")
    flag.BoolVar(&cfg.DryRun, "dry-run", false, "Show what would be done without making changes")
    flag.Parse()

    if envLevels := os.Getenv("VARLINK_LEVELS_ABOVE"); envLevels != "" {
        if l, err := strconv.Atoi(envLevels); err == nil {
            cfg.LevelsAbove = l
        }
    }
    if envMaxDepth := os.Getenv("VARLINK_MAX_SEARCH_DEPTH"); envMaxDepth != "" {
        if d, err := strconv.Atoi(envMaxDepth); err == nil {
            cfg.MaxSearchDepth = d
        }
    }

    return cfg
}
EOF

# 4. Add unit tests (create empty files for now)
touch internal/config/config_test.go
touch internal/env/env_test.go
touch internal/tfvars/parser_test.go

# 5. Update Makefile
cat << EOF >> Makefile

lint:
	golangci-lint run

fmt:
	go fmt ./...

test-structure:
	./create_test_structure.sh
EOF

# 6. Update README.md (you may want to edit this manually with more details)
cat << EOF > README.md
# Varlink

Varlink is a tool for managing Terraform variables and environments. It simplifies the process of setting up Terraform environment variables based on .tfvars files in your project structure.

## Features

- Automatically finds and parses .tfvars files
- Sets Terraform environment variables (TF_VAR_*)
- Supports multiple environments (dev, staging, prod, etc.)
- Allows specifying search depth for .tfvars files
- Provides a dry-run mode to preview actions
- Can deactivate all Terraform-related environment variables

## Usage

\`\`\`
varlink [flags]

Flags:
  -levels int
        Levels above 'environments' to search for tfvars (default 1)
  -max-depth int
        Maximum directory depth to search (default 10)
  -deactivate
        Deactivate Terraform environment variables
  -dry-run
        Show what would be done without making changes
\`\`\`

## Building

To build the project, run:

\`\`\`
make build
\`\`\`

This will create the \`varlink\` binary in the \`bin\` directory.

## Testing

To run tests, use:

\`\`\`
make test
\`\`\`

To create a test directory structure:

\`\`\`
make test-structure
\`\`\`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
EOF

# 7. Create .github directory and add a basic workflow
mkdir -p .github/workflows
cat << EOF > .github/workflows/go.yml
name: Go

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:

  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.21

    - name: Build
      run: go build -v ./...

    - name: Test
      run: go test -v ./...
EOF

echo "Additional setup for varlink project completed."

