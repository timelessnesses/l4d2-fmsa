package parser

type CannotParseTextError struct{}

func ParseText(path string) (*BannedIPs, *CannotReadFileError)
