package kaiju
import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"errors"
)

type FailService struct {
	startError error
	stopError error
	runError error
	health Health
}

func (s FailService) Start() error {
	return s.startError
}

func (s FailService) Stop() error {
	return s.stopError
}

func (s FailService) Run() error {
	return s.runError
}

func (s FailService) Health() Health {
	return s.health
}

func TestScheduledServiceSuccess(t *testing.T) {
	count := 0
	started := false
	stopped := false
	s := NewScheduledService(0, time.Second, func() error {
		count+=1
		return nil
	}, func() error {
		started = true
		return nil
	}, func() {
		stopped = true
	})

	err := s.Start()
	assert.NoError(t, err)
	assert.True(t, started)
	time.Sleep(time.Second)
	err = s.Stop()
	assert.NoError(t, err)
	assert.True(t, count > 0)
	assert.True(t, stopped)
}

func TestScheduledServiceFailedStart(t *testing.T) {
	s := NewScheduledService(0, time.Second, func() error {
		return nil
	}, func() error {
		return errors.New("blah")
	}, nil)

	err := s.Start()
	assert.Error(t, err)
}

func TestScheduledServiceFailedRun(t *testing.T) {
	s := NewScheduledService(0, time.Second, func() error {
		return errors.New("blah")
	}, nil, nil)
	err := s.Start()
	assert.NoError(t, err)
	time.Sleep(time.Second)
	assert.False(t, s.Health().Healthy)
	assert.Error(t, s.Health().Error)
}

func TestIdleServiceSuccess(t *testing.T) {
	started := false
	stopped := false
	s := NewIdleService(func() error {
		started = true
		return nil
	}, func() {
		stopped = true
	})

	err := s.Start()
	assert.NoError(t, err)
	assert.True(t, s.Health().Healthy)
	assert.Nil(t, s.Health().Error)
	assert.True(t, started)
	err = s.Stop()
	assert.NoError(t, err)
	assert.True(t, stopped)
}

func TestIdleServiceStartFailed(t *testing.T) {
	s := NewIdleService(func() error {
		return errors.New("blah")
	}, nil)

	err := s.Start()
	assert.Error(t, err)
}

func TestServiceManagerSuccess(t *testing.T) {
	sm := NewServiceManager()
	s := NewIdleService(nil, nil)
	sm.AddService(s)
	assert.Equal(t, 1, sm.Size())
	err := sm.Stop()
	assert.Equal(t, ErrServiceManagerNotStarted, err)
	err = sm.Start()
	assert.NoError(t, err)
	err = sm.Start()
	assert.Equal(t, ErrServiceManagerAlreadyStarted, err)
	err = sm.Stop()
	assert.NoError(t, err)
}

func TestServiceManagerFailStart(t *testing.T) {
	sm := NewServiceManager()
	s := &FailService{startError:errors.New("blah"), health:Health{Healthy:true}}
	sm.AddService(s)
	err := sm.Start()
	assert.Error(t, err)
}

func TestServiceManager1FailStop(t *testing.T) {
	sm := NewServiceManager()
	s := &FailService{stopError:errors.New("blah"), health:Health{Healthy:true}}
	sm.AddService(s)
	sm.Start()
	err := sm.Stop()
	assert.Error(t, err)
}