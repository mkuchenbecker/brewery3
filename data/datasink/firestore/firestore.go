package firestore

import (
	"context"

	"github.com/mkuchenbecker/brewery3/data/datasink"
	data "github.com/mkuchenbecker/brewery3/data/gomodel"
	"github.com/mkuchenbecker/brewery3/data/logger"
	"github.com/mkuchenbecker/brewery3/data/utils"
	"github.com/pkg/errors"
)

func NewStore(collection string,
	client datasink.FirestoreClient,
	log logger.Log) datasink.DataSink {
	return &firestoreSink{
		collection: collection,
		client:     client,
		log:        log,
	}
}

type firestoreSink struct {
	collection string
	client     datasink.FirestoreClient
	log        logger.Log
}

func (s *firestoreSink) Send(ctx context.Context, in *data.DataObject) (res *data.SendResponse, err error) {
	s.log.Level(logger.Info).With("request", in).Log(ctx, "send request received")
	defer func() { s.log.LogIfError(ctx, err, "send encountered an error") }()
	defer utils.PanicRecover(&err)
	row := make(map[string]interface{})
	for key, val := range in.Fields {
		switch t := val.Value.(type) {
		case *data.Value_Bool:
			row[key] = t.Bool
		case *data.Value_Bytes:
			row[key] = t.Bytes
		case *data.Value_Double:
			row[key] = t.Double
		case *data.Value_Float:
			row[key] = t.Float
		case *data.Value_Int32:
			row[key] = t.Int32
		case *data.Value_Int64:
			row[key] = t.Int64
		case *data.Value_String_:
			row[key] = t.String_
		case *data.Value_Uint32:
			row[key] = t.Uint32
		case *data.Value_Uint64:
			row[key] = t.Uint64
		default:
			return &data.SendResponse{}, errors.New("bad data")
		}
	}
	err = s.client.Send(ctx, s.collection, in.Key, row)
	return &data.SendResponse{}, errors.Wrap(nil, "error saving data")
}

func (s *firestoreSink) Get(ctx context.Context, in *data.GetRequest) (res *data.GetResponse, err error) {
	s.log.Level(logger.Info).With("request", in).Log(ctx, "get request received")
	defer func() { s.log.LogIfError(ctx, err, "get encountered an error") }()
	defer utils.PanicRecover(&err)

	response := &data.GetResponse{Data: []*data.DataObject{}}
	lom, err := s.client.Get(ctx, s.collection, in.Key)
	if err != nil {
		return response, err
	}

	do := &data.DataObject{Fields: make(map[string]*data.Value)}

	do.Key = in.Key

	for _, mapStringInterface := range lom {
		for k, v := range mapStringInterface {
			switch x := v.(type) {
			case bool:
				do.Fields[k] = &data.Value{Value: &data.Value_Bool{Bool: x}}
			case []byte:
				do.Fields[k] = &data.Value{Value: &data.Value_Bytes{Bytes: x}}
			case float64:
				do.Fields[k] = &data.Value{Value: &data.Value_Double{Double: x}}
			case float32:
				do.Fields[k] = &data.Value{Value: &data.Value_Float{Float: x}}
			case int32:
				do.Fields[k] = &data.Value{Value: &data.Value_Int32{Int32: x}}
			case int64:
				do.Fields[k] = &data.Value{Value: &data.Value_Int64{Int64: x}}
			case string:
				do.Fields[k] = &data.Value{Value: &data.Value_String_{String_: x}}
			case uint32:
				do.Fields[k] = &data.Value{Value: &data.Value_Uint32{Uint32: x}}
			case uint64:
				do.Fields[k] = &data.Value{Value: &data.Value_Uint64{Uint64: x}}
			}
		}
		response.Data = append(response.Data, do)
	}

	return response, nil
}
