package shuttle

import (
	"context"

	protos "github.com/golang001/deltav/model/gomodel"
)

/*
A shuttle is a front-end client that lives in a main process that is constantly
communicating with other elements. LocalSensors live on the client process but
remote sensors are network addresses. The network latency can be compensated for
by subracting the latency from the injected latency.
*/

type PlayerShuttle struct {
	vessel protos.Vessel

	commandCenter protos.WorldModelClient
}

func NewPlayerShuttle(vessel protos.Vessel,
	commandCenter protos.WorldModelClient) (*PlayerShuttle, error) {
	ps := &PlayerShuttle{vessel: vessel, commandCenter: commandCenter}
	err := ps.Initialize()
	return ps, err
}

func (ps *PlayerShuttle) Initialize() error {
	res, err := ps.commandCenter.Initialize(context.Background(),
		&protos.InitializeRequest{Vessel: &ps.vessel})
	if err != nil {
		return err
	}
	ps.vessel = *res.Vessel
	return nil
}
