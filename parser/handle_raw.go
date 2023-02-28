package parser

import "strings"

type CannotParseRawDataError struct{}

func ParseRaw(data string) (*BannedIPs, *CannotParseRawDataError) {
	reason := "User Added"
	var ips []IP
	for _, ip := range strings.Split(data, ",") {
		ips = append(ips, IP{ip, reason})
	}
	return &BannedIPs{ips}, nil
}
