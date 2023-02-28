package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/timelessnesses/l4d2-fmsa/firewall"
	"github.com/timelessnesses/l4d2-fmsa/parser"
	"github.com/visualfc/atk/tk"
)

type Window struct {
	*tk.Window
}

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
	state := false
	box.OnUpdate(func() {
		if len(box.Text()) > 0 {
			button.SetText("Add IP(s)")
			state = false
		} else {
			button.SetText("Open Text File")
			state = true
		}
	})
	button.OnCommand(func() {
		handle(w, pack, state)
	})
	pack.AddWidgets(
		tk.NewLabel(
			w,
			"Welcome to L4D2 MSF Application!",
			tk.LabelFrameAttrPadding(tk.Pad{X: 20, Y: 20}),
		),
		get_firewalled_ip(w),
		tk.NewLabel(
			w,
			"Enter new IP addresses (Separated by spaces) or enter nothing and click the button to open text file.",
		),
		box,
		button,
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
		time.Sleep(5)
		f.Destroy()
	}()
}

func handle(w *Window, pack *tk.PackLayout, state bool) {
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
			pack.AddWidget(
				tk.NewLabel(
					w,
					"Error: "+err.Error(),
				),
			)
		}
		res, err := parser.Parse(path)
		if err != nil {
			report(w, err.Error(), pack)
		}

		go func() {
			for _, ip := range res.IPs {
				firewall.AddFirewallIP(ip.IP, ip.Type_banned)
			}
		}()

	}
}

func get_firewalled_ip(w *Window) *tk.Label {
	assemble := "IP Addresses that are currently firewalled:\n"
	for _, j := range firewall.GetFirewallIPs().IPs {
		assemble += j.IP + " Firewalled Because: " + j.Type_banned + "\n"
	}
	return tk.NewLabel(
		w,
		assemble,
	)
}

func cleanup() bool {
	firewall.Cleanup()
	return true
}
