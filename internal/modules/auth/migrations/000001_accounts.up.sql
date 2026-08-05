CREATE TABLE auth.accounts
(
    id            uuid,
    actor         text        NOT NULL,
    email         text        NOT NULL,
    password_hash text        NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT pk_accounts PRIMARY KEY (id)
);

CREATE UNIQUE INDEX uidx_account_actor_email ON auth.accounts (actor, email);

INSERT INTO auth.accounts (id, actor, email, password_hash)
VALUES ('2dc9fb5d-28a2-4060-b615-b64d58df8377',
        'employee',
        'demo@demo.ru',
        '$2a$10$1h6.JneqyaOoLCOnKhecle9U/FwL2n2xQPG59Y5aBGcwOVrQPDDmO');