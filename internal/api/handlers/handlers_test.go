package handlers

import (
	"testing"
	"time"
)

// Whoever stops the process waits on WaitBackground for the work the handlers
// started: it must hold until that work ends, and not a moment longer.
func TestWaitBackgroundWaitsForTheWorkItStarted(t *testing.T) {
	h := &Handler{}
	release := make(chan struct{})
	h.startBackground("a blocked job", func() { <-release })

	expectWaiting(t, h.WaitBackground, release)
}

// The shutdown waits through CloseAndWait the same way.
func TestCloseAndWaitWaitsForTheWorkItStarted(t *testing.T) {
	h := &Handler{}
	release := make(chan struct{})
	h.startBackground("a blocked job", func() { <-release })

	expectWaiting(t, h.CloseAndWait, release)
}

// A request still running when the server gave up waiting for it must not
// start a job behind the wait: the job is dropped, never run.
func TestCloseAndWaitRefusesLaterWork(t *testing.T) {
	h := &Handler{}
	h.CloseAndWait()

	ran := make(chan struct{}, 1)
	h.startBackground("a job asked for too late", func() { ran <- struct{}{} })
	h.background.Wait()

	select {
	case <-ran:
		t.Fatal("a job started after CloseAndWait")
	default:
	}
}

// expectWaiting checks that wait blocks while a job is held, and returns
// once release lets the job end.
func expectWaiting(t *testing.T, wait func(), release chan struct{}) {
	t.Helper()

	waited := make(chan struct{})
	go func() {
		wait()
		close(waited)
	}()

	select {
	case <-waited:
		t.Fatal("the wait returned while the work was still running")
	case <-time.After(50 * time.Millisecond):
	}

	close(release)
	select {
	case <-waited:
	case <-time.After(5 * time.Second):
		t.Fatal("the wait still blocks after the work ended")
	}
}
