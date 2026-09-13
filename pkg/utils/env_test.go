package utils

import (
	"sync"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	previous := v
	v = viper.New()
	envOnce.Do(func() {})
	t.Cleanup(func() {
		v = previous
		envOnce = sync.Once{}
	})
	v.Set("test.value", "foobar")
	assert.Equal(t, "foobar", GetEnv("test.value"))
	assert.Empty(t, GetEnv("test.empty"))
	assert.Equal(t, "defaultbar", GetEnvWithDefault("test.missing", "defaultbar"))
}
