package tools

import (
	"blumfield/internal/models"
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"reflect"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/sirupsen/logrus"
)

type Tools struct {
	log *logrus.Logger
}

const (
	green   = "\033[32m" // ANSI code for green
	magenta = "\033[35m" // ANSI code for magenta
	reset   = "\033[0m"  // ANSI code to reset color
)

func NewTools() *Tools {
	t := &Tools{
		log: &logrus.Logger{
			Out:   os.Stdout,
			Level: logrus.DebugLevel,
			Formatter: &logrus.TextFormatter{
				// FullTimestamp:   true,
				// TimestampFormat: "2006-01-02 15:04:05",
				ForceColors: true,
				// PadLevelText: true,
			},
		},
	}
	return t
}

func (t *Tools) LogTask(task *models.Task, status string) {
	t.log.WithFields(logrus.Fields{
		"Status:": status,
		// "Title:":  task.Title,
		"ID:":     task.ID,
		"Reward:": task.Reward,
	}).Info(fmt.Sprintf("%s[TASK] %s%s", green, reset, task.Title))
}

func (t *Tools) LogGameStatus(balance *models.BalanceResponse, status string) {
	t.log.WithFields(logrus.Fields{
		"Balance:":    balance.AvailableBalance,
		"PlayPasses:": balance.PlayPasses,
	}).Info(fmt.Sprintf("%s[GAME] %sStatus: %s", magenta, reset, status))
}

func (t *Tools) Delay(sec int) {
	time.Sleep(time.Duration(rand.IntN(2)+sec) * time.Second)
}

func (tools *Tools) DelayWithContext(ctx context.Context, duration int) bool {
	timer := time.NewTimer(time.Duration(duration) * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false // Context canceled
	case <-timer.C:
		return true // Delay completed
	}
}

func (t *Tools) GetRandomUserAgent(filepath string) (string, error) {
	if err := checkTXTFile(filepath); err != nil {
		useragent, err := saveUserAgent(filepath)
		if err != nil {
			return "", err
		}
		return useragent, err
	}

	useragent, err := os.ReadFile(filepath)
	if err != nil {
		return "nil", err
	}
	return string(useragent), err

}

func checkTXTFile(path string) error {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}
	if fileInfo.Size() == 0 {
		return fmt.Errorf("file is empty")
	}
	return nil
}

func saveUserAgent(file string) (string, error) {
	fakeua, _ := faker.GetUserAgent().UserAgent(reflect.Value{})
	if err := os.WriteFile(file, []byte(fakeua.(string)), 0644); err != nil {
		return "", err
	}
	return fakeua.(string), nil
}
