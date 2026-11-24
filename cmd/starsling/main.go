package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danni2019/starsling-tui/internal/doctor"
	"github.com/danni2019/starsling-tui/internal/logging"
	"github.com/danni2019/starsling-tui/internal/router"
	"github.com/danni2019/starsling-tui/internal/tui"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

var doctorCollectFn = doctor.Collect

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		return runCommand(args, stdout, stderr)
	}
	return runInteractive()
}

func runCommand(args []string, stdout, stderr io.Writer) int {
	switch args[0] {
	case "doctor":
		return runDoctorCommand(args[1:], stdout, stderr)
	case "demo":
		if len(args) != 1 {
			fmt.Fprintln(stderr, "usage: starsling demo")
			return 2
		}
		if err := tui.RunDemo(); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command: %s\nusage: starsling [demo|doctor]\n", args[0])
		return 2
	}
}

func runDoctorCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		fmt.Fprintln(stderr, "usage: starsling doctor")
		return 2
	}
	report := doctorCollectFn()
	report.WriteTo(stdout)
	return report.ExitCode()
}

func runInteractive() int {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// TUI owns the terminal; suppress background slog output to avoid screen corruption.
	logger := logging.NewDiscard("INFO")
	routerState := router.NewState()
	routerServer, err := router.Start(ctx, "127.0.0.1:0", routerState, logger)
	if err != nil {
		fmt.Fprintln(os.Stderr, "start router:", err)
		return 1
	}
	defer func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer stopCancel()
		_ = routerServer.Stop(stopCtx)
	}()

	return tui.Run(ctx, routerServer.Addr(), logger)
}
