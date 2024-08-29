package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	t.Skip(("refactoring"))
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	want := "Hello, World " + in
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}

	cancel()
	if err := eg.Wait(); err != nil {
		t.Errorf("failed to wait: %v", err)
	}
}
