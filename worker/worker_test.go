// Package worker provides a simple asynchronous worker queue.
package worker_test

import (
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/motki/core/log"
	"github.com/motki/core/worker"
)

var delay = 100 * time.Microsecond

// TestScheduler tests basic functionality
func TestScheduler(t *testing.T) {
	sched := worker.NewWithTick(log.New(log.Config{Level: "fatal"}), delay)
	defer sched.Shutdown()

	var i int64
	err := sched.ScheduleFunc(func() error {
		atomic.StoreInt64(&i, 42)
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	for j := 0; atomic.LoadInt64(&i) != 42; j++ {
		if j >= 10 {
			t.Errorf("did not process func in time")
			return
		}
		time.Sleep(delay)
	}
}

// TestSchedulerDoesntPanic tests that the scheduler doesn't panic when immediately shutdown.
func TestSchedulerDoesntPanic(t *testing.T) {
	sched := worker.NewWithTick(log.New(log.Config{Level: "fatal"}), delay)
	if err := sched.Shutdown(); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

// TestSchedulerShutdownTimeout tests the shutdown timeout.
func TestSchedulerShutdownTimeout(t *testing.T) {
	sched := worker.NewWithTick(log.New(log.Config{Level: "fatal"}), delay)
	sched.ShutdownTimeout = 10 * time.Millisecond
	sched.ScheduleFunc(func() error {
		time.Sleep(1 * time.Second)
		return nil
	})
	time.Sleep(10 * time.Millisecond)
	err := sched.Shutdown()
	if err == nil || !strings.Contains(err.Error(), "shutdown timed out") {
		t.Errorf("expected error, got %s", err)
	}
}

// testBuf is a locking "buffer" to ensure the test does not race.
type testBuf struct {
	sync.Mutex
	out [][]byte
}

func (t *testBuf) Write(b []byte) (int, error) {
	t.Lock()
	defer t.Unlock()
	t.out = append(t.out, b)
	return len(b), nil
}

func (t *testBuf) String() string {
	t.Lock()
	defer t.Unlock()
	if len(t.out) == 0 {
		return ""
	}
	return string(t.out[0])
}

func (t *testBuf) Len() int {
	return len(t.out)
}

// TestSchedulerError tests that errors are returned and logged.
func TestSchedulerError(t *testing.T) {
	logger := logrus.New()
	buf := &testBuf{}
	logger.Out = buf
	logger.Formatter = &logrus.TextFormatter{DisableColors: true, DisableTimestamp: true}

	sched := worker.NewWithTick(logger, delay)
	defer sched.Shutdown()

	done := make(chan struct{}, 0)

	err := sched.ScheduleFunc(func() error {
		return errors.New("error with running the job")
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	err = sched.ScheduleFunc(func() error {
		defer func() {
			// Jobs are performed sequentially, so this ensures the previous job has
			// completed before the channel is closed.
			close(done)
		}()
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	select {
	case <-done:
		if strings.Contains(buf.String(), "error with running the job") {
			return
		}
		if buf.Len() != 1 {
			t.Errorf("expected 1 error, got %d", buf.Len())
			return
		}
		t.Errorf("expected error, got none")
	case <-time.Tick(250 * time.Millisecond):
		t.Errorf("did not process func in time")
	}
}

// TestScheduleAt tests scheduling a job at a specific time in the future.
func TestScheduleAt(t *testing.T) {
	sched := worker.NewWithTick(log.New(log.Config{Level: "fatal"}), delay)
	defer sched.Shutdown()

	done := make(chan struct{})

	err := sched.ScheduleAt(worker.JobFunc(func() error {
		close(done)
		return nil
	}), time.Now().Add(delay))
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	select {
	case <-done:
		return

	case <-time.Tick(250 * time.Millisecond):
		t.Error("did not process func in time")
	}
}

// TestRepeatEvery tests repeating jobs.
func TestRepeatEvery(t *testing.T) {
	sched := worker.NewWithTick(log.New(log.Config{Level: "fatal"}), delay)
	defer sched.Shutdown()

	q := make(chan struct{})

	err := sched.ScheduleAt(
		sched.RepeatFuncEvery(
			func() error {
				q <- struct{}{}
				return nil
			}, delay),
		time.Now().Add(delay))
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	i := 0
	for {
		if i > 5 {
			break
		}
		select {
		case <-q:
			i += 1
			continue

		case <-time.Tick(250 * time.Millisecond):
			t.Error("did not process func in time", time.Now())
			return
		}
	}
}
