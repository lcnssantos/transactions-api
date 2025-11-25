package responses

import (
	"time"
	"transactions/internal/domain"

	"github.com/google/uuid"
)

type Transaction struct {
	UID           uuid.UUID                       `json:"uid"`
	OperationType domain.TransactionOperationType `json:"operation_type"`
	Amount        int64                           `json:"amount"`
	ExternalID    uuid.UUID                       `json:"external_id"`
	EventDate     time.Time                       `json:"event_date"`
	Account       Account                         `json:"account"`
}

func (Transaction) FromDomain(transaction domain.Transaction) Transaction {
	return Transaction{
		UID:           transaction.UUID,
		OperationType: transaction.OperationType,
		Amount:        transaction.Amount,
		ExternalID:    transaction.ExternalID,
		EventDate:     transaction.EventDate,
		Account:       Account{}.FromDomain(transaction.Account),
	}
}
