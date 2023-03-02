package export

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/timelessnesses/l4d2-fmsa/firewall"
)

type format struct {
	BannedIPs []IP
}

type IP struct {
	IP          string `json:IP,string`
	Type_banned string `json:type_banned,string`
}

func ExportJSON(path string) error {
	convert := format{}
	for _, ip := range firewall.GetFirewallIPs().IPs {
		convert.BannedIPs = append(convert.BannedIPs, IP{IP: ip.IP, Type_banned: ip.Type_banned})
	}
	final, err := json.Marshal(convert)
	fmt.Println(final)
	if err != nil {
		return (err)
	}
	f, err := os.Create(path)
	if err != nil {
		return (err)
	}
	f.Write(final)
	return nil
}
