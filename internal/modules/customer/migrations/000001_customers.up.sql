CREATE TABLE customer.customers
(
    id          uuid,
    account_id     uuid,
    first_name  text        NOT NULL,
    last_name   text        NOT NULL,
    patronymic  text,
    phone       text        NOT NULL,
    email       text        NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    archived_at timestamptz,
    CONSTRAINT pk_customers PRIMARY KEY (id)
);
CREATE UNIQUE INDEX uidx_customers_account_id ON customer.customers (account_id) WHERE account_id IS NOT NULL;
CREATE UNIQUE INDEX uidx_customers_phone ON customer.customers (phone);
CREATE UNIQUE INDEX uidx_customers_email ON customer.customers (email);