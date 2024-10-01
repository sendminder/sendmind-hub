package config_test

import (
	"sendmind-hub/pkg/config"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg := config.NewConfig()
	if cfg.PostGresAddr != "" {
		t.Fail()
	}
}
