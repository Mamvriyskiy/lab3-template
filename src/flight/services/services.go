package services

import (
	"github.com/Mamvriyskiy/lab3-template/src/flight/model"
	"github.com/Mamvriyskiy/lab3-template/src/flight/repository"
)

type Flight interface {
	GetInfoAboutFlight(page, size int) (model.FlightResponse, error)
	GetInfoAboutFlightByFlightNumber(flightNumber string) (model.Flight, error)
}

type Services struct {
	Flight
}

func NewServices(repo *repository.Repository) *Services {
	return &Services{
		Flight: NewFlightService(repo.RepoFlight),
	}
}
