package cmd

import (
	"flag"
	"fmt"
	"github.com/austinlparker/dropsonde/ui"
	"github.com/austinlparker/dropsonde/ui/model"
	"os"
)

var (
	tapEndpoint   = flag.String("tap-endpoint", "", "opentelemetry collector remote tap endpoint")
	opAmpEndpoint = flag.String("opamp-endpoint", "", "opentelemetry collector opamp endpoint")
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
	if *opAmpEndpoint == "" {
		fmt.Fprintf(os.Stderr, "OpAMP endpoint not set; Will not be enabled.\n")
	}

	ui.RenderUI(*tapEndpoint, *opAmpEndpoint)
}
