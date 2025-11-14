package repository

import (
	model "github.com/Mamvriyskiy/lab3-template/src/bonus/model"
	"github.com/jmoiron/sqlx"
)

type RepoBonus interface {
	GetInfoAboutUserPrivilege(username string) (model.PrivilegeResponse, error)
	UpdateBonus(username, ticketUID string, price int) (model.PrivilegeInfo, error)
	UpdateBonusBonus(username, ticketUid string, price int) error
}

type Repository struct {
	RepoBonus
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		RepoBonus: NewBonusPostgres(db),
	}
}
