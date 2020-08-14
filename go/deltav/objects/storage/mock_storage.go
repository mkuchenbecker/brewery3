package storage

// import (
// 	"context"

// 	model "github.com/golang001/deltav/model/gomodel"
// 	"github.com/golang001/deltav/utils"
// )

// // MockReactor returns a rote ReactResponse or error. Used for testing.
// type MockStorage struct {
// 	withdrawResp model.StorageResponse
// 	addResp      model.StorageResponse
// 	err          error
// }

// // WithdrawStorage implaments the GRPC interface for a Storage.
// func (mr *MockStorage) WithdrawStorage(ctx context.Context, req *model.StorageRequest) (res *model.StorageResponse,
// 	err error) {
// 	defer utils.CapturePanic(ctx, res, &err)
// 	*res = mr.withdrawResp
// 	err = mr.err
// 	return
// }

// // WithdrawStorage implaments the GRPC interface for a Storage.
// func (mr *MockStorage) AddStorage(ctx context.Context, req *model.StorageRequest) (res *model.StorageResponse,
// 	err error) {
// 	defer utils.CapturePanic(ctx, res, &err)
// 	*res = mr.addResp
// 	err = mr.err
// 	return
// }
