package parser

import (
	"errors"
	"strings"
)

func Parse(path string) (*BannedIPs, error) {
	if strings.Contains(path, ".txt") {
		if v, err := ParseText(path); err != nil {
			return nil, errors.New("CannotParseTextError")
		} else {
			return v, nil
		}
	} else if strings.Contains(path, ".json") {
		if v, err, err2 := ParseJson(path); err != nil {
			return nil, errors.New("CannotReadFileError")
		} else if err2 != nil {
			return nil, errors.New("CannotParseJsonError")
		} else {
			return v, nil
		}
	} else if strings.Contains(path, ".fmsa") {
		if v, err, err2 := ParseFMSA(path); err != nil {
			return nil, errors.New("CannotReadFileError")
		} else if err2 != nil {
			return nil, errors.New("CannotParseFMSAError")
		} else {
			return v, nil
		}
	} else {
		if v, err := ParseRaw(path); err != nil {
			return nil, errors.New("CannotParseRawError")
		} else {
			return v, nil
		}
	}
}
