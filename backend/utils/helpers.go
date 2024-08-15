package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"

	"github.com/mwdev22/WebIDE/backend/storage"
)

var LanguageMap = map[string]string{
	"cpp": "g++",
	"py":  "python3",
	"go":  "go run",
	"sh":  "./",
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
	if cmd, exists := LanguageMap[extName]; exists {
		return cmd
	}
	return ""
}

func RunCode(file *storage.File) ProgramResult {
	tempDir, err := os.MkdirTemp("", "code_run")
	if err != nil {
		return ProgramResult(fmt.Sprintf("Failed to create temp directory: %v", err))
	}
	defer os.RemoveAll(tempDir)

	tempFilePath := filepath.Join(tempDir, file.Name+"."+file.Extension)
	if err := os.WriteFile(tempFilePath, []byte(file.Content), 0644); err != nil {
		return ProgramResult(fmt.Sprintf("Failed to write temp file: %v", err))
	}

	switch file.Extension {
	case "cpp":
		cmdStr := LanguageMap[file.Extension]
		cmd := exec.Command(cmdStr, tempFilePath)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return ProgramResult(fmt.Sprintf("Failed to run code: %v\nOutput: %s", err, output))
		}
		binFilePath := filepath.Join(tempDir, file.Name+".exe")
		runBinFile := exec.Command("./", binFilePath)
		res, err := runBinFile.CombinedOutput()

	case "py":
	case "sh":
	case "go":
	}

	cmd := exec.Command(cmdStr, tempFilePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ProgramResult(fmt.Sprintf("Failed to run code: %v\nOutput: %s", err, output))
	}

	return ProgramResult(output)
}
