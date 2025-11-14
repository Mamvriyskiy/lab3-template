package services

import (
	model "github.com/Mamvriyskiy/lab3-template/src/ticket/model"
	repository "github.com/Mamvriyskiy/lab3-template/src/ticket/repository"
)

type TicketService struct {
	repo repository.RepoTicket
}

func NewTicketService(repo repository.RepoTicket) *TicketService {
	return &TicketService{repo: repo}
}

func (s *TicketService) GetInfoAboutTiket(ticketUID string) (model.Ticket, error) {
	return s.repo.GetInfoAboutTiket(ticketUID)
}

func (s *TicketService) GetInfoAboutTikets(username string) ([]model.Ticket, error) {
	return s.repo.GetInfoAboutTikets(username)
}

func (s *TicketService) UpdateStatusTicket(ticket string) error {
	return s.repo.UpdateStatusTicket(ticket)
}

func (s *TicketService) CreateTicket(username, flightNumber string, price int) (string, error) {
	return s.repo.CreateTicket(username, flightNumber, price)
}
