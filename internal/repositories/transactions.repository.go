package repositories

import (
	"context"
	"transactions/internal/domain"

	"gorm.io/gorm"
)

type Transaction struct {
	db *gorm.DB
}

func (t Transaction) Save(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	err := t.db.WithContext(ctx).Omit("Accounts").FirstOrCreate(&transaction, domain.Transaction{ExternalID: transaction.ExternalID}).Error

	return transaction, err
}

func NewTransaction(db *gorm.DB) Transaction {
	return Transaction{
		db: db,
	}
}
