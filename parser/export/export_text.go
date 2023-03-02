package export

import (
	"os"

	"github.com/timelessnesses/l4d2-fmsa/firewall"
)

type CannotParseTextError struct{}

func ExportText(path string) error {
	j, err := os.Create(path)
	if err != nil {
		return err
	}

	h := ""

	for _, ip := range firewall.GetFirewallIPs().IPs {
		h += ip.IP + "," + ip.Type_banned + "\n"
	}
	j.Write([]byte(h))
	return nil
}
