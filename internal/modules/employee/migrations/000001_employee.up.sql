CREATE TABLE employee.employees
(
    id          uuid,
    account_id     uuid,
    type        text        NOT NULL,
    first_name  text        NOT NULL,
    last_name   text        NOT NULL,
    patronymic  text        NOT NULL,
    phone       text        NOT NULL,
    birth_date  date        NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    archived_at timestamptz,
    CONSTRAINT pk_employees PRIMARY KEY (id),
    CONSTRAINT check_employees_type CHECK (type IN ('owner', 'manager', 'foreman', 'loader', 'assembler'))
);
CREATE UNIQUE INDEX uidx_employees_account_id ON employee.employees (account_id) WHERE account_id IS NOT NULL;
CREATE UNIQUE INDEX uidx_employees_phone ON employee.employees (phone);

CREATE TABLE employee.passports
(
    employee_id     uuid,
    series          text        NOT NULL,
    number          text        NOT NULL,
    issued_by       text        NOT NULL,
    issued_at       date        NOT NULL,
    department_code text        NOT NULL,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT pk_passports PRIMARY KEY (employee_id),
    CONSTRAINT fk_passports_employee FOREIGN KEY (employee_id) REFERENCES employee.employees (id) ON DELETE CASCADE
);

INSERT INTO employee.employees (id, account_id, type, first_name, last_name, patronymic, phone, birth_date)
VALUES (gen_random_uuid(), '2dc9fb5d-28a2-4060-b615-b64d58df8377', 'owner', 'Михай', 'Бусуёк', '', '977 667 50 15', '2026-01-01'::date);
