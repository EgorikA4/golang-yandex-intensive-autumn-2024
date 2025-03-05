package utils

import (
    "os"
    "fmt"
)

func CheckEnvVars(requiredVars []string) error {
    for _, envVar := range requiredVars {
        if _, ok := os.LookupEnv(envVar); !ok {
            return fmt.Errorf("variable %s does not exist", envVar)
        }
    }
    return nil
}

