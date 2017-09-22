// Package worker provides a simple asynchronous worker queue.
package worker

import (
	"sync"
	"time"

	"github.com/motki/motki/log"
)

// Job defines the interface for performing asynchronous work.
type Job interface {
	// Perform invokes the job.
	Perform() error
}

// JobFunc allows bare functions to implement the Job interface.
type JobFunc func() error

func (j JobFunc) Perform() error {
	return j()
}

// Scheduler is the entry-point for scheduling jobs to run asynchronously.
type Scheduler struct {
	tick    *time.Ticker
	waiting chan Job
	quit    chan struct{}
	wg      sync.WaitGroup
	logger  log.Logger
}

// New creates a new scheduler, ready to use.
func New(logger log.Logger) *Scheduler {
	s := &Scheduler{
		tick:    time.NewTicker(5 * time.Second),
		waiting: make(chan Job, 5),
		quit:    make(chan struct{}, 0),
		wg:      sync.WaitGroup{},
		logger:  logger,
	}
	go s.Loop()
	return s
}

// Schedule adds a job to be performed.
func (s *Scheduler) Schedule(j Job) error {
	s.waiting <- j
	return nil
}

// ScheduleFunc is a convenience method accepting a function as a job.
func (s *Scheduler) ScheduleFunc(j func() error) error {
	s.waiting <- JobFunc(j)
	return nil
}

// Loop starts the scheduler.
func (s *Scheduler) Loop() {
	s.wg.Add(1)
	for {
		select {
		case <-s.quit:
			s.wg.Done()
			return

		case <-s.tick.C:
			select {
			case j := <-s.waiting:
				err := j.Perform()
				if err != nil {
					s.logger.Warnf("scheduler: job returned error: %s", err.Error())
				}
			default:
				// nothing to work on, wait until next tick
			}
		}
	}
}

// Shutdown performs a graceful shutdown of the scheduler.
func (s *Scheduler) Shutdown() error {
	close(s.quit)
	s.wg.Wait()
	return nil
}
