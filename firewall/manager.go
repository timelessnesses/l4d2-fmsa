package firewall

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/timelessnesses/l4d2-fmsa/database"
	"github.com/timelessnesses/l4d2-fmsa/parser"
	"golang.org/x/sys/windows"
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

func AddIPs(ips []parser.IP) {
	required := []parser.IP{}
	for _, ip := range ips {
		// check if the ip already in the db
		res, err := global_database.Fetch("SELECT * FROM banned WHERE ip=\"" + ip.IP + "\"")
		if err != nil {
			panic(err)
		}
		if len(res) <= 0 {
			global_database.Execute("INSERT INTO banned VALUES (\"" + ip.IP + "\",\"" + ip.Type_banned + "\")")
			required = append(required, ip)
		}
	}
	// build command
	done := []string{}
	for _, ip := range required {
		done = append(done, "netsh advfirewall firewall add rule name=\"FMSA "+ip.IP+"\" dir=in action=block remoteip="+ip.IP)
	}
	j := strings.Join(done, " && ")
	err := windows.ShellExecute(0, windows.StringToUTF16Ptr("runas"), windows.StringToUTF16Ptr("cmd"), windows.StringToUTF16Ptr("/c "+j), nil, 1)
	if err != nil {
		println("Error: " + err.Error())
	}
}

// func is_elevated() bool {
// 	// try access PHYSICALDRIVE0 something
// 	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
// 	return err == nil
// }

func RemoveIPs([]parser.IP) {
	// netsh delete firewall rule something
	done := []string{}
	for _, ip := range GetFirewallIPs().IPs {
		done = append(done, "netsh advfirewall firewall delete rule name=\"FMSA "+ip.IP+"\"")
	}
	j := strings.Join(done, " && ")
	err := windows.ShellExecute(0, windows.StringToUTF16Ptr("runas"), windows.StringToUTF16Ptr("cmd"), windows.StringToUTF16Ptr("/c "+j), nil, 1)
	if err != nil {
		println("Error: " + err.Error())
	}
}

func Cleanup() {
	global_database.Close()
}
