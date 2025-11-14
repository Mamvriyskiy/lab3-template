package repository

import (
	"time"

	model "github.com/Mamvriyskiy/lab3-template/src/flight/model"
	"github.com/jmoiron/sqlx"
)

type FlightPostgres struct {
	db *sqlx.DB
}

func NewFlightPostgres(db *sqlx.DB) *FlightPostgres {
	return &FlightPostgres{db: db}
}

func (r *FlightPostgres) GetInfoAboutFlightByFlightNumber(flightNumber string) (model.Flight, error) {
	var flight model.Flight

	query := `
		SELECT 
			f.flight_number,
			af.city || ' ' || af.name AS from_airport,
			at.city || ' ' || at.name AS to_airport,
			f.datetime,
			f.price
		FROM flight f
		JOIN airport af ON f.from_airport_id = af.id
		JOIN airport at ON f.to_airport_id = at.id
		WHERE f.flight_number = $1
	`

	var dt time.Time
	row := r.db.QueryRow(query, flightNumber)
	err := row.Scan(
		&flight.FlightNumber,
		&flight.FromAirport,
		&flight.ToAirport,
		&dt,
		&flight.Price,
	)

	loc, _ := time.LoadLocation("Europe/Moscow")
	flight.Datetime = dt.In(loc).Format("2006-01-02 15:04")

	if err != nil {
		return model.Flight{}, err
	}

	return flight, nil
}

func (r *FlightPostgres) GetFlights(page, size int) (model.FlightResponse, error) {
	offset := (page - 1) * size

	rows, err := r.db.Query(`
        SELECT
            f.flight_number,
            a_from.name AS from_airport,
            a_to.name AS to_airport,
			a_from.city as from_city,
			a_to.city AS to_city,
            f.datetime,
            f.price
        FROM flight f
        JOIN airport a_from ON f.from_airport_id = a_from.id
        JOIN airport a_to ON f.to_airport_id = a_to.id
        ORDER BY f.id
        LIMIT $1 OFFSET $2`, size, offset)
	if err != nil {
		return model.FlightResponse{}, err
	}
	defer rows.Close()

	var items []model.FlightItem
	for rows.Next() {
		var item model.FlightItem
		var dt time.Time
		var from_airport, to_airport, from_city, to_city string
		if err := rows.Scan(&item.FlightNumber, &from_airport, &to_airport, &from_city, &to_city, &dt, &item.Price); err != nil {
			return model.FlightResponse{}, err
		}
		item.FromAirport = from_city + " " + from_airport
		item.ToAirport = to_city + " " + to_airport

		loc, _ := time.LoadLocation("Europe/Moscow")
		item.Date = dt.In(loc).Format("2006-01-02 15:04")

		items = append(items, item)
	}

	var total int
	err = r.db.QueryRow("SELECT COUNT(*) FROM flight").Scan(&total)
	if err != nil {
		return model.FlightResponse{}, err
	}

	return model.FlightResponse{
		Page:          page,
		PageSize:      len(items),
		TotalElements: total,
		Items:         items,
	}, nil
}
