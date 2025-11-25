package services

import (
	"context"
	"transactions/internal/domain"
	"transactions/internal/repositories"

	"github.com/google/uuid"
)

type Account struct {
	accountRepository repositories.Account
}

func (a Account) FindByUUID(ctx context.Context, uid uuid.UUID) (domain.Account, error) {
	return a.accountRepository.FindByUUID(ctx, uid)
}

func (a Account) Create(ctx context.Context, account domain.Account) (domain.Account, error) {
	return a.accountRepository.Create(ctx, account)
}

func NewAccount(accountRepository repositories.Account) Account {
	return Account{
		accountRepository: accountRepository,
	}
}
