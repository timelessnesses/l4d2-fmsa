package parser

type CannotParseFMSAError struct{}

func ParseFMSA(path string) (*BannedIPs, *CannotReadFileError, *CannotParseFMSAError) {
	s, e, e2 := ParseJson(path)
	if e != nil {
		return nil, &CannotReadFileError{}, nil
	} else if e2 != nil {
		return nil, nil, &CannotParseFMSAError{}
	}
	return s, nil, nil
}
