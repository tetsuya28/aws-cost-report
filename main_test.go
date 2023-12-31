//go:build integration

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	err := handler()
	assert.NoError(t, err)
}
