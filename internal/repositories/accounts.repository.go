package repositories

import (
	"context"
	"errors"
	"transactions/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrAccountDocumentAlreadyExist = errors.New("account_document_already_exists")
var ErrAccountNotFound = errors.New("account_not_found")

type Account struct {
	db *gorm.DB
}

func (a Account) FindByUUID(ctx context.Context, uid uuid.UUID) (domain.Account, error) {
	account := domain.Account{}

	err := a.db.WithContext(ctx).Where("uuid = ?", uid).First(&account).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Account{}, ErrAccountNotFound
	}

	return account, err
}

func (a Account) Create(ctx context.Context, account domain.Account) (domain.Account, error) {
	err := a.db.WithContext(ctx).Save(&account).Error

	pgErr, ok := err.(*pgconn.PgError)

	if ok {
		if pgErr.Code == "23505" && pgErr.ConstraintName == "accounts_document_number_key" {
			return domain.Account{}, ErrAccountDocumentAlreadyExist
		}
	}

	return account, err
}

func NewAccount(db *gorm.DB) Account {
	return Account{
		db: db,
	}
}
