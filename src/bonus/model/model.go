package model

type HistoryItem struct {
	Date          string `json:"date"`
	TicketUid     string `json:"ticketUid"`
	BalanceDiff   int    `json:"balanceDiff"`
	OperationType string `json:"operationType"`
}

type PrivilegeResponse struct {
	Balance int           `json:"balance"`
	Status  string        `json:"status"`
	History []HistoryItem `json:"history"`
}

type PrivilegeInfo struct {
    Status      string `db:"status"`
    Balance     int    `db:"balance"`
    BalanceDiff int    `db:"balance_diff"`
}
