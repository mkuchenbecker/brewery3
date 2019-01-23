package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDeferredReturn(t *testing.T) {
	t.Parallel()
	errToReturn := fmt.Errorf("error")

	f := func() (err error) {
		defer DeferErrReturn(func() error { return errToReturn }, &err)
		return nil
	}

	err := f()
	assert.Equal(t, errToReturn, err)
}

func TestDeferredReturnExistingErr(t *testing.T) {
	t.Parallel()
	errToReturn := fmt.Errorf("error")

	f := func() (err error) {
		defer DeferErrReturn(func() error { return nil }, &err)
		return errToReturn
	}

	err := f()
	assert.Equal(t, errToReturn, err)
}

func TestDeferredReturnMultipleErrs(t *testing.T) {
	t.Parallel()
	errToReturn := fmt.Errorf("error")

	f := func() (err error) {
		defer DeferErrReturn(func() error { return errToReturn }, &err)
		return errToReturn
	}

	err := f()
	assert.NotEqual(t, errToReturn, err) // Errors get concatonated.
	assert.Error(t, err)
}

type FakePrinter struct {
	lastPrint string
}

func (fp *FakePrinter) Print(p string) {
	fp.lastPrint = p
}

func TestPrinter(t *testing.T) {
	t.Parallel()
	fp := FakePrinter{}
	LogError(&fp, fmt.Errorf("error"), "encountered an error")
	assert.Equal(t, "encountered an error : error", fp.lastPrint)
}

func TestRunLoop(t *testing.T) {
	t.Parallel()

	i := 0

	f := func() error {
		i++
		return nil
	}

	err := RunLoop(100*time.Millisecond, 10*time.Millisecond, f)
	assert.NoError(t, err)
	assert.Equal(t, 10, i)
}

func TestRunLoopError(t *testing.T) {
	t.Parallel()
	i := 0
	f := func() error {
		i++
		return fmt.Errorf("error")
	}
	err := RunLoop(100*time.Millisecond, 10*time.Millisecond, f)
	assert.NoError(t, err)
	assert.Equal(t, 10, i)
}

func TestBackgroundErrReturn(t *testing.T) {
	t.Parallel()
	fp := FakePrinter{}
	BackgroundErrReturn(&fp, func() error { return fmt.Errorf("error") })
	assert.Equal(t, "background function encountered error : error", fp.lastPrint)
}
