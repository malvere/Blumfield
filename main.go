package main

import (
	"blumfield/internal/blumfield"
	"context"
	"flag"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := flag.String("config", "config", "config file without extension (example: -config=dev)")
	flag.Parse()

	logger := &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
			PadLevelText:    true,
		},
	}
	// Create Blumfield instance
	blum, err := blumfield.NewBlumfield(logger, *cfg)
	if err != nil {
		logger.Fatal("Error creating Blumfield instance: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Graceful shutdown on signal
	go func() {
		sig := <-sigChan
		logger.Info("Received signal: ", sig)
		cancel()
	}()

	if blum.Config.Settings.Daemon {
		runDaemonMode(ctx, blum, logger)
	} else {
		runOnce(ctx, blum, logger)
	}

	logger.Info("Application has shut down! Goodbye!")
}

func runDaemonMode(ctx context.Context, blum *blumfield.Blumfield, logger *logrus.Logger) {
	logger.Info("Daemon mode is enabled, running tasks perpetually with 8-hour intervals...")
	for {
		select {
		case <-ctx.Done():
			logger.Info("Shutdown requested, exiting daemon loop.")
			return
		default:
			if err := blum.Start(ctx); err != nil {
				if err == context.Canceled {
					logger.Info("Task execution interrupted by shutdown request.")
					return
				}
				logger.Error("Error during task execution: ", err)
			}

			logger.Info("Waiting for 8 hours before the next task execution...")
			sleepDuration := 8*time.Hour + time.Duration(rand.IntN(15)+2)*time.Minute
			timer := time.NewTimer(sleepDuration)

			select {
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}
				logger.Info("Shutdown requested during sleep, exiting loop.")
				return
			case <-timer.C:
				logger.Info("8 hours have passed, continuing to the next execution...")
			}
		}
	}
}

func runOnce(ctx context.Context, blum *blumfield.Blumfield, logger *logrus.Logger) {
	if err := blum.Start(ctx); err != nil {
		if err == context.Canceled {
			logger.Info("Task execution interrupted by shutdown request.")
		} else {
			logger.Error("Error during task execution: ", err)
		}
	}
}
