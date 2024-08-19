package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"

	"github.com/mwdev22/WebIDE/backend/storage"
)

var languageMap = map[string]string{
	"cpp": "g++",
	"py":  "python3",
	"go":  "go run",
	"sh":  "./",
}

var formatterMap = map[string]string{
	"go": "gofmt",
	"py": "black",
}

type ProgramResult string

func CheckAndUpdate[T any, G any](payload T, entity *G) error {
	payloadValue := reflect.ValueOf(payload)
	entityValue := reflect.ValueOf(entity).Elem()

	if payloadValue.Kind() != reflect.Struct || entityValue.Kind() != reflect.Struct {
		return fmt.Errorf("invalid data types for update")
	}

	for i := 0; i < payloadValue.NumField(); i++ {
		fieldValue := payloadValue.Field(i)
		fieldName := payloadValue.Type().Field(i).Name

		if fieldValue.IsValid() && !fieldValue.IsZero() {
			entityField := entityValue.FieldByName(fieldName)
			if entityField.IsValid() && entityField.CanSet() {
				entityField.Set(fieldValue)
			}
		}
	}
	return nil
}

func GetRunCmd(extName string) string {
	if cmd, exists := languageMap[extName]; exists {
		return cmd
	}
	return ""
}

func FormatCode(file *storage.File) error {
	formatter, exists := formatterMap[file.Extension]
	if !exists {
		return nil
	}

	tempDir, err := os.MkdirTemp("", "code_format")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	tempFilePath := filepath.Join(tempDir, file.Name+"."+file.Extension)
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
	tempDir, err := os.MkdirTemp("", "code_run")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	tempFilePath := filepath.Join(tempDir, file.Name+"."+file.Extension)
	if err := os.WriteFile(tempFilePath, []byte(file.Content), 0644); err != nil {
		return "", fmt.Errorf("failed to write temp file: %s", err)
	}

	var res []byte

	switch file.Extension {
	case "cpp":
		cmdStr := languageMap[file.Extension]
		cmd := exec.Command(cmdStr, tempFilePath)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to compile C++ code: %v\nOutput: %s", err, output)
		}
		binFilePath := filepath.Join(tempDir, file.Name+".exe")
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
		_, err := format.CombinedOutput()
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
		cmd := exec.Command(cmdStr, "run", tempFilePath)
		res, err = cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to run code: %v", err)
		}
	}

	return ProgramResult(res), nil
}
