package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/yourusername/varlink/internal/config"
    "github.com/yourusername/varlink/internal/env"
    "github.com/yourusername/varlink/internal/tfvars"
)

func main() {
    cfg := config.Parse()

    if cfg.Deactivate {
        env.Deactivate()
        return
    }

    environment, tfvarsFiles, err := findEnvironmentAndVars(".", cfg.LevelsAbove, cfg.MaxSearchDepth)
    if err != nil {
        fmt.Printf("I apologize, but I encountered an issue while searching for your environment and tfvars files:\n%s\n", err)
        fmt.Println("Please ensure you're within a project structure that includes an 'environments' directory,")
        fmt.Println("or consider adjusting the search parameters if your project structure is non-standard.")
        os.Exit(1)
    }

    vars, err := tfvars.ParseFiles(tfvarsFiles)
    if err != nil {
        fmt.Printf("Error parsing tfvars files: %s\n", err)
        os.Exit(1)
    }

    if cfg.DryRun {
        fmt.Printf("Would set environment to: %s\n", environment)
        for k, v := range vars {
            fmt.Printf("Would set TF_VAR_%s=%s\n", k, v)
        }
    } else {
        env.SetVars(vars)
        fmt.Printf("Environment set to: %s\n", environment)
        fmt.Printf("Set %d environment variables\n", len(vars))
    }
}

func findEnvironmentAndVars(startDir string, levelsAbove int, maxSearchDepth int) (string, []string, error) {
    currentDir, err := filepath.Abs(startDir)
    if err != nil {
        return "", nil, fmt.Errorf("error getting absolute path: %w", err)
    }

    var environment string
    var tfvarsFiles []string
    environmentFound := false
    levelsAboveCount := 0
    searchDepth := 0

    for {
        searchDepth++
        if searchDepth > maxSearchDepth {
            return "", nil, fmt.Errorf("exceeded maximum search depth of %d directories", maxSearchDepth)
        }

        files, err := filepath.Glob(filepath.Join(currentDir, "*.tfvars"))
        if err != nil {
            return "", nil, fmt.Errorf("error searching for tfvars files: %w", err)
        }
        tfvarsFiles = append(tfvarsFiles, files...)

        parentDir := filepath.Dir(currentDir)

        if !environmentFound {
            if filepath.Base(parentDir) == "environments" {
                environment = filepath.Base(currentDir)
                environmentFound = true
            }
        } else {
            levelsAboveCount++
            if levelsAboveCount > levelsAbove {
                break
            }
        }

        if currentDir == parentDir {
            break
        }
        currentDir = parentDir
    }

    if !environmentFound && filepath.Base(currentDir) == "environments" {
        environment = filepath.Base(startDir)
        environmentFound = true
    }

    if !environmentFound {
        return "", nil, fmt.Errorf("couldn't find 'environments' directory in the path")
    }

    return environment, tfvarsFiles, nil
}

