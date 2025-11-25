package responses

import (
	"time"
	"transactions/internal/domain"

	"github.com/google/uuid"
)

type Account struct {
	UID            uuid.UUID `json:"uid"`
	DocumentNumber string    `json:"document_number"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Account) FromDomain(account domain.Account) Account {
	return Account{
		UID:            account.UUID,
		DocumentNumber: account.DocumentNumber,
		CreatedAt:      account.CreatedAt,
		UpdatedAt:      account.UpdatedAt,
	}
}
