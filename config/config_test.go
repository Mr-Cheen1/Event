package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := LoadConfig("../config.json")
	if config.NotificationTime == "" || config.Timezone == "" {
		t.Error("Expected config to be loaded with values, but got empty fields")
	}
}
