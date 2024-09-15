package main

import (
	"blumfield/internal/blumfield"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
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

	// Usage
	blum, err := blumfield.NewBlumfield(logger)
	if err != nil {
		logger.Fatal("Error creating Blumfield instance: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logger.Info("Recieved signal: ", sig)
		cancel()
	}()

	if err := blum.Start(ctx); err != nil {
		logger.Fatal("Something went wrong: ", err)
	}

	time.Sleep(2 * time.Second)
	logger.Info("Application has shut down! Goodbye!")
}
