package main

import (
	"dominus-project/internal/boostrap"

	"flag"
	"fmt"
	"os"
)

var (
	mode       *bool
	st         chan os.Signal
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
		info := "[*] ***dominus-project*** [*]\n"
		info += "mode: boolean\n"
		info += "banner: boolean\n"
		fmt.Fprintf(os.Stderr, "%s\n", info)
		flag.PrintDefaults()
	}
}

// @title dominus-project server
// @description <h3>This server is a bidirectional queue using gRCP. Manager section</h3>
// @contact.name Maikel Barrios
// @contact.email ingmbi8807@gmail.com
// @version 1.2.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-api-key
// @host localhost:8000
// @BasePath /
func main() {
	//Receives commands from cli
	flag.Parse()
	boostrap.RunApp(mode, showBanner, banner)
}
