package config_test

import (
	"fmt"
	"testing"

	"example.com/sample/go_todo_app/config"
)

func TestConfig(t *testing.T) {
	wantPort := 3333
	t.Setenv("PORT", fmt.Sprintf("%d", wantPort))

	got, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	if got.Port != wantPort {
		t.Errorf("want %d, got %d", wantPort, got.Port)
	}

	wantEnv := "dev"
	if got.Env != wantEnv {
		t.Errorf("want %s, got %s", wantEnv, got.Env)
	}
}
