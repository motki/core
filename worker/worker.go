// Package worker provides a simple asynchronous worker queue.
package worker

import (
	"fmt"
	"sync"
	"time"
)

// Job defines the interface for performing asynchronous work.
type Job interface {
	// Perform invokes the job.
	Perform() error
}

// Scheduler is the entry-point for scheduling jobs to run asynchronously.
type Scheduler struct {
	tick    *time.Ticker
	waiting chan Job
	quit    chan struct{}
	wg      sync.WaitGroup
}

// New creates a new scheduler, ready to use.
func New() *Scheduler {
	s := &Scheduler{
		tick:    time.NewTicker(5 * time.Second),
		waiting: make(chan Job, 5),
		quit:    make(chan struct{}, 0),
		wg:      sync.WaitGroup{},
	}
	go s.Loop()
	return s
}

// Schedule adds a job to be performed.
func (s *Scheduler) Schedule(j Job) error {
	s.waiting <- j
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
					// TODO: use logger
					fmt.Println("scheduler: job returned error:", err.Error())
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
