package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func panics() (err error) {
	defer PanicRecover(&err)
	panic("err")
}

func noPanic() (err error) {
	defer PanicRecover(&err)
	return err
}

func TestPanics(t *testing.T) {
	t.Parallel()
	err := panics()
	assert.Equal(t, "encountered a panic: err", err.Error())

	err = noPanic()
	assert.NoError(t, err)
}
