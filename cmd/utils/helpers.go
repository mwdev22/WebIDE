package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mwdev22/WebIDE/cmd/storage"
)

var languageMap = map[string]string{
	"cpp": "g++",
	"py":  "python3",
	"go":  "go",
	"sh":  "./",
}

var formatterMap = map[string]string{
	"go": "gofmt",
	"py": "black",
}

var TempDir string = getFileDir("tmp")

type ProgramResult string

func GetRunCmd(extName string) string {
	if cmd, exists := languageMap[extName]; exists {
		return cmd
	}
	return ""
}

func getFileDir(dirname string) string {
	cmd := exec.Command("go", "env", "GOMOD")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	goModPath := strings.TrimSpace(string(output))

	if goModPath == "" {
		log.Fatalf("go.mod file not found")
	}

	moduleDir := filepath.Dir(goModPath)

	path := filepath.Join(moduleDir, dirname)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	return path
}

func FormatCode(file *storage.File) error {
	formatter, exists := formatterMap[file.Extension]
	if !exists {
		return nil
	}

	TempDir, err := os.MkdirTemp("", "code_format")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %s", err)
	}

	tempFilePath := filepath.Join(TempDir, file.Name)
	if err := os.WriteFile(tempFilePath, []byte(file.Content), 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %s", err)
	}

	cmd := exec.Command(formatter, tempFilePath)
	formattedOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to format code: %v\nOutput: %s", err, formattedOutput)
	}

	file.Content = string(formattedOutput)
	return nil
}

func RunCode(file *storage.File) (ProgramResult, error) {

	tempFilePath := filepath.Join(TempDir, file.Name)
	defer os.Remove(tempFilePath)
	if err := os.WriteFile(tempFilePath, []byte(file.Content), 0644); err != nil {
		return "", fmt.Errorf("failed to write temp file: %s", err)
	}

	var res []byte
	var err error

	switch file.Extension {
	case "cpp":
		cmdStr := languageMap[file.Extension]
		cmd := exec.Command(cmdStr, tempFilePath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to compile C++ code: %v\nOutput: %s", err, output)
		}
		binFilePath := filepath.Join(TempDir, file.Name+".exe")
		runBinFile := exec.Command(binFilePath)
		res, err = runBinFile.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to run code: %v", err)
		}

	case "py":
		cmdStr := languageMap[file.Extension]
		cmd := exec.Command(cmdStr, tempFilePath)
		res, err = cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to run code: %v", err)
		}

	case "sh":
		format := exec.Command("dos2unix", tempFilePath)
		_, err = format.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to convert file format: %v", err)
		}
		cmdStr := languageMap[file.Extension]
		cmd := exec.Command(cmdStr, tempFilePath)
		res, err = cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to run code: %v", err)
		}

	case "go":
		cmdStr := languageMap[file.Extension]
		fmt.Println(cmdStr, tempFilePath)

		cmd := exec.Command(cmdStr, "run", tempFilePath)
		res, err = cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to run code: %v", err)
		}
	}

	return ProgramResult(res), nil
}
