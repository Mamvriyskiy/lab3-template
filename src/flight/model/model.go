package model

type FlightItem struct {
	FlightNumber string `json:"flightNumber"`
	FromAirport  string `json:"fromAirport"`
	ToAirport    string `json:"toAirport"`
	Date         string `json:"date"`
	Price        int    `json:"price"`
}

type FlightResponse struct {
	Page          int          `json:"page"`
	PageSize      int          `json:"pageSize"`
	TotalElements int          `json:"totalElements"`
	Items         []FlightItem `json:"items"`
}

type Flight struct {
	ID           int    `json:"id"`
	FlightNumber string `json:"flightNumber"`
	Datetime     string `json:"datetime"`
	FromAirport  string `json:"fromAirport"`
	ToAirport    string `json:"toAirport"`
	Price        int    `json:"price"`
}
