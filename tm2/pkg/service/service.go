package service

import (
	"errors"
	"fmt"
	"log/slog"
	"sync/atomic"

	"github.com/gnolang/gno/tm2/pkg/log"
)

var (
	// ErrAlreadyStarted is returned when somebody tries to start an already
	// running service.
	ErrAlreadyStarted = errors.New("already started")
	// ErrAlreadyStopped is returned when somebody tries to stop an already
	// stopped service (without resetting it).
	ErrAlreadyStopped = errors.New("already stopped")
	// ErrNotStarted is returned when somebody tries to stop a not running
	// service.
	ErrNotStarted = errors.New("not started")
)

// Service defines a service that can be started, stopped, and reset.
type Service interface {
	// Start the service.
	// If it's already started or stopped, will return an error.
	// If OnStart() returns an error, it's returned by Start()
	Start() error
	OnStart() error

	// Stop the service.
	// If it's already stopped, will return an error.
	// OnStop must never error.
	Stop() error
	OnStop()

	// Reset the service.
	// Panics by default - must be overwritten to enable reset.
	Reset() error
	OnReset() error

	// Return true if the service is running
	IsRunning() bool

	// Quit returns a channel, which is closed once service is stopped.
	Quit() <-chan struct{}

	// String representation of the service
	String() string

	// SetLogger sets a logger.
	SetLogger(*slog.Logger)
}

/*
Classical-inheritance-style service declarations. Services can be started, then
stopped, then optionally restarted.

Users can override the OnStart/OnStop methods. In the absence of errors, these
methods are guaranteed to be called at most once. If OnStart returns an error,
service won't be marked as started, so the user can call Start again.

A call to Reset will panic, unless OnReset is overwritten, allowing
OnStart/OnStop to be called again.

The caller must ensure that Start and Stop are not called concurrently.

It is ok to call Stop without calling Start first.

Typical usage:

	type FooService struct {
		BaseService
		// private fields
	}

	func NewFooService() *FooService {
		fs := &FooService{
			// init
		}
		fs.BaseService = *NewBaseService(log, "FooService", fs)
		return fs
	}

	func (fs *FooService) OnStart() error {
		fs.BaseService.OnStart() // Always call the overridden method.
		// initialize private fields
		// start subroutines, etc.
	}

	func (fs *FooService) OnStop() error {
		fs.BaseService.OnStop() // Always call the overridden method.
		// close/destroy private fields
		// stop subroutines, etc.
	}
*/
type BaseService struct {
	Logger  *slog.Logger
	name    string
	started atomic.Bool
	stopped atomic.Bool
	quit    chan struct{}

	// The "subclass" of BaseService
	impl Service
}

// NewBaseService creates a new BaseService.
func NewBaseService(logger *slog.Logger, name string, impl Service) *BaseService {
	if logger == nil {
		logger = log.NewNoopLogger()
	}

	return &BaseService{
		Logger: logger,
		name:   name,
		quit:   make(chan struct{}),
		impl:   impl,
	}
}

// SetLogger implements Service by setting a logger.
func (bs *BaseService) SetLogger(l *slog.Logger) {
	bs.Logger = l
}

// Start implements Service by calling OnStart (if defined). An error will be
// returned if the service is already running or stopped. Not to start the
// stopped service, you need to call Reset.
func (bs *BaseService) Start() error {
	if !bs.started.CompareAndSwap(false, true) {
		bs.Logger.Debug(fmt.Sprintf("Not starting %v -- already started", bs.name), "impl", bs.impl)
		return ErrAlreadyStarted
	}

	if bs.stopped.Load() {
		bs.Logger.Error(fmt.Sprintf("Not starting %v -- already stopped", bs.name), "impl", bs.impl)
		// revert started flag
		bs.started.Store(false)
		return ErrAlreadyStopped
	}

	bs.Logger.Info(fmt.Sprintf("Starting %v", bs.name), "impl", bs.impl)
	err := bs.impl.OnStart()
	if err != nil {
		// revert started flag
		bs.started.Store(false)
		return err
	}
	return nil
}

// OnStart implements Service by doing nothing.
// NOTE: Do not put anything in here,
// that way users don't need to call BaseService.OnStart()
func (bs *BaseService) OnStart() error { return nil }

// Stop implements Service by calling OnStop (if defined) and closing quit
// channel. An error will be returned if the service is already stopped.
func (bs *BaseService) Stop() error {
	if !bs.stopped.CompareAndSwap(false, true) {
		bs.Logger.Debug(fmt.Sprintf("Stopping %v (ignoring: already stopped)", bs.name), "impl", bs.impl)
		return ErrAlreadyStopped
	}

	if !bs.started.Load() {
		bs.Logger.Warn(fmt.Sprintf("Not stopping %v -- have not been started yet", bs.name), "impl", bs.impl)
		// revert stopped flag
		bs.stopped.Store(false)
		return ErrNotStarted
	}

	bs.Logger.Info(fmt.Sprintf("Stopping %v", bs.name), "impl", bs.impl)
	bs.impl.OnStop()
	close(bs.quit)

	return nil
}

// OnStop implements Service by doing nothing.
// NOTE: Do not put anything in here,
// that way users don't need to call BaseService.OnStop()
func (bs *BaseService) OnStop() {}

// Reset implements Service by calling OnReset callback (if defined). An error
// will be returned if the service is running.
func (bs *BaseService) Reset() error {
	if !bs.stopped.CompareAndSwap(true, false) {
		bs.Logger.Debug(fmt.Sprintf("Can't reset %v. Not stopped", bs.name), "impl", bs.impl)
		return fmt.Errorf("can't reset running %s", bs.name)
	}

	// whether or not we've started, we can reset
	bs.started.CompareAndSwap(true, false)

	bs.quit = make(chan struct{})
	return bs.impl.OnReset()
}

// OnReset implements Service by panicking.
func (bs *BaseService) OnReset() error {
	panic("The service cannot be reset")
}

// IsRunning implements Service by returning true or false depending on the
// service's state.
func (bs *BaseService) IsRunning() bool {
	return bs.started.Load() && !bs.stopped.Load()
}

// Wait blocks until the service is stopped.
func (bs *BaseService) Wait() {
	<-bs.quit
}

// String implements Service by returning a string representation of the service.
func (bs *BaseService) String() string {
	return bs.name
}

// Quit Implements Service by returning a quit channel.
func (bs *BaseService) Quit() <-chan struct{} {
	return bs.quit
}
