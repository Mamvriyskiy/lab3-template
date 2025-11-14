package services

import (
	model "github.com/Mamvriyskiy/lab3-template/src/bonus/model"
	repository "github.com/Mamvriyskiy/lab3-template/src/bonus/repository"
)

type BonusService struct {
	repo repository.RepoBonus
}

func NewBonusService(repo repository.RepoBonus) *BonusService {
	return &BonusService{repo: repo}
}

func (s *BonusService) UpdateBonusBonus(username, ticketUid string, price int) error {
	return s.repo.UpdateBonusBonus(username, ticketUid, price)
}

func (s *BonusService) GetInfoAboutUserPrivilege(username string) (model.PrivilegeResponse, error) {
	return s.repo.GetInfoAboutUserPrivilege(username)
}

func (s *BonusService) UpdateBonus(username, ticketUID string, price int) (model.PrivilegeInfo, error) {
	return s.repo.UpdateBonus(username, ticketUID, price)
}
