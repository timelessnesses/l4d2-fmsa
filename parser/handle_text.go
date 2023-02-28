package parser

import "strings"

type CannotParseTextError struct{}

func ParseText(path string) (*BannedIPs, *CannotReadFileError) {
	// kinda looks something like this
	// IP, reason
	// Might as well make reason optional and make it as IPs separated by newline
	a, e := ReadFile(path)
	if e != nil {
		return nil, e
	}
	var splitted [][]string
	for _, line := range strings.Split(a, "\n") {
		splitted = append(splitted, strings.Split(line, ","))
	}
	var b BannedIPs
	for _, line := range splitted {
		b.IPs = append(b.IPs, IP{
			IP:          line[0],
			Type_banned: line[1],
		})
	}
	return &b, nil
}
