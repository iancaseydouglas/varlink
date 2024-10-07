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

