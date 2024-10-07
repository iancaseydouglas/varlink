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

