package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/timelessnesses/l4d2-fmsa/firewall"
	"github.com/timelessnesses/l4d2-fmsa/parser"
	"github.com/visualfc/atk/tk"
)

type Window struct {
	*tk.Window
}

var state bool

func main() {
	e := make(chan os.Signal, 1)
	signal.Notify(e, os.Interrupt)
	go func() {
		for range e {
			firewall.Cleanup()
			os.Exit(0)
		}
	}()
	tk.MainLoop(func() {
		initialize()
	})
}

var s *tk.Label

func initialize() {
	firewall.Init()
	w := &Window{}
	w.Window = tk.RootWindow()
	w.SetTitle("L4D2 Fuck Modded Server")
	w.SetSize(tk.Size{Width: 400, Height: 300})

	pack := tk.NewPackLayout(w, tk.SideTop)

	// Components
	box := tk.NewEntry(
		w,
		tk.EntryAttrWidth(50),
	)
	button := tk.NewButton(w, "Open Text File")
	state = true
	box.OnUpdate(func() {
		if len(box.Text()) >= 0 {
			button.SetText("Add IP(s)")
			state = false
		} else {
			button.SetText("Open Text File")
			state = true
		}
	})
	button.OnCommand(func() {
		handle(w, pack, state, box.Text())
	})
	remove_ip := tk.NewButton(w, "Remove IP(s)")
	view_banned_ips := tk.NewButton(w, "View Banned IPs")
	remove_ip.OnCommand(func() {
		remove_ip_from_firewall(w, pack, box.Text())
	})
	view_banned_ips.OnCommand(func() {
		view_banned_ips_from_firewall()
	})
	s = tk.NewLabel(w, get_firewalled_ip_text())
	pack.AddWidgets(
		tk.NewLabel(
			w,
			"Welcome to L4D2 MSF Application!",
			tk.LabelFrameAttrPadding(tk.Pad{X: 20, Y: 20}),
		),
		s,
		tk.NewLabel(
			w,
			"Enter new IP addresses (Separated by spaces) or enter nothing and click the button to open text file.",
		),
		box,
		button,
		remove_ip,
		view_banned_ips,
	)

	// might as well detect if app is exited

	w.OnClose(cleanup)

	w.ShowNormal()
}

func report(w *Window, msg string, pack *tk.PackLayout) {
	f := tk.NewLabel(
		w,
		"Error: "+msg,
	)
	pack.AddWidget(
		f,
	)
	// delete the label after 5 seconds
	go func() {
		time.Sleep(5 * 1000)
		f.Destroy()
	}()
}

func handle(w *Window, pack *tk.PackLayout, state bool, text_box string) {
	supported_exts := []tk.FileType{
		{
			Info: "Text Files",
			Ext:  "*.txt",
		},
		{
			Info: "All Files",
			Ext:  "*",
		},
		{
			Info: "FMSA Files",
			Ext:  "*.FMSA",
		},
		{
			Info: "JSON Files",
			Ext:  "*.json",
		},
	}
	if state {
		// Open file dialog
		path, err := tk.GetOpenFile(
			w,
			"Open a file contain IP addresses",
			supported_exts,
			"",
			"",
		)
		if err != nil {
			report(w, err.Error(), pack)
			return
		}
		println(path)
		if len(strings.Trim(path, " ")) <= 0 {
			report(w, errors.New("PathEmpty").Error(), pack)
			return
		}
		res, err := parser.Parse(path)
		if err != nil {
			report(w, err.Error(), pack)
			return
		}

		go func() {
			firewall.AddIPs(res.IPs)
		}()

		s.SetText(get_firewalled_ip_text())

	} else {
		_, err := parser.ParseRaw(text_box)
		if err != nil {
			report(w, errors.New("CannotParseRawDataError").Error(), pack)
		}
	}
}

func get_firewalled_ip_text() string {
	assemble := "IP Addresses that are currently firewalled:\n"
	for _, j := range firewall.GetFirewallIPs().IPs[0:4] {
		assemble += j.IP + " Firewalled Because: " + j.Type_banned + "\n"
	}
	if len(firewall.GetFirewallIPs().IPs[:4]) >= 0 {
		assemble += "And " + fmt.Sprint(len(firewall.GetFirewallIPs().IPs)-5) + " More IPs! Please view all of those with \"Reveal banned IPs\" button"
	}
	return assemble
}

func cleanup() bool {
	firewall.Cleanup()
	return true
}

func view_banned_ips_from_firewall() {
	panic("Not implemented")
}

func remove_ip_from_firewall(w *Window, pack *tk.PackLayout, ip string) {
	panic("Not implemented")
}
