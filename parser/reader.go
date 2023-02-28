package parser

import (
	"os"
)

type CannotReadFileError struct{}

func ReadFile(path string) (string, *CannotReadFileError) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", &CannotReadFileError{}
	}
	return string(data), nil
}
