package firestore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mkuchenbecker/brewery3/data/datasink"
	mock "github.com/mkuchenbecker/brewery3/data/datasink/mock"
	"github.com/mkuchenbecker/brewery3/data/gomodel/data"
	"github.com/stretchr/testify/assert"
)

var genericID = "1234"

func TestFirestoreSink(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	id := genericID
	collection := "testCollection"

	t.Run("Send Success", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ts := time.Now()

		key := fmt.Sprintf("%s:%d", id, ts.UnixNano())

		fields := map[string]*data.Value{
			"id":     {Value: &data.Value_String_{String_: id}},
			"bool":   {Value: &data.Value_Bool{Bool: true}},
			"bytes":  {Value: &data.Value_Bytes{Bytes: []byte("1")}},
			"float":  {Value: &data.Value_Float{Float: 1}},
			"double": {Value: &data.Value_Double{Double: 1}},
			"int32":  {Value: &data.Value_Int32{Int32: 1}},
			"int64":  {Value: &data.Value_Int64{Int64: 1}},
			"uint32": {Value: &data.Value_Uint32{Uint32: 1}},
			"uint64": {Value: &data.Value_Uint64{Uint64: 1}},
			"string": {Value: &data.Value_String_{String_: "1"}},
		}
		rows := map[string]interface{}{
			"id":     id,
			"bool":   true,
			"bytes":  []byte("1"),
			"float":  float32(1),
			"double": float64(1),
			"int32":  int32(1),
			"int64":  int64(1),
			"uint32": uint32(1),
			"uint64": uint64(1),
			"string": "1",
		}

		mockFirestoreClient := mock.NewMockFirestoreClient(mockCtrl)
		mockFirestoreClient.EXPECT().Send(ctx, collection, key, rows).Return(nil).Times(1)

		store := NewStore(collection, mockFirestoreClient)

		req := &data.DataObject{Key: key, Fields: fields}
		_, err := store.Send(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("Get Success", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ts := time.Now()

		key := fmt.Sprintf("%s:%d", id, ts.UnixNano())

		fields := map[string]*data.Value{
			"id":     {Value: &data.Value_String_{String_: id}},
			"bool":   {Value: &data.Value_Bool{Bool: true}},
			"bytes":  {Value: &data.Value_Bytes{Bytes: []byte("1")}},
			"float":  {Value: &data.Value_Float{Float: 1}},
			"double": {Value: &data.Value_Double{Double: 1}},
			"int32":  {Value: &data.Value_Int32{Int32: 1}},
			"int64":  {Value: &data.Value_Int64{Int64: 1}},
			"uint32": {Value: &data.Value_Uint32{Uint32: 1}},
			"uint64": {Value: &data.Value_Uint64{Uint64: 1}},
			"string": {Value: &data.Value_String_{String_: "1"}},
		}
		rows := map[string]interface{}{
			"id":     id,
			"bool":   true,
			"bytes":  []byte("1"),
			"float":  float32(1),
			"double": float64(1),
			"int32":  int32(1),
			"int64":  int64(1),
			"uint32": uint32(1),
			"uint64": uint64(1),
			"string": "1",
		}

		mockFirestoreClient := mock.NewMockFirestoreClient(mockCtrl)
		mockFirestoreClient.EXPECT().Get(ctx, collection, key).Return([]datasink.ColValueMap{rows}, nil).Times(1)

		var store datasink.DataSink = NewStore(collection, mockFirestoreClient)

		response, err := store.Get(ctx, &data.GetRequest{Key: key})
		assert.NoError(t, err)

		expected := data.DataObject{Key: key, Fields: fields}
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, expected, *response.Data[0])
	})

}
