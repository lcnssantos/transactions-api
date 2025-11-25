package requests

import "github.com/google/uuid"

type FindAccount struct {
	AccountID uuid.UUID `param:"account_id" validate:"required,uuid"`
}
