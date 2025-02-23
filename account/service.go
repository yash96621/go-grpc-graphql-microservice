package account

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &accountService{repo}
}

func (s *accountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	a := &Account{Name: name, ID: uuid.New().String()}
	if err := s.repo.PutAccount(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if take > 100 || (take == 0 && skip == 0) {
		take = 100
	}
	return s.repo.ListAccounts(ctx, skip, take)
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	return s.repo.GetAccountByID(ctx, id)
}
