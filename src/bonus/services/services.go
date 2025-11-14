package services

import (
	model "github.com/Mamvriyskiy/lab3-template/src/bonus/model"
	repository "github.com/Mamvriyskiy/lab3-template/src/bonus/repository"
)

type Bonus interface {
	GetInfoAboutUserPrivilege(username string) (model.PrivilegeResponse, error)
	UpdateBonus(username, ticketUID string, price int) (model.PrivilegeInfo, error)
	UpdateBonusBonus(username, ticketUid string, price int) error
}

type Services struct {
	Bonus
}

func NewServices(repo *repository.Repository) *Services {
	return &Services{
		Bonus: NewBonusService(repo.RepoBonus),
	}
}
