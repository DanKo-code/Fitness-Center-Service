package postgres

import (
	"Service/internal/dtos"
	customErrors "Service/internal/errors"
	"Service/internal/models"
	"Service/pkg/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ServiceRepository struct {
	db *sqlx.DB
}

func NewServiceRepository(db *sqlx.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (serviceRep *ServiceRepository) CreateService(ctx context.Context, service *models.Service) error {
	_, err := serviceRep.db.NamedExecContext(ctx, `
		INSERT INTO "service" (id, title, photo, created_time, updated_time)
		VALUES (:id, :title, :photo, :created_time, :updated_time)`, *service)
	if err != nil {
		logger.ErrorLogger.Printf("Error CreateService: %v", err)
		return err
	}

	return nil
}

func (serviceRep *ServiceRepository) GetServiceById(ctx context.Context, id uuid.UUID) (*models.Service, error) {
	service := &models.Service{}
	err := serviceRep.db.GetContext(ctx, service, `SELECT id, title, photo, created_time, updated_time FROM "service" WHERE id = $1`, id)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetServiceById: %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.ServiceNotFound
		}

		return nil, err
	}

	return service, nil
}

func (serviceRep *ServiceRepository) UpdateService(ctx context.Context, cmd *dtos.UpdateServiceCommand) error {

	setFields := map[string]interface{}{}

	if cmd.Title != "" {
		setFields["name"] = cmd.Title
	}
	if cmd.Photo != "" {
		setFields["photo"] = cmd.Photo
	}
	setFields["updated_time"] = cmd.UpdatedTime

	if len(setFields) == 0 {
		logger.InfoLogger.Printf("No fields to update for service Id: %v", cmd.Id)
		return nil
	}

	query := `UPDATE "service" SET `

	var params []interface{}
	i := 1
	for field, value := range setFields {
		if i > 1 {
			query += ", "
		}

		query += fmt.Sprintf(`%s = $%d`, field, i)
		params = append(params, value)
		i++
	}
	query += fmt.Sprintf(` WHERE id = $%d`, i)
	params = append(params, cmd.Id)

	_, err := serviceRep.db.ExecContext(ctx, query, params...)
	if err != nil {
		logger.ErrorLogger.Printf("Error UpdateService: %v", err)
		return err
	}

	return nil
}

func (serviceRep *ServiceRepository) DeleteService(ctx context.Context, id uuid.UUID) error {
	_, err := serviceRep.db.ExecContext(ctx, `DELETE FROM "service" WHERE id = $1`, id)

	if err != nil {
		logger.ErrorLogger.Printf("Error DeleteService: %v", err)
		return err
	}

	return nil
}

func (serviceRep *ServiceRepository) GetServices(ctx context.Context) ([]*models.Service, error) {
	var services []*models.Service

	err := serviceRep.db.SelectContext(ctx, &services, `SELECT id, title, photo, created_time, updated_time FROM "service"`)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetServices: %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.ServiceNotFound
		}

		return nil, err
	}

	return services, nil
}

func (serviceRep *ServiceRepository) CreateCoachServices(ctx context.Context, cmd *dtos.CreateCoachServicesCommand) error {

	query := `
	INSERT INTO "coach_service" (coach_id, service_id)
	VALUES (:coach_id, :service_id)
`
	var values []map[string]interface{}
	for _, serviceId := range cmd.ServicesIds {
		values = append(values, map[string]interface{}{
			"coach_id":   cmd.CoachId,
			"service_id": serviceId,
		})
	}

	_, err := serviceRep.db.NamedExecContext(ctx, query, values)
	if err != nil {
		logger.ErrorLogger.Printf("Error CreateCoachServices: %v", err)
		return err
	}

	return nil
}

func (serviceRep *ServiceRepository) GetServicesByIds(ctx context.Context, ids []uuid.UUID) ([]*models.Service, error) {
	query := `SELECT id, title, photo, created_time, updated_time 
			  FROM "service"
			  WHERE id IN (?)`

	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to bind ids: %w", err)
	}

	query = serviceRep.db.Rebind(query)

	var services []*models.Service

	err = serviceRep.db.SelectContext(ctx, &services, query, args...)
	if err != nil {
		return nil, customErrors.ServiceNotFound
	}

	return services, nil
}

func (serviceRep *ServiceRepository) GetCoachServices(ctx context.Context, id uuid.UUID) ([]*models.Service, error) {
	var services []*models.Service
	err := serviceRep.db.SelectContext(ctx, &services, `SELECT id, title, photo, created_time, updated_time FROM "service" WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	return services, nil
}
