package reactor

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	mock_gomodel "github.com/golang001/deltav/model/gomock"
	models "github.com/golang001/deltav/model/gomodel"
	"github.com/stretchr/testify/assert"
)

func TestReactor(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	energyDesired := float64(2157000)
	efficency := float64(0.5)

	expectedH3Request := models.WithdrawStorageRequest{Storage: &models.Storage{Amount: float64(2.00)}}
	expectedH3Response := models.WithdrawStorageResponse{
		Storage: &models.Storage{
			Amount: float64(2.00),
			Type:   models.StorageType_HELIUM3_KG,
		}}
	h3Mock := mock_gomodel.NewMockStorageTankClient(mockCtrl)
	h3Mock.EXPECT().WithdrawStorage(context.Background(),
		&expectedH3Request).Return(&expectedH3Response, nil).Times(1)

	expectedDeutRequest := models.WithdrawStorageRequest{Storage: &models.Storage{Amount: float64(1.32)}}
	expectedDeutResponse := models.WithdrawStorageResponse{
		Storage: &models.Storage{
			Amount: float64(1.32),
			Type:   models.StorageType_DEUTERIUM_KG,
		}}
	deutMock := mock_gomodel.NewMockStorageTankClient(mockCtrl)
	deutMock.EXPECT().WithdrawStorage(context.Background(),
		&expectedDeutRequest).Return(&expectedDeutResponse, nil).Times(1)

	reactor := FusionReactor{heliumBuffer: h3Mock, deuteriumBuffer: deutMock, maxEfficency: efficency}

	req := models.ReactRequest{DesiredEnergyTeraJoules: energyDesired}
	res, err := reactor.React(context.Background(), &req)
	assert.NoError(t, err)

	expectedResponse := &models.ReactResponse{
		Outputs: []*models.ReactorOutput{
			{
				Type:   models.ReactorOutput_ENERGY_TJOULES,
				Amount: energyDesired},
			{
				Type:   models.ReactorOutput_HEAT_TJOULES,
				Amount: energyDesired}, // With an efficency of 50%, half the enregy becomes heat.
		},
	}
	assert.Equal(t, expectedResponse, res)
}
