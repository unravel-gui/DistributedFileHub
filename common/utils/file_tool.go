package utils

import (
	"os"
)

// FileExists 判断文件是否存在
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// EnsureFileExists 确保路径存在
func EnsureFileExists(filePath string) error {
	exists := FileExists(filePath)
	if !exists {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
