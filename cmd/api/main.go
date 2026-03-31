package main

import (
	"dominus-broker/internal/boostrap"

	"flag"
	"fmt"
	"os"
)

var (
	mode       *bool
	showBanner *bool
	banner     = `
==========================================================
██████    ██████  ███     ██ ██ ███    ██ ██    ██ ███████
██   ██  ██    ██ ████  ████ ██ ████   ██ ██    ██ ██
██   ██  ██    ██ ██ ████ ██ ██ ██ ██  ██ ██    ██ ███████
██   ██  ██    ██ ██  ██  ██ ██ ██  ██ ██ ██    ██      ██
██████    ██████  ██      ██ ██ ██   ████  ██████  ███████
==========================================================

👉 Github: https://github.com/MBI-88
🔧 Press CTRL+C to terminate the server

dominus server is running on`
)

// catches inital variables
func init() {
	mode = flag.Bool("prod", false, "set operation mode")
	showBanner = flag.Bool("banner", true, "show banner")

	flag.Usage = func() {
		info := "[*] ***dominus-broker*** [*]\n"
		info += "mode: boolean\n"
		info += "banner: boolean\n"
		fmt.Fprintf(os.Stderr, "%s\n", info)
		flag.PrintDefaults()
	}
}

func main() {
	//Receives commands from cli
	flag.Parse()
	boostrap.RunApp(mode, showBanner, banner)
}
