package repository

import (
	model "github.com/Mamvriyskiy/lab3-template/src/ticket/model"
	"github.com/jmoiron/sqlx"
)

type RepoTicket interface {
	GetInfoAboutTiket(ticketUID string) (model.Ticket, error)
	GetInfoAboutTikets(username string) ([]model.Ticket, error)
	UpdateStatusTicket(ticket string) error
	CreateTicket(username, flightNumber string, price int) (string, error)
}

type Repository struct {
	RepoTicket
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		RepoTicket: NewTicketPostgres(db),
	}
}
