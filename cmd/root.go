package cmd

import (
	"flag"
	"fmt"
	"github.com/austinlparker/dropsonde/ui"
	"github.com/austinlparker/dropsonde/ui/model"
	"os"
)

var (
	tapEndpoint = flag.String("tap-endpoint", "", "opentelemetry collector remote tap endpoint")
)

func Execute() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Tap into an OpenTelemetry Collector.\n\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n", model.HelpText)
	}
	flag.Parse()
	if *tapEndpoint == "" {
		fmt.Fprintf(os.Stderr, "Tap endpoint is required\n")
		os.Exit(1)
	}

	ui.RenderUI(*tapEndpoint)
}
