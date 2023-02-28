package parser

import (
	"github.com/valyala/fastjson"
)

type CannotParseJsonError struct{}

func ParseJson(path string) (*BannedIPs, *CannotReadFileError, *CannotParseJsonError) {
	d, e1 := ReadFile(path)
	if e1 != nil {
		return nil, e1, nil
	}
	j, e2 := _ParseJson([]byte(d))
	if e2 != nil {
		return nil, nil, e2
	}
	return j, nil, nil
}

func _ParseJson(data []byte) (*BannedIPs, *CannotParseJsonError) {
	var p BannedIPs
	pv := fastjson.Parser{}
	v, err := pv.ParseBytes(data)
	if err != nil {
		// get what kind of error
		return nil, &CannotParseJsonError{}
	}
	for _, ip := range v.GetArray("BannedIPs") {
		p.IPs = append(p.IPs, IP{
			IP:          string(ip.GetStringBytes("IP")),
			type_banned: string(ip.GetStringBytes("type_banned")),
		})
	}
	return &p, nil
}
