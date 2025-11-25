package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionOperationType string

const (
	TransactionOperationTypePurchase                 = TransactionOperationType("PURCHASE")
	TransactionOperationTypeWithdrawal               = TransactionOperationType("WITHDRAWAL")
	TransactionOperationTypeCreditVoucher            = TransactionOperationType("CREDIT_VOUCHER")
	TransactionOperationTypePurchaseWithInstallments = TransactionOperationType("PURCHASE_WITH_INSTALLMENTS")
)

type Transaction struct {
	Base
	OperationType TransactionOperationType
	Amount        int64
	ExternalID    uuid.UUID
	EventDate     time.Time
	AccountID     int64
	Account       Account `gorm:"foreignKey:AccountID"`
}
