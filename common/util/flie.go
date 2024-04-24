package util

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func WriteToFile(filePath, content string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	// 将内容写入文件
	_, err = file.WriteString(content + "\n")
	return err
}

// 在固定位置添加内容
func AddAtCustomLocation(filePath, content, location string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	file1, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	defer file1.Close()
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(file1)
	if err != nil {
		return err
	}
	if len(bytes) <= 0 {
		_, err = file1.WriteString(content)
		return err
	}

	file2, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	defer file2.Close()
	if err != nil {
		return err
	}

	lines := strings.Split(string(bytes), "\n")

	list := make([]string, 0)

	for _, line := range lines {
		if strings.Contains(line, location) {
			list = append(list, content)
		}
		list = append(list, line)
	}

	fileContent := strings.Join(list, "\n")

	_, err = file2.WriteString(fileContent)
	return err
}

// 查看某个文件路径是否存在
func FileIsExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 读取文件内容
func ReadFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return "", err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
