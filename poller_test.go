package poller

import (
	"os"
	"testing"
	"time"
)

func TestPoller(t *testing.T) {
	var (
		f   *os.File
		p   *Poller
		err error
	)

	if p, err = New("", nil); err != ErrEmptyFilename {
		t.Fatalf("expected error of \"%s\" and received \"%s\"", ErrEmptyFilename, err)
	}

	if p, err = New("./.test_file", nil); err != ErrEmptyCallback {
		t.Fatalf("expected error of \"%s\" and received \"%s\"", ErrEmptyCallback, err)
	}

	var count int
	waiter := make(chan struct{}, 1)
	if p, err = New("./.test_file", func(e Event) {
		var targetEvent Event
		switch count {
		case 0:
			targetEvent = EventCreate
		case 1:
			targetEvent = EventWrite
		case 2:
			targetEvent = EventChmod
		case 3:
			targetEvent = EventRemove
			waiter <- struct{}{}
		}

		if targetEvent != e {
			t.Fatalf("invalid event, expected \"%s\" and received \"%s\"", targetEvent.String(), e.String())
		}

		count++
	}); err != nil {
		t.Fatal(err)
	}

	go p.Run(0)

	time.Sleep(300 * time.Millisecond)

	if f, err = os.Create("./.test_file"); err != nil {
		t.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)

	if _, err = f.WriteString("Hello world!"); err != nil {
		t.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)

	if err = f.Close(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)

	if err = os.Chmod("./.test_file", 0655); err != nil {
		t.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)

	if err = os.Remove("./.test_file"); err != nil {
		t.Fatal(err)
	}

	<-waiter
	close(waiter)
	p.Close()
}
