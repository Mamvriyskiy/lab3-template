package repository

import (
	"fmt"
	"time"

	"strings"

	model "github.com/Mamvriyskiy/lab3-template/src/bonus/model"
	"github.com/jmoiron/sqlx"
)

type BonusPostgres struct {
	db *sqlx.DB
}

func NewBonusPostgres(db *sqlx.DB) *BonusPostgres {
	return &BonusPostgres{db: db}
}

func (r *BonusPostgres) UpdateBonusBonus(username, uid string, price int) (int, error) {
	bonusAmount := price / 10

	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	ticketUID := strings.Trim(uid, `"`)

	var privilegeID int
	err = tx.QueryRow(`
        SELECT id FROM privilege WHERE username = $1
    `, username).Scan(&privilegeID)

	if err != nil {
		return 0, fmt.Errorf("failed to get privilege ID: %w", err)
	}

	var updatedBalance int64

	err = tx.QueryRow(`
			UPDATE privilege
			SET balance = balance + $1
			WHERE username = $2
			RETURNING balance
		`, bonusAmount, username).Scan(&updatedBalance)

	if err != nil {
		return 0, fmt.Errorf("failed to update bonus balance: %w", err)
	}

	_, err = tx.Exec(`
        INSERT INTO privilege_history 
        (privilege_id, ticket_uid, datetime, balance_diff, operation_type) 
        VALUES ($1, $2, $3, $4, $5)
    `, privilegeID, ticketUID, time.Now(), bonusAmount, "FILL_IN_BALANCE")

	if err != nil {
		return 0, fmt.Errorf("failed to insert into privilege history: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return int(updatedBalance), nil
}

func (r *BonusPostgres) GetInfoAboutUserPrivilege(username string) (model.PrivilegeResponse, error) {
	var resp model.PrivilegeResponse
	var privilegeID int

	err := r.db.QueryRow(`
        SELECT id, balance, status
        FROM privilege
        WHERE username = $1
    `, username).Scan(&privilegeID, &resp.Balance, &resp.Status)
	if err != nil {
		return resp, err
	}

	rows, err := r.db.Query(`
        SELECT datetime, ticket_uid, balance_diff, operation_type
        FROM privilege_history
        WHERE privilege_id = $1
        ORDER BY datetime DESC
    `, privilegeID)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.HistoryItem
		var dt time.Time
		if err := rows.Scan(&dt, &item.TicketUid, &item.BalanceDiff, &item.OperationType); err != nil {
			return resp, err
		}
		item.Date = dt.Format(time.RFC3339)

		resp.History = append(resp.History, item)
	}

	return resp, nil
}

func (r *BonusPostgres) UpdateBonus(username, uid string, price int) (model.PrivilegeInfo, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return model.PrivilegeInfo{}, err
	}
	defer tx.Rollback()
	ticketUID := strings.Trim(uid, `"`)

	// 1. Получаем ID и текущий баланс пользователя
	var privilegeID, balance int
	query := `SELECT id, balance FROM privilege WHERE username = $1`
	err = tx.QueryRow(query, username).Scan(&privilegeID, &balance)
	if err != nil {
		return model.PrivilegeInfo{}, fmt.Errorf("user not found: %w", err)
	}

	var balance_diff int
	if price%balance == 0 {
		balance_diff = balance - price
	} else {
		balance_diff = price
	}

	insertHistory := `
		INSERT INTO privilege_history (
			privilege_id, ticket_uid, datetime, balance_diff, operation_type
		)
		VALUES ($1, $2, NOW(), $3, 'DEBIT_THE_ACCOUNT')`

	_, err = tx.Exec(insertHistory, privilegeID, ticketUID, balance_diff)
	if err != nil {
		return model.PrivilegeInfo{}, fmt.Errorf("failed to insert privilege history: %w", err)
	}

	var info model.PrivilegeInfo

	queryInfo := `
		SELECT p.status, p.balance, ph.balance_diff
		FROM privilege p
		JOIN privilege_history ph ON p.id = ph.privilege_id
		WHERE p.username = $1
		ORDER BY ph.datetime DESC
		LIMIT 1;
	`

	err = tx.QueryRow(queryInfo, username).Scan(&info.Status, &info.Balance, &info.BalanceDiff)
	if err != nil {
		return model.PrivilegeInfo{}, fmt.Errorf("failed to select privilege info: %w", err)
	}
	tx.Commit()

	return info, err
}

func (r *BonusPostgres) UpdateBonusDelete(username string, price int) error {
    _, err := r.db.Exec(`
        UPDATE privilege
        SET balance = balance - $1
        WHERE username = $2
    `, price, username)
	
	print("kkkkkk", err)
    if err != nil {
        return fmt.Errorf("failed to subtract bonus: %w", err)
    }

    return nil
}

