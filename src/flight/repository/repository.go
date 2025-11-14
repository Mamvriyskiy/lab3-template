package repository

import (
	model "github.com/Mamvriyskiy/lab3-template/src/flight/model"
	"github.com/jmoiron/sqlx"
)

type RepoFlight interface {
	GetFlights(page, size int) (model.FlightResponse, error)
	GetInfoAboutFlightByFlightNumber(flightNumber string) (model.Flight, error)
}

type Repository struct {
	RepoFlight
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		RepoFlight: NewFlightPostgres(db),
	}
}
