package reactor

import (
	"context"

	model "github.com/golang001/deltav/model/gomodel"
	"github.com/golang001/deltav/utils"
)

// MockReactor returns a rote ReactResponse or error. Used for testing.
type MockReactor struct {
	response model.ReactResponse
	err      error
}

// React implaments the GRPC interface for a Reactor.
func (mr *MockReactor) React(ctx context.Context, req *model.ReactRequest) (res *model.ReactResponse,
	err error) {
	defer utils.CapturePanic(ctx, res, &err)
	*res = mr.response
	err = mr.err
	return
}
