BEGIN;

CREATE TABLE customer.customers (
    id uuid,
    name text NOT NULL,
    phone text NOT NULL,
    email text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    archived_at timestamptz,
    CONSTRAINT pk_customers PRIMARY KEY (id)
);

CREATE UNIQUE INDEX uidx_customers_phone_unique ON customer.customers (phone);
CREATE UNIQUE INDEX uidx_customers_email_unique ON customer.customers (email);
CREATE INDEX idx_customers_active ON customer.customers (name) WHERE archived_at IS NULL;

COMMIT;
