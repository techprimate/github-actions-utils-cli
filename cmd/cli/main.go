package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"log/slog"

	"github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
	"github.com/techprimate/github-actions-utils-cli/internal/cli/cmd"
	"github.com/techprimate/github-actions-utils-cli/internal/logging"
)

// version is set at build time via ldflags
var version = "dev"

// sentryRelease is set at build time via ldflags to embed the unified release identifier
var sentryRelease string

func main() {
	ctx := context.Background()

	// Allow disabling Sentry via environment variable (useful for development/testing)
	sentryEnabled := os.Getenv("TELEMETRY_ENABLED")
	if sentryEnabled != "false" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              "https://445c4c2185068fa980b83ddbe4bf1fd7@o188824.ingest.us.sentry.io/4510306572828672",
			Debug:            false,
			Environment:      "production",
			Release:          getSentryRelease(),
			AttachStacktrace: true,
			SendDefaultPII:   true,
			SampleRate:       1.0,
			EnableLogs:       true,
			EnableTracing:    true,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "sentry.Init: %s\n", err)
		}

		// Set up Sentry logging handler to capture errors
		sentryHandler := sentryslog.Option{
			EventLevel: []slog.Level{slog.LevelError},
			LogLevel:   []slog.Level{slog.LevelWarn, slog.LevelInfo},
		}.NewSentryHandler(ctx)

		// Combine Sentry handler with terminal handler
		// This gives us error tracking while maintaining local visibility
		terminalHandler := logging.NewTerminalHandler()
		multiHandler := logging.NewMultiHandler(sentryHandler, terminalHandler)
		logger := slog.New(multiHandler)
		slog.SetDefault(logger)

		// Flush buffered events before the program terminates
		defer sentry.Flush(2 * time.Second)
	}

	// Execute CLI
	if err := cmd.Execute(); err != nil {
		// Capture error in Sentry before exiting
		sentry.CaptureException(err)
		sentry.Flush(2 * time.Second)

		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// getSentryRelease returns the release identifier for Sentry
func getSentryRelease() string {
	if sentryRelease != "" {
		return sentryRelease
	}
	return "github-actions-utils-cli@" + version
}
