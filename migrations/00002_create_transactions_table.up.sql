CREATE TYPE transaction_operation_type AS ENUM (
    'PURCHASE',
    'WITHDRAWAL',
    'CREDIT_VOUCHER',
    'PURCHASE_WITH_INSTALLMENTS'
);

CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    external_id UUID NOT NULL UNIQUE,
    operation_type transaction_operation_type NOT NULL,
    amount BIGINT NOT NULL,
    event_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);