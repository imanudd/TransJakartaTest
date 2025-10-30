package repository

import (
	"app/internal/domain"
	"context"
	"database/sql"
	"errors"
)

type VehicleRepositoryImpl interface {
	GetHistory(ctx context.Context, vehicleId string, req *domain.GetVehicleHistoryRequest) ([]*domain.VehicleLocation, error)
	GetLatestLocation(ctx context.Context, vehicleId string) (*domain.VehicleLocation, error)
	InsertVehicleLocation(ctx context.Context, vehicle *domain.VehicleLocation) error
	Get(ctx context.Context) ([]*domain.Vehicle, error)
}

type VehicleRepository struct {
	db *sql.DB
}

func NewVehicleRepository(db *sql.DB) VehicleRepositoryImpl {
	return &VehicleRepository{
		db: db,
	}
}

func (r *VehicleRepository) GetHistory(ctx context.Context, vehicleId string, req *domain.GetVehicleHistoryRequest) ([]*domain.VehicleLocation, error) {
	var vehicles []*domain.VehicleLocation

	rows, err := r.db.QueryContext(ctx, ` SELECT vehicle_id, latitude, longitude, timestamp FROM vehicle_locations WHERE vehicle_id = $1 AND timestamp BETWEEN $2 AND $3`, vehicleId, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data domain.VehicleLocation
		err = rows.Scan(
			&data.VehicleId,
			&data.Latitude,
			&data.Longitude,
			&data.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		vehicles = append(vehicles, &data)
	}

	return vehicles, nil
}

func (r *VehicleRepository) GetLatestLocation(ctx context.Context, vehicleId string) (*domain.VehicleLocation, error) {
	var result domain.VehicleLocation

	selectQuery := `SELECT vehicle_id, latitude, longitude, timestamp FROM vehicle_locations WHERE vehicle_id = $1 order by timestamp desc limit 1;`

	row := r.db.QueryRowContext(ctx, selectQuery, vehicleId)
	err := row.Scan(&result.VehicleId, &result.Latitude, &result.Longitude, &result.Timestamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &result, err
}

func (r *VehicleRepository) InsertVehicleLocation(ctx context.Context, data *domain.VehicleLocation) error {
	insertQuery := ` INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp) VALUES ($1,$2,$3,$4); `

	_, err := r.db.ExecContext(ctx, insertQuery, data.VehicleId, data.Latitude, data.Longitude, data.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (r *VehicleRepository) Get(ctx context.Context) ([]*domain.Vehicle, error) {
	var vehicles []*domain.Vehicle

	selectQuery := `SELECT vehicle_id, vehicle_type, vehicle_code FROM vehicles`

	rows, err := r.db.QueryContext(ctx, selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vehicle domain.Vehicle
		err = rows.Scan(
			&vehicle.VehicleId,
			&vehicle.VehicleType,
			&vehicle.VehicleCode,
		)
		if err != nil {
			return nil, err
		}

		vehicles = append(vehicles, &vehicle)
	}

	return vehicles, nil
}
