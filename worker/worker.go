// Package worker provides a simple asynchronous worker queue.
package worker

import (
	"sync"
	"time"

	"github.com/pkg/errors"

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
	delay time.Duration // Delay between ticks.

	waiting chan Job // Jobs ready to be performed.

	// Scheduled jobs are stored as slices truncated to "delay" intervals.
	// For example, if delay is 5 * time.Millisecond, the scheduled map will
	// contain a slice of jobs for every 5 milliseconds, rounded down.
	scheduled  map[time.Time][]Job
	step       time.Time  // The current step in time.
	schedMutex sync.Mutex // Guards scheduled and step.

	quit chan struct{} // Closed when shutting down.
	done chan struct{} // Closed when finished shutting down.

	workers    int64      // Number of active worker goroutines.
	countMutex sync.Mutex // Guards workers.

	logger log.Logger

	// Acceptable time to wait before forcefully quitting when shutting down
	// gracefully.
	ShutdownTimeout time.Duration
}

// New creates a new scheduler, ready to use.
func New(logger log.Logger) *Scheduler {
	return NewWithTick(logger, 5*time.Second)
}

// NewWithTick creates a new scheduler with a tick duration of delay.
func NewWithTick(logger log.Logger, delay time.Duration) *Scheduler {
	s := &Scheduler{
		delay: delay,

		waiting: make(chan Job, 5),

		scheduled:  make(map[time.Time][]Job),
		step:       time.Now().Truncate(delay),
		schedMutex: sync.Mutex{},

		quit: make(chan struct{}, 0),
		done: make(chan struct{}, 0),

		workers:    0,
		countMutex: sync.Mutex{},

		logger: logger,

		ShutdownTimeout: 1 * time.Second,
	}
	go s.Loop()
	go s.loopSchedule()
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

// ScheduleAt adds a job to be performed at a specific time.
func (s *Scheduler) ScheduleAt(j Job, t time.Time) error {
	s.schedMutex.Lock()
	defer s.schedMutex.Unlock()

	t = t.Truncate(s.delay)
	if !t.After(s.step) {
		return errors.New("cannot schedule a job in the past")
	}
	s.scheduled[t] = append(s.scheduled[t], j)

	return nil
}

// ScheduleFuncAt is a convenience method for adding a bare func as a job.
func (s *Scheduler) ScheduleFuncAt(j func() error, t time.Time) error {
	return s.ScheduleAt(JobFunc(j), t)
}

// RepeatEvery wraps a job, rescheduling it after each successful run.
func (s *Scheduler) RepeatEvery(j Job, d time.Duration) Job {
	var res Job
	res = JobFunc(func() error {
		err := j.Perform()
		if err != nil {
			return err
		}
		return s.ScheduleAt(res, time.Now().Add(d))
	})
	return res
}

// RepeatFuncEvery is a convenience method for wrapping a bare func as a repeated job.
func (s *Scheduler) RepeatFuncEvery(j func() error, d time.Duration) Job {
	return s.RepeatEvery(JobFunc(j), d)
}

// inc atomically increments the number of workers.
func (s *Scheduler) inc() {
	s.countMutex.Lock()
	defer s.countMutex.Unlock()
	s.workers += 1
}

// dec atomically decrements the number of workers.
//
// The first time workers is reduced to 0, the done channel is closed, allowing
// the scheduler to shut down gracefully.
func (s *Scheduler) dec() {
	s.countMutex.Lock()
	defer s.countMutex.Unlock()
	s.workers -= 1
	if s.workers == 0 {
		select {
		case <-s.done:
			// Select on s.done to guard from closing the channel twice.
			return
		default:
			close(s.done)
		}
	}
}

// loopSchedule moves jobs to the waiting channel as their scheduled time is reached.
func (s *Scheduler) loopSchedule() {
	s.inc()
	defer s.dec()
	tick := time.Tick(s.delay)
	for {
		select {
		case <-s.quit:
			// Quit signal received, return.
			return
		case t := <-tick:
			// Chunk the time into intervals separated by s.delay.
			t = t.Truncate(s.delay)
			// Iterate over each step previous to t.
			for s.step.Before(t) {
				s.schedMutex.Lock()
				// Read any jobs for the current time step, if any.
				jobs := s.scheduled[s.step]
				// Delete the index in the map. If it didn't exist, it's a no-op.
				delete(s.scheduled, s.step)
				s.schedMutex.Unlock()

				// Move any jobs to the waiting channel.
				for _, j := range jobs {
					s.waiting <- j
				}

				s.schedMutex.Lock()
				// Increment the current step.
				s.step = s.step.Add(s.delay)
				s.schedMutex.Unlock()
			}
		}
	}
}

// Loop begins a worker goroutine that takes care of running any jobs.
func (s *Scheduler) Loop() {
	s.inc()
	defer s.dec()
	tick := time.Tick(s.delay)
	for {
		select {
		case <-s.quit:
			// Quit signal received, return.
			return

		case <-tick:
			select {
			case j := <-s.waiting:
				// Received a waiting job, perform the work.
				if err := j.Perform(); err != nil {
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
	// Signal workers to shutdown.
	close(s.quit)
	select {
	case <-s.done:
		// Block until workers decrements to 0 and the done channel is closed;
		// the scheduler is done shutting down, and can return.
		break

	case <-time.Tick(s.ShutdownTimeout):
		// Or until ^ duration elapses, in which case, return an error.
		return errors.New("scheduler shutdown timed out")
	}
	return nil
}
