package service

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrServiceManagerAlreadyStarted error = errors.New("This ServiceManager has already been started. ServiceManagers can only be started once.")
var ErrServiceManagerNotStarted error = errors.New("This ServiceManager is not started.")

type Health struct {
	Healthy bool
	Error   error
}

// Service is an interface for executing long running code asynchronously. Services should be used
// when you want to run asynchronous code indefinitely. If you need short-lived asynchronous execution, then a raw
// go-routine is probably the better choice.
type Service interface {
	// Starts the service. An optional error is returned if the service failed to start.
	Start() error

	// Stops the service. An optional error is returned if the service failed to stop.
	Stop() error

	// The main execution point for the service. This method is called by Service#Start() and acts as the
	// workhorse for the service.
	//
	// An optional error is returned if something goes wrong when running the service. A returned error will cause
	// the service to panic and execution to cease.
	Run() error

	// Returns the health of the service.
	Health() Health
}

// ScheduledService is a Service implementation that periodically executes a task at a predefined interval.
type ScheduledService struct {
	task         func() error
	onStart      func() error
	onStop       func()
	initialDelay time.Duration
	interval     time.Duration
	health       Health
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
	s.health = Health{Healthy: true}
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
					s.health.Healthy = false
					s.health.Error = err
					return
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

func (s ScheduledService) Health() Health {
	return s.health
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

// IdleService's Run() method is a noop, so it is technically always healthy.
func (s IdleService) Health() Health {
	return Health{Healthy: true}
}

// ServiceManager is a registry of Services. It can be used to start and stop a number of services.
type ServiceManager struct {
	services []Service
	started  bool
}

// Creates a pointer to a new ServiceManager instance and initializes its list of services.
// The new ServiceManager instance will be empty after creation.
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services: []Service{},
	}
}

// Registers a new service with the ServiceManager.
func (m *ServiceManager) AddService(service Service) {
	m.services = append(m.services, service)
}

// Starts all of the currently registered services. Returns an error if the ServiceManager is already started, or
// one of the registered Services fails to start. After all services are successfully stared, a new go-routine is
// started that monitors the health of all of the services. If one is found to be unhealthy, the application panics.
func (m *ServiceManager) Start() error {
	if m.started {
		return ErrServiceManagerAlreadyStarted
	}

	m.started = true
	for _, service := range m.services {
		if err := service.Start(); err != nil {
			return err
		}
	}

	go func() {
		for {
			for _, service := range m.services {
				health := service.Health()
				if !health.Healthy {
					panic(fmt.Sprintf("Service became unhealthy: %s", health.Error))
				}
			}
		}
	}()

	return nil
}

// Stops all of the currently registered Services. Returns an error if the ServiceManager is not started.
// If one or more of the registered Services fails to stop, returned errors from the Service#Stop() methods of the
// failed Services are aggregated into a single error and returned.
func (m *ServiceManager) Stop() error {
	if !m.started {
		return ErrServiceManagerNotStarted
	}
	aggErr := NewAggregateError()
	for _, service := range m.services {
		if err := service.Stop(); err != nil {
			aggErr.AddError(err)
		}
	}

	if !aggErr.Empty() {
		return aggErr
	}

	return nil
}

// Returns the number of currently registered Services.
func (m ServiceManager) Size() int {
	return len(m.services)
}
