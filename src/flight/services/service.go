package services

import (
	model "github.com/Mamvriyskiy/lab3-template/src/flight/model"
	"github.com/Mamvriyskiy/lab3-template/src/flight/repository"
)

type FlightService struct {
	repo repository.RepoFlight
}

func NewFlightService(repo repository.RepoFlight) *FlightService {
	return &FlightService{repo: repo}
}

func (s *FlightService) GetInfoAboutFlight(page, size int) (model.FlightResponse, error) {
	return s.repo.GetFlights(page, size)
}

func (s *FlightService) GetInfoAboutFlightByFlightNumber(flightNumber string) (model.Flight, error) {
	return s.repo.GetInfoAboutFlightByFlightNumber(flightNumber)
}
