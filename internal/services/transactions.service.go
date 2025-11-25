package services

import (
	"context"
	"transactions/internal/domain"
	"transactions/internal/repositories"
)

type Transaction struct {
	transactionRepository repositories.Transaction
}

func (t Transaction) Save(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	return t.transactionRepository.Save(ctx, transaction)
}

func NewTransaction(transactionRepository repositories.Transaction) Transaction {
	return Transaction{
		transactionRepository: transactionRepository,
	}
}
