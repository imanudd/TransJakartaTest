package usecase

import (
	"app/config"
	"app/internal/domain"
	"app/internal/repository"
	"app/pkg/rabbitMQ"
	"app/utils"
	"context"
)

type VehicleUsecaseImpl interface {
	GetHistory(ctx context.Context, vehicleId string, req *domain.GetVehicleHistoryRequest) ([]*domain.VehicleLocation, *utils.CustomError)
	GetLatestLocation(ctx context.Context, vehicleId string) (*domain.VehicleLocation, *utils.CustomError)
}

type VehicleUsecase struct {
	cfg  *config.MainConfig
	repo repository.RepositoryImpl
	rmq  rabbitMQ.RabbitMQImpl
}

func NewVehicleUsecase(cfg *config.MainConfig, repo repository.RepositoryImpl, rmq rabbitMQ.RabbitMQImpl) VehicleUsecaseImpl {
	return &VehicleUsecase{
		cfg:  cfg,
		repo: repo,
		rmq:  rmq,
	}
}

func (u *VehicleUsecase) GetHistory(ctx context.Context, vehicleId string, req *domain.GetVehicleHistoryRequest) ([]*domain.VehicleLocation, *utils.CustomError) {
	startTime := utils.GetTime(req.Start)
	endTime := utils.GetTime(req.End)

	if startTime.After(endTime) {
		return nil, utils.ErrBadRequest("start time must be before end time")
	}

	vehicleLocations, err := u.repo.GetVehicleRepo().GetHistory(ctx, vehicleId, req)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}

	if len(vehicleLocations) == 0 {
		vehicleLocations = []*domain.VehicleLocation{}
		return vehicleLocations, nil
	}

	return vehicleLocations, nil
}

func (u *VehicleUsecase) GetLatestLocation(ctx context.Context, vehicleId string) (*domain.VehicleLocation, *utils.CustomError) {
	vehicleLocation, err := u.repo.GetVehicleRepo().GetLatestLocation(ctx, vehicleId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}

	if vehicleLocation == nil {
		return nil, utils.ErrNotFound("vehicle")
	}

	return vehicleLocation, nil
}
