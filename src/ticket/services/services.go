package services

import (
	model "github.com/Mamvriyskiy/lab3-template/src/ticket/model"
	repository "github.com/Mamvriyskiy/lab3-template/src/ticket/repository"
)

type Ticket interface {
	GetInfoAboutTiket(ticketUID string) (model.Ticket, error)
	GetInfoAboutTikets(username string) ([]model.Ticket, error)
	UpdateStatusTicket(ticket string) error
	CreateTicket(username, flightNumber string, price int) (string, error)
}

type Services struct {
	Ticket
}

func NewServices(repo *repository.Repository) *Services {
	return &Services{
		Ticket: NewTicketService(repo.RepoTicket),
	}
}
