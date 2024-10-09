
package main

import (
    "os"
)

func SetFilePermissions(filePath string, permissions os.FileMode) error {
    return os.Chmod(filePath, permissions)
}
