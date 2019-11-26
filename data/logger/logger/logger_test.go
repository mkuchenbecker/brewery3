package logger

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/mkuchenbecker/brewery3/data/logger"
	mock "github.com/mkuchenbecker/brewery3/data/logger/mock"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	t.Run("Log", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockGetter := mock.NewMockLoggerGetter(mockCtrl)
		mockLog := mock.NewMockLogger(mockCtrl)
		mockLog.EXPECT().Printf(
			`message
	1: 1
	error: error`).Times(1)
		mockGetter.EXPECT().Get(logger.Error).Return(mockLog).Times(1)

		l := New(mockGetter)
		l.Level(logger.Error)
		l.With("1", 1)
		l.WithError(errors.New("error"))
		l.Log(context.Background(), "message")
	})

	t.Run("printf", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockGetter := mock.NewMockLoggerGetter(mockCtrl)
		mockLog := mock.NewMockLogger(mockCtrl)
		mockLog.EXPECT().Printf("%s%d", "1", 2).Times(1)
		mockGetter.EXPECT().Get(logger.Info).Return(mockLog).Times(1)

		l := &standardLogger{
			get:  mockGetter,
			sev:  logger.Info,
			with: make(map[string]interface{}),
		}
		l.Printf("%s%d", "1", 2)
	})

	t.Run("with no error", func(t *testing.T) {
		t.Parallel()

		l := &standardLogger{
			with: make(map[string]interface{}),
		}
		l.WithError(nil)
		assert.Equal(t, 0, len(l.with))
	})

	t.Run("log if err", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockGetter := mock.NewMockLoggerGetter(mockCtrl)
			mockLog := mock.NewMockLogger(mockCtrl)
			mockLog.EXPECT().Printf("message\n\terror: error").Times(1)
			mockGetter.EXPECT().Get(logger.Error).Return(mockLog).Times(1)

			l := &standardLogger{
				get:  mockGetter,
				sev:  logger.Info,
				with: make(map[string]interface{}),
			}
			l.LogIfError(context.Background(), errors.New("error"), "message")
		})
		t.Run("no error", func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockGetter := mock.NewMockLoggerGetter(mockCtrl)

			l := &standardLogger{
				get:  mockGetter,
				sev:  logger.Info,
				with: make(map[string]interface{}),
			}
			l.LogIfError(context.Background(), nil, "message")
		})
	})

}
