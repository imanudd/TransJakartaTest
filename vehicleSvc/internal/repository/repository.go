package repository

import (
	"database/sql"
)

type RepositoryImpl interface {
	GetVehicleRepo() VehicleRepositoryImpl
}

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) RepositoryImpl {
	return &Repository{
		db: database,
	}
}

func (r *Repository) GetVehicleRepo() VehicleRepositoryImpl { return NewVehicleRepository(r.db) }
