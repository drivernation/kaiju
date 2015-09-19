package kaiju

import (
	"sync"
	"time"
)

// Service is an interface for executing long running code asynchronously. Services should be used
// when you want to run asynchronous code indefinitely. If you need short-lived asynchronous execution, then a raw
// go-routine is probably the better choice.
type Service interface {
	// Starts the service. An optional error is returned if the service failed to start.
	Start() error

	// Stops the service. An optional error is returned if the service failed to stop.
	Stop() error

	// The main execution point for the service. This method is called by Service.Start() and acts as the
	// workhorse for the service.
	//
	// An optional error is returned if something goes wrong when running the service. A returned error will cause
	// the service to panic and execution to cease.
	Run() error
}

// ScheduledService is a Service implementation that periodically executes a task at a predefined interval.
type ScheduledService struct {
	task         func() error
	onStart      func() error
	onStop       func()
	initialDelay time.Duration
	interval     time.Duration
	channel      chan int
	waiter       sync.WaitGroup
}

// Creates a new ScheduledService instance. The service, once started, will wait for initialDelay before executing
// task for the first time. Hereafter, task will be executed every interval until either an error occurs, or Stop()
// is called explicitly. onStart() is called exactly once when the service is started. It can be used to perform any
// further initialization. onStop() is called exactly once when the service is stopped. It should be used to perform
// any cleanup that may need done.
func NewScheduledService(initalDelay, interval time.Duration, task, onStart func() error, onStop func()) *ScheduledService {
	s := new(ScheduledService)
	s.task = task
	s.onStart = onStart
	s.onStop = onStop
	s.initialDelay = initalDelay
	s.interval = interval
	s.channel = make(chan int)
	s.waiter = sync.WaitGroup{}
	return s
}

// Starts the scheduled service. A error may be returned if onStart returns an error. This function will call Run() in
// an indefinite loop every interval, after the initial delay. If Run() returns an error, the go routine will panic and
// execution of the service will cease. If this happens, onStop will still be called so that the applicaiton can
// perform any required cleanup.
func (s *ScheduledService) Start() error {
	if s.onStart != nil {
		if err := s.onStart(); err != nil {
			return err
		}
	}

	s.waiter.Add(1)
	go func() {
		if s.onStop != nil {
			defer s.onStop()
		}

		for {
			time.Sleep(s.initialDelay)
			select {
			case <-s.channel:
				s.waiter.Done()
				return
			default:
				if err := s.Run(); err != nil {
					panic(err)
				}
				time.Sleep(s.interval)
			}
		}
	}()
	return nil
}

// Runs one iteration of the service. This simply runs the task once. This function should not be called explicitly.
// Start() should be used instead.
func (s *ScheduledService) Run() error {
	return s.task()
}

// Explicitly stops the service. Note that if an iteration of the task is already underway, it will finish
// before stopping. This will cause onStop to be called.
func (s *ScheduledService) Stop() error {
	s.channel <- 1
	s.waiter.Wait()
	return nil
}

// IdleService is a service implementation that performs a task at startup and then
// simply idles until told to stop, when it will perform a task before shutting down.
type IdleService struct {
	onStart func() error
	onStop  func()
	channel chan int
	waiter  sync.WaitGroup
}

// Creates a new IdleService instance.
func NewIdleService(onStart func() error, onStop func()) *IdleService {
	s := new(IdleService)
	s.onStart = onStart
	s.onStop = onStop
	s.channel = make(chan int)
	s.waiter = sync.WaitGroup{}
	return s
}

// Starts the idle service. A error may be returned if onStart returns an error. This function will call Run() in
// an indefinite loop.
func (s *IdleService) Start() error {
	if s.onStart != nil {
		if err := s.onStart(); err != nil {
			return err
		}
	}

	s.waiter.Add(1)

	go func() {
		if s.onStop != nil {
			defer s.onStop()
		}

		for {
			select {
			case <-s.channel:
				s.waiter.Done()
				return
			default:
				s.Run()
			}
		}
	}()

	return nil
}

// IdleService's run implementation is a simple no-op.
func (s *IdleService) Run() error {
	return nil
}

// Explicitly stops the service. Note that if an iteration of the task is already underway, it will finish
// before stopping. This will cause onStop to be called.
func (s *IdleService) Stop() error {
	s.channel <- 1
	s.waiter.Wait()
	return nil
}
