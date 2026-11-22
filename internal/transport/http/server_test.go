package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func testLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

func testConfig(addr string, shutdownTimeout time.Duration) Config {
	return Config{
		Addr:              addr,
		ShutDownTimeout:   shutdownTimeout,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       time.Second,
	}
}

func listenLocal(t *testing.T) net.Listener {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })

	return ln
}

func waitForReady(t *testing.T, url string, done <-chan error) {
	t.Helper()

	client := &http.Client{Timeout: 200 * time.Millisecond}
	deadline := time.Now().Add(5 * time.Second)

	for time.Now().Before(deadline) {
		select {
		case err := <-done:
			t.Fatalf("server stopped before becoming ready: %v", err)
		default:
		}

		resp, err := client.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	t.Fatal("server did not become ready")
}

func waitForDone(t *testing.T, done <-chan error) error {
	t.Helper()

	select {
	case err := <-done:
		return err
	case <-time.After(5 * time.Second):
		t.Fatal("server did not stop")
		return nil
	}
}

func TestHTTPServer_Serve_GracefulShutdown(t *testing.T) {
	t.Parallel()

	ln := listenLocal(t)
	url := "http://" + ln.Addr().String()

	h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	srv := New(h, testConfig(ln.Addr().String(), time.Second), testLogger())

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	done := make(chan error, 1)
	go func() { done <- srv.serve(ctx, ln) }()

	waitForReady(t, url, done)
	cancel()

	if err := waitForDone(t, done); err != nil {
		t.Fatalf("expected graceful shutdown, got %v", err)
	}

	if conn, err := net.DialTimeout("tcp", ln.Addr().String(), 100*time.Millisecond); err == nil {
		_ = conn.Close()
		t.Fatal("listener is still accepting connections after shutdown")
	}
}

func TestHTTPServer_Serve_ShutdownTimeoutForcesClose(t *testing.T) {
	t.Parallel()

	ln := listenLocal(t)
	url := "http://" + ln.Addr().String()

	handlerStarted := make(chan struct{})
	handlerFinished := make(chan struct{})

	h := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		close(handlerStarted)
		<-r.Context().Done()
		close(handlerFinished)
	})
	srv := New(h, testConfig(ln.Addr().String(), 100*time.Millisecond), testLogger())

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	done := make(chan error, 1)
	go func() { done <- srv.serve(ctx, ln) }()

	requestDone := make(chan struct{})
	go func() {
		defer close(requestDone)
		resp, err := http.Get(url)
		if err == nil {
			_ = resp.Body.Close()
		}
	}()

	select {
	case <-handlerStarted:
	case <-time.After(5 * time.Second):
		t.Fatal("handler did not start")
	}

	cancel()

	err := waitForDone(t, done)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected shutdown deadline error, got %v", err)
	}

	select {
	case <-handlerFinished:
	case <-time.After(5 * time.Second):
		t.Fatal("active handler was not interrupted by forced close")
	}

	select {
	case <-requestDone:
	case <-time.After(5 * time.Second):
		t.Fatal("request goroutine did not finish")
	}
}

func TestHTTPServer_Run_ListenError(t *testing.T) {
	t.Parallel()

	ln := listenLocal(t)

	srv := New(http.NotFoundHandler(), testConfig(ln.Addr().String(), time.Second), testLogger())

	err := srv.Run(context.Background())
	if err == nil {
		t.Fatal("expected listen error")
	}

	if !strings.Contains(err.Error(), "listen HTTP") {
		t.Fatalf("expected wrapped listen error, got %v", err)
	}
}

func TestHTTPServer_Serve_ServeError(t *testing.T) {
	t.Parallel()

	ln := listenLocal(t)
	addr := ln.Addr()

	if err := ln.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}

	srv := New(http.NotFoundHandler(), testConfig(addr.String(), time.Second), testLogger())

	err := srv.serve(context.Background(), ln)
	if !errors.Is(err, net.ErrClosed) {
		t.Fatalf("expected closed listener error, got %v", err)
	}
}
