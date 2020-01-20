package reactor

import (
	"context"
	"math"

	model "github.com/golang001/deltav/model/gomodel"
	"github.com/golang001/deltav/utils"
)

// FusionReactor is a reactor that implaments the GRPC Reactor interface.
/*
It works by reacting helium and deuterium and produces heat, energy.
TODO: For future realism, it should produce a small amount of gamma radiation
and should slowly degrade over time due to neutron bombardment because a
small number of deuterium-deuterium reactions should occur.
*/
type FusionReactor struct {
	heliumBuffer    model.StorageTankClient
	deuteriumBuffer model.StorageTankClient
	maxEfficency    float64 // Should never be above 0.7
}

const (
	// http://www.asi.org/adb/02/09/he3-intro.html
	tJPerKgH3                  float64 = 2157000
	tJPerKgDeut                float64 = tJPerKgH3 / 0.66
	inverseTeraJoulesPerKgH3           = 1 / tJPerKgH3
	inverseTeraJoulesPerKgDeut         = 1 / tJPerKgDeut
)

// React implaments the GRPC interface for a Reactor.
func (mr *FusionReactor) React(ctx context.Context, req *model.ReactRequest) (res *model.ReactResponse,
	err error) {
	defer utils.CapturePanic(ctx, res, &err)

	// To get the desired energy, we need to increase the reaction mass in proportion to the efficency of the reaction.
	invMaxEfficency := 1 / mr.maxEfficency

	getFuel := func(fuelTank model.StorageTankClient, desiredAmount float64) (float64, error) {
		res, err := fuelTank.WithdrawStorage(ctx,
			&model.WithdrawStorageRequest{Storage: &model.Storage{Amount: desiredAmount}},
		)
		if err != nil {
			return 0, err
		}
		// We might get less back than anticipated.
		fraction := res.Storage.Amount / desiredAmount // Expensive, consider inverse math.
		if fraction > 0.99 {                           // Correct for floating point errors.
			fraction = 1.0
		}
		return fraction, nil
	}

	h3Desired := req.DesiredEnergyTeraJoules * inverseTeraJoulesPerKgH3 * invMaxEfficency
	h3Fraction, err := getFuel(mr.heliumBuffer, h3Desired)
	if err != nil {
		return
	}

	deutDesired := req.DesiredEnergyTeraJoules * inverseTeraJoulesPerKgDeut * invMaxEfficency
	deutFraction, err := getFuel(mr.deuteriumBuffer, deutDesired)
	if err != nil {
		return
	}

	// We take the min energy fraction as the true fraction.
	energyFraction := math.Min(h3Fraction, deutFraction)
	energy := energyFraction * req.DesiredEnergyTeraJoules
	heat := energy/mr.maxEfficency - energy // A fraction of the total output is heat.

	res = &model.ReactResponse{
		Outputs: []*model.ReactorOutput{
			{Type: model.ReactorOutput_ENERGY_TJOULES, Amount: energy},
			{Type: model.ReactorOutput_HEAT_TJOULES, Amount: heat},
		},
	}
	err = nil

	return
}

/*
TODO:
Options-
 this reactor needs to be topped off (self-contained, harder to manage, needs to implemnt storage interface)
 we can say fuck the buffer and it can request fuel from the main tanks directly (chaos dependencies, maybe a good thing)
 or we can have a monolith service for the ship to make requests from (maybe doesn't scale)
*/
