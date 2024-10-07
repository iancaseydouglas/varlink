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

