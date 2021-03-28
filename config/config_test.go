package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitViper(t *testing.T) {
	cfg := InitViper()
	assert.NotNil(t, cfg)
}
