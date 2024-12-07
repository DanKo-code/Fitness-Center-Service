package service_usecase

import (
	"Service/internal/dtos"
	customErrors "Service/internal/errors"
	"Service/internal/models"
	"Service/internal/repository"
	"context"
	coachGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.coach"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type ServiceUseCase struct {
	serviceRepo repository.ServiceRepository
	coachClient *coachGRPC.CoachClient
}

func NewServiceUseCase(
	serviceRepo repository.ServiceRepository,
	coachClient *coachGRPC.CoachClient,
) *ServiceUseCase {
	return &ServiceUseCase{
		serviceRepo: serviceRepo,
		coachClient: coachClient,
	}
}

func (u *ServiceUseCase) CreateService(ctx context.Context, cmd *dtos.CreateServiceCommand) (*models.Service, error) {

	service := &models.Service{
		Id:          uuid.New(),
		Title:       cmd.Title,
		Photo:       cmd.Photo,
		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}

	err := u.serviceRepo.CreateService(ctx, service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (u *ServiceUseCase) GetServiceById(ctx context.Context, id uuid.UUID) (*models.Service, error) {
	service, err := u.serviceRepo.GetServiceById(ctx, id)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (u *ServiceUseCase) UpdateService(ctx context.Context, cmd *dtos.UpdateServiceCommand) (*models.Service, error) {

	err := u.serviceRepo.UpdateService(ctx, cmd)
	if err != nil {
		return nil, err
	}

	service, err := u.serviceRepo.GetServiceById(ctx, cmd.Id)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (u *ServiceUseCase) DeleteServiceById(ctx context.Context, id uuid.UUID) (*models.Service, error) {

	service, err := u.serviceRepo.GetServiceById(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.serviceRepo.DeleteService(ctx, id)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (u *ServiceUseCase) GetServices(ctx context.Context) ([]*models.Service, error) {
	services, err := u.serviceRepo.GetServices(ctx)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (u *ServiceUseCase) CreateCoachServices(ctx context.Context, cmd *dtos.CreateCoachServicesCommand) ([]*models.Service, error) {

	getCoachByIdRequest := &coachGRPC.GetCoachByIdRequest{Id: cmd.CoachId.String()}

	_, err := (*u.coachClient).GetCoachById(ctx, getCoachByIdRequest)
	if err != nil {

		st, ok := status.FromError(err)

		if !ok {
			return nil, nil
		}

		switch st.Code() {
		case codes.NotFound:
			return nil, customErrors.CoachNotFound
		default:
			return nil, customErrors.InternalCoachServerError
		}
	}

	_, err = u.serviceRepo.GetServicesByIds(ctx, cmd.ServicesIds)
	if err != nil {
		return nil, err
	}

	err = u.serviceRepo.CreateCoachServices(ctx, cmd)
	if err != nil {
		return nil, err
	}

	services, err := u.serviceRepo.GetCoachServices(ctx, cmd.CoachId)
	if err != nil {
		return nil, err
	}

	return services, nil
}
