package main

import (
	"context"
	"net"
	"net/http"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// jobsFunc lets a test stand in for the handler a shutdown waits on.
type jobsFunc func()

func (f jobsFunc) CloseAndWait() { f() }

// journal records the order in which the steps of a shutdown happened.
type journal struct {
	mu     sync.Mutex
	events []string
}

func (j *journal) add(event string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.events = append(j.events, event)
}

func (j *journal) list() []string {
	j.mu.Lock()
	defer j.mu.Unlock()
	return slices.Clone(j.events)
}

func listen(t *testing.T) net.Listener {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	return ln
}

// startServe runs serve in the background and returns what it will return.
func startServe(ctx context.Context, srv *http.Server, ln net.Listener, jobs backgroundJobs, timeouts stopTimeouts) <-chan error {
	served := make(chan error, 1)
	go func() { served <- serve(ctx, srv, ln, jobs, timeouts) }()
	return served
}

// get sends one request and returns its status, or the error that replaced it.
func get(ln net.Listener) <-chan error {
	answered := make(chan error, 1)
	go func() {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get("http://" + ln.Addr().String())
		if err != nil {
			answered <- err
			return
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			answered <- &unexpectedStatus{resp.StatusCode}
			return
		}
		answered <- nil
	}()
	return answered
}

type unexpectedStatus struct{ code int }

func (e *unexpectedStatus) Error() string { return http.StatusText(e.code) }

func receive[T any](t *testing.T, ch <-chan T, what string) T {
	t.Helper()

	select {
	case v := <-ch:
		return v
	case <-time.After(10 * time.Second):
		t.Fatalf("still waiting for %s", what)
		panic("unreachable")
	}
}

// A request already being served when the shutdown starts gets its answer,
// and only then are the background jobs closed: a request that confirms a
// payment starts the delivery of the order on its way out.
func TestServeLetsTheRequestInFlightFinishBeforeClosingTheJobs(t *testing.T) {
	ln := listen(t)
	var steps journal
	entered, release := make(chan struct{}), make(chan struct{})

	srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-release
		steps.add("request answered")
	})}

	ctx, shutdown := context.WithCancel(context.Background())
	served := startServe(ctx, srv, ln, jobsFunc(func() { steps.add("jobs closed") }),
		stopTimeouts{requests: 10 * time.Second, background: 10 * time.Second})

	answered := get(ln)
	receive(t, entered, "the request to reach the handler")
	shutdown()

	select {
	case err := <-served:
		t.Fatalf("serve returned (%v) while a request was still being answered", err)
	case <-time.After(100 * time.Millisecond):
	}

	close(release)
	if err := receive(t, answered, "the answer"); err != nil {
		t.Fatalf("the request in flight must be answered, got %v", err)
	}
	if err := receive(t, served, "serve to return"); err != nil {
		t.Fatalf("a clean shutdown returns no error, got %v", err)
	}

	want := []string{"request answered", "jobs closed"}
	if got := steps.list(); !slices.Equal(got, want) {
		t.Fatalf("steps = %v, want %v", got, want)
	}
}

// A request that outlives the grace period has its connection closed: its
// client gets an error rather than a 200, so it does not believe an outcome
// the server never reached. The jobs are still waited for.
func TestServeClosesTheRequestsThatOutliveTheGracePeriod(t *testing.T) {
	ln := listen(t)
	entered, release := make(chan struct{}), make(chan struct{})
	t.Cleanup(func() { close(release) })

	srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		select {
		case <-release:
		case <-r.Context().Done():
		}
	})}

	closed := make(chan struct{})
	ctx, shutdown := context.WithCancel(context.Background())
	served := startServe(ctx, srv, ln, jobsFunc(func() { close(closed) }),
		stopTimeouts{requests: 100 * time.Millisecond, background: 10 * time.Second})

	answered := get(ln)
	receive(t, entered, "the request to reach the handler")
	shutdown()

	if err := receive(t, answered, "the connection to be closed"); err == nil {
		t.Fatal("a request cut by the shutdown must not read an answer")
	}
	if err := receive(t, served, "serve to return"); err != nil {
		t.Fatalf("the jobs finished, so serve returns no error, got %v", err)
	}
	select {
	case <-closed:
	default:
		t.Fatal("the jobs were not closed")
	}
}

// A job that never ends cannot hold the process forever: serve gives up and
// says so, so that the exit is logged as an abandoned job.
func TestServeGivesUpOnAJobThatDoesNotFinish(t *testing.T) {
	stuck := make(chan struct{})
	t.Cleanup(func() { close(stuck) })

	ctx, shutdown := context.WithCancel(context.Background())
	shutdown()

	err := receive(t, startServe(ctx, &http.Server{}, listen(t), jobsFunc(func() { <-stuck }),
		stopTimeouts{requests: time.Second, background: 100 * time.Millisecond}), "serve to return")
	if err == nil || !strings.Contains(err.Error(), "still running") {
		t.Fatalf("serve must report the job it abandoned, got %v", err)
	}
}

// A listener can fail after hours of service, while a delivery is running.
// serve must report the failure, but only once that job is finished: exiting
// under it would cut it between the mint and the update of the order.
func TestServeWaitsForTheJobsWhenTheListenerFails(t *testing.T) {
	ln := listen(t)
	ln.Close()

	waited := false
	err := receive(t, startServe(context.Background(), &http.Server{}, ln, jobsFunc(func() { waited = true }),
		defaultStopTimeouts), "serve to return")
	if err == nil {
		t.Fatal("serve must report why it could not serve")
	}
	if !waited {
		t.Fatal("serve returned without waiting for the background jobs")
	}
}
