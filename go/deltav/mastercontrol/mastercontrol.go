package mastercontrol

import (
	"fmt"

	"github.com/golang001/deltav/detectable"
	protos "github.com/golang001/deltav/model/gomodel"
	"github.com/satori/go.uuid"
)

type MasterControl struct {
	database detectable.DetectableDatabase

	//Map of ids to their objects.'s.
	objects map[string]interface{}
}

func (mc *MasterControl) insert(key string, obj interface{}) error {
	_, ok := mc.objects[key]
	if ok {
		return fmt.Errorf("key in use")
	}
	mc.objects[key] = obj
	return nil
}

func (mc *MasterControl) Initialize(req *protos.InitializeRequest) error {
	id := uuid.NewV4()
	err := mc.insert(id.String(), req.Vessel)
	if err != nil {
		return err
	}

	return nil
}
