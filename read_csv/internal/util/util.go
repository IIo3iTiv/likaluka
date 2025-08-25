package util

import (
	"os"
	"path/filepath"
)

func GetExecuteFilePath() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}
