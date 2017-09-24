// Package worker provides a simple asynchronous worker queue.
package worker_test

import (
	"errors"
	"testing"
	"time"

	"bytes"

	"strings"

	"github.com/motki/motki/log"
	"github.com/motki/motki/worker"
	"github.com/sirupsen/logrus"
)

var delay = 1 * time.Millisecond

func TestScheduler(t *testing.T) {
	sched := worker.NewWithTick(log.New(log.Config{Level: "fatal"}), delay)
	defer sched.Shutdown()
	i := new(int)
	err := sched.ScheduleFunc(func() error {
		*i = 42
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	j := 0
	for {
		if *i == 42 {
			return
		}
		if j >= 10 {
			t.Errorf("did not process func in time")
			return
		}
		j += 1
		time.Sleep(delay)
	}
}

func TestSchedulerError(t *testing.T) {
	logger := logrus.New()
	buf := &bytes.Buffer{}
	logger.Out = buf
	logger.Formatter = &logrus.TextFormatter{DisableColors: true, DisableTimestamp: true}
	sched := worker.NewWithTick(logger, delay)
	defer sched.Shutdown()
	err := sched.ScheduleFunc(func() error {
		return errors.New("error with running the job")
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	j := 0
	for {
		if strings.Contains(buf.String(), "error with running the job") {
			return
		}
		if j >= 10 {
			t.Errorf("did not process func in time")
			return
		}
		j += 1
		time.Sleep(delay)
	}
}
