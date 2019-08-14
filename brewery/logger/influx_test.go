package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mkuchenbecker/brewery3/brewery/utils"
)

func TestFake(t *testing.T) {
	// This test is simply exercising the fake code to ensure a fake can be called
	// with no setup required, and that the current default is a fake.
	fake := NewFake()
	assert.NoError(t, fake.InsertTemperature(context.Background(), utils.NewTemperaturePointSink()))
	fake2 := NewDefault()
	assert.NoError(t, fake2.InsertTemperature(context.Background(), utils.NewTemperaturePointSink()))
}
