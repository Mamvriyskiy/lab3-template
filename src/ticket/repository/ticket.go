package repository

import (
	"fmt"

	model "github.com/Mamvriyskiy/lab3-template/src/ticket/model"
	"github.com/jmoiron/sqlx"
)

type TicketPostgres struct {
	db *sqlx.DB
}

func NewTicketPostgres(db *sqlx.DB) *TicketPostgres {
	return &TicketPostgres{db: db}
}

func (r *TicketPostgres) UpdateStatusTicket(ticketUid string) error {
	query := `
        UPDATE ticket
        SET status = 'CANCELED'
        WHERE ticket_uid = $1
          AND status = 'PAID'
    `

	res, err := r.db.Exec(query, ticketUid)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no ticket updated: ticket either does not exist or is not PAID")
	}

	return nil
}

func (r *TicketPostgres) GetInfoAboutTiket(ticketUID string) (model.Ticket, error) {
	query := `
        SELECT ticket_uid, username, flight_number, price, status
        FROM ticket
        WHERE ticket_uid = $1;
    `
	var ticket model.Ticket
	err := r.db.QueryRow(query, ticketUID).Scan(
		&ticket.TicketUID,
		&ticket.Username,
		&ticket.FlightNumber,
		&ticket.Price,
		&ticket.Status,
	)

	if err != nil {
		return model.Ticket{}, err
	}

	return ticket, nil
}

func (r *TicketPostgres) GetInfoAboutTikets(username string) ([]model.Ticket, error) {
	query := `
        SELECT ticket_uid, username, flight_number, price, status
        FROM ticket
        WHERE username = $1;
    `
	rows, err := r.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.Ticket
	for rows.Next() {
		var t model.Ticket
		if err := rows.Scan(&t.TicketUID, &t.Username, &t.FlightNumber, &t.Price, &t.Status); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, nil
}

func (r *TicketPostgres) CreateTicket(username, flightNumber string, price int) (string, error) {
	query := `
        INSERT INTO ticket (ticket_uid, username, flight_number, price, status) 
        VALUES (gen_random_uuid(), $1, $2, $3, 'PAID') 
        RETURNING ticket_uid
    `

	var ticketUID string
	err := r.db.QueryRow(query, username, flightNumber, price).Scan(&ticketUID)
	if err != nil {
		return "", fmt.Errorf("failed to create ticket: %w", err)
	}

	return ticketUID, nil
}
