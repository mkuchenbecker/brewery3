package datasink

import (
	"context"

	"github.com/mkuchenbecker/brewery3/data/gomodel/data"
)

//go:generate mockgen -destination=./mock/mock.go github.com/mkuchenbecker/brewery3/data/datasink DataSink,FirestoreClient

type DataSink interface {
	data.DataProcessorServer
}

type FirestoreClient interface {
	Send(ctx context.Context, collection string, docName string, doc map[string]interface{}) error
	Get(ctx context.Context, collection string, id string) (RowList, error)
}

type ColValueMap map[string]interface{}

type RowList []ColValueMap
