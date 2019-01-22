package utils

import (
	"fmt"
	"testing"

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