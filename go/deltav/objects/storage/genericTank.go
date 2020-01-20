package storage

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	model "github.com/golang001/deltav/model/gomodel"
	"github.com/golang001/deltav/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GenericTank struct {
	capacity    float64
	current     float64
	storageType model.StorageType

	mux sync.RWMutex
}

func startGasTank(gt *GenericTank, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serve := grpc.NewServer()
	model.RegisterStorageTankServer(serve, gt)
	// Register reflection service on gRPC server.
	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getGasTankClient(port int) model.StorageTankClient {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	return model.NewStorageTankClient(conn)
}

func NewGasTank(capacity float64, initialVolume float64, storageType model.StorageType, port int) model.StorageTankClient {
	go startGasTank(&GenericTank{capacity: capacity, current: initialVolume, storageType: storageType}, port)
	time.Sleep(time.Second)

	return getGasTankClient(port)
}

// WithdrawStorage implaments the GRPC interface for a Storage.
func (gt *GenericTank) WithdrawStorage(ctx context.Context, req *model.WithdrawStorageRequest) (res *model.WithdrawStorageResponse,
	err error) {
	gt.mux.Lock()
	defer gt.mux.Unlock()
	defer utils.CapturePanic(ctx, res, &err)
	if gt.current < req.Storage.Amount {
		gt.current = 0
		return &model.WithdrawStorageResponse{Storage: &model.Storage{Amount: gt.current, Type: gt.storageType}}, nil
	}
	gt.current -= req.Storage.Amount
	return &model.WithdrawStorageResponse{Storage: &model.Storage{Amount: req.Storage.Amount, Type: gt.storageType}}, nil
}

// AddStorage implaments the GRPC interface for a Storage.
func (gt *GenericTank) AddStorage(ctx context.Context, req *model.AddStorageRequest) (res *model.AddStorageResponse,
	err error) {
	gt.mux.Lock()
	defer gt.mux.Unlock()
	defer utils.CapturePanic(ctx, res, &err)
	spareCap := gt.capacity - req.Storage.Amount + gt.current
	if spareCap < 0 {
		gt.current = gt.capacity
		return &model.AddStorageResponse{Storage: &model.Storage{Amount: req.Storage.Amount + spareCap, Type: gt.storageType}}, nil
	}
	return &model.AddStorageResponse{Storage: &model.Storage{Amount: float64(0), Type: gt.storageType}}, nil
}

// Status implaments the GRPC interface for a Storage.
func (gt *GenericTank) Status(ctx context.Context, req *model.StorageStatusRequest) (res *model.StorageStatusResponse,
	err error) {
	gt.mux.RLock()
	defer gt.mux.RUnlock()
	defer utils.CapturePanic(ctx, res, &err)

	return &model.StorageStatusResponse{Capacity: gt.capacity, Current: gt.current, Type: gt.storageType}, nil
}
