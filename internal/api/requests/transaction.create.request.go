package requests

import (
	"slices"
	"time"
	"transactions/internal/domain"

	"github.com/google/uuid"
)

type CreateTransaction struct {
	ExternalID    uuid.UUID                       `json:"external_id" validate:"required,uuid"`
	AccountID     uuid.UUID                       `json:"account_id" validate:"required,uuid"`
	OperationType domain.TransactionOperationType `json:"operation_type" validate:"required,oneof=PURCHASE WITHDRAWAL CREDIT_VOUCHER PURCHASE_WITH_INSTALLMENTS"`
	Amount        int64                           `json:"amount" validate:"required,gt=0"`
}

func (c CreateTransaction) Domain(account domain.Account) domain.Transaction {
	multiplier := 1

	negativeAmountOperations := []domain.TransactionOperationType{
		domain.TransactionOperationTypePurchase,
		domain.TransactionOperationTypePurchaseWithInstallments,
		domain.TransactionOperationTypeWithdrawal,
	}

	if slices.Contains(negativeAmountOperations, c.OperationType) {
		multiplier = -1
	}

	return domain.Transaction{
		OperationType: c.OperationType,
		Amount:        c.Amount * int64(multiplier),
		ExternalID:    c.ExternalID,
		EventDate:     time.Now().UTC(),
		Account:       account,
		AccountID:     int64(account.ID),
	}
}
