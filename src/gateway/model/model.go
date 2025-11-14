package model

type TicketInfo struct {
    TicketUID    string `json:"ticketUid"`
    FlightNumber string `json:"flightNumber"`
    FromAirport  string `json:"fromAirport"`
    ToAirport    string `json:"toAirport"`
    Date         string `json:"date"`
    Price        int    `json:"price"`
    Status       string `json:"status"`
}

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

type Ticket struct {
    TicketUID    string `json:"ticketUid"`
    Username     string `json:"username"`
    FlightNumber string `json:"flightNumber"`
    Price        int    `json:"price"`
    Status       string `json:"status"`
}

type PrivilegeResponse struct {
	Balance int           `json:"balance"`
	Status  string        `json:"status"`
}
