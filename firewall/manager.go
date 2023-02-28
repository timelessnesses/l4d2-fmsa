package firewall

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/timelessnesses/l4d2-fmsa/database"
	"github.com/timelessnesses/l4d2-fmsa/parser"
)

//go:embed banned.fmsa
var banned []byte

var global_database database.Database

func Init() {
	_, err := os.Stat("fmsa.db")
	if err != nil {
		_global_database, err := database.CreateIfNotExistsDatabase("fmsa.db")
		if err != nil {
			panic(err)
		}
		global_database = *_global_database
		global_database.Execute(`
			CREATE TABLE IF NOT EXISTS banned(ip TEXT, reason STRING)
		`)
		add_copy_of_banned_servers()
	} else {
		_global_database, err := database.CreateIfNotExistsDatabase("fmsa.db")
		if err != nil {
			panic(err)
		}
		global_database = *_global_database
	}
	// check if fmsa.db is already exists so we don't overwrite it
}

func add_copy_of_banned_servers() {
	value, e1 := parser.ParseRawFMSA(banned)
	if e1 != nil {
		panic(e1)
	}
	for _, ip := range value.IPs {
		fmt.Println("Adding " + ip.IP + " to firewall")
		global_database.Execute("INSERT INTO banned VALUES (\"" + ip.IP + "\",\"" + ip.Type_banned + "\")")
	}
}

func GetFirewallIPs() *parser.BannedIPs {
	res, err := global_database.Fetch("SELECT * FROM banned")
	if err != nil {
		panic(err)
	}
	p := parser.BannedIPs{}

	for _, group := range res {
		p.IPs = append(p.IPs, parser.IP{IP: group[0], Type_banned: group[1]})
	}

	return &p
}

func AddFirewallIP(ip string, reason string) error {
	// This command requires elevation. What should I do?
	// Where can I just you know, send a UAC request in a middle of the program?
	// Or maybe I could try elevate it at the start
	// I mean I could manifest it but
	// golang.org/x/sys/windows does have elevated process for this specific use case
	// eh whatever I am just going to require user to elevate me first!

	cmd := exec.Command(
		"netsh advfirewall firewall add rule name=\"L4D2 FMSA " + ip + "\" dir=in action=block remoteip=" + ip + " enable=yes description=\"L4D2 FMSA " + reason + "\"",
	)
	if !is_elevated() {
		panic("You need to elevate this program first!")
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func is_elevated() bool {
	// try access PHYSICALDRIVE0 something
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func Cleanup() {
	global_database.Close()
}
