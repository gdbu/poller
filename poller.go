package poller

import (
	"os"
	"sync"
	"time"

	"github.com/gdbu/atoms"
	"github.com/hatchify/errors"
)

const (
	// MinimumInterval is the minimum interval which can be set
	MinimumInterval = time.Millisecond * 10
	// DefaultInterval is the default interval, which is set when the provided interval is less than MinimumInterval
	DefaultInterval = time.Millisecond * 100
)

const (
	// ErrEmptyFilename is returned when an empty filename is provided
	ErrEmptyFilename = errors.Error("filename cannot be empty")
	// ErrEmptyCallback is returned when an empty callback is provided
	ErrEmptyCallback = errors.Error("event callback cannot be empty")
)

// New will return a new instance of poller
func New(filename string, fn OnEvent) (pp *Poller, err error) {
	var p Poller
	if len(filename) == 0 {
		err = ErrEmptyFilename
		return
	}

	if fn == nil {
		err = ErrEmptyCallback
		return
	}

	p.filename = filename
	p.onEvent = fn

	// Set initial file information
	info, _ := os.Stat(p.filename)
	p.updateLast(info)

	pp = &p
	return
}

// Poller watches a file
type Poller struct {
	mu sync.RWMutex

	filename string
	onEvent  OnEvent

	last       os.FileInfo
	writeState bool

	closed atoms.Bool
}

func (p *Poller) updateLast(info os.FileInfo) {
	p.last = info
}

// poll will poll a file
func (p *Poller) poll() (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed.Get() {
		return errors.ErrIsClosed
	}

	info, _ := os.Stat(p.filename)
	defer p.updateLast(info)

	switch {
	case areBothEmpty(p.last, info):
		// Both are empty, no event needed
	case wasCreated(p.last, info):
		// Fire create event
		p.onEvent(EventCreate)
	case wasRemoved(p.last, info):
		// Fire remove event
		p.onEvent(EventRemove)

	default:
		// If we passed the provided cases, we know that both p.last and info exist.
		// At this point, EventWrite and/or EventChmod could fire. We must utilize
		// if statements for these two cases because we could have a scenario where
		// both need to fire

		if wasUpdated(p.last, info) {
			p.writeState = true
		} else if p.writeState {
			p.writeState = false
			// Fire write event
			p.onEvent(EventWrite)
		}

		if hadPermissionsChange(p.last, info) {
			// Fire chmod event
			p.onEvent(EventChmod)
		}
	}

	return
}

// Run will run the poller until closed
func (p *Poller) Run(interval time.Duration) {
	var err error
	if interval < MinimumInterval {
		interval = DefaultInterval
	}

	for {
		if err = p.poll(); err != nil {
			return
		}

		time.Sleep(interval)
	}
}

// Close will close the poller
func (p *Poller) Close() (err error) {
	if !p.closed.Set(true) {
		return errors.ErrIsClosed
	}

	return
}
