package requests

import "transactions/internal/domain"

type CreateAccount struct {
	DocumentNumber string `json:"document_number" validate:"required,min=11,max=14"`
}

func (c CreateAccount) Domain() domain.Account {
	return domain.Account{
		DocumentNumber: c.DocumentNumber,
	}
}
