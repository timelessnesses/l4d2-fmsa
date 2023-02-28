package main

import (
	"time"

	"github.com/timelessnesses/l4d2-fmsa/parser"
	"github.com/visualfc/atk/tk"
)

type Window struct {
	*tk.Window
}

func main() {
	tk.MainLoop(func() {
		initialize()
	})
}

func initialize() {
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
			tk.LabelFrameAttrPadding(tk.Pad{20, 20}),
		),
		get_firewalled_ip(w),
		tk.NewLabel(
			w,
			"Enter new IP addresses (Separated by spaces) or enter nothing and click the button to open text file.",
		),
		box,
		button,
	)

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
				firewall.AddFirewallIP(ip)
			}
		}()

	}
}

func get_firewalled_ip(w *Window) *tk.Label {
	assemble := "IP Addresses that are currently firewalled:\n"
	for _, j := range firewall.GetFirewallIPs() {
		assemble += j.ip + "Firewalled Because: " + j.type_reason + "\n"
	}
}
