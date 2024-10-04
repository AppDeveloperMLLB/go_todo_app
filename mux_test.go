package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/sample/go_todo_app/config"
)

func TestNewMux(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	sut, cleanup, err := NewMux(context.Background(), cfg)
	defer cleanup()
	if err != nil {
		t.Fatal(err)
	}
	sut.ServeHTTP(w, r)
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusOK {
		t.Error("want 200 OK, got", resp.Status)
	}

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"status":"ok"}`
	if string(got) != want {
		t.Errorf("want %s, got %s", want, got)
	}
}
