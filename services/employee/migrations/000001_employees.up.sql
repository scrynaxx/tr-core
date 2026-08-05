BEGIN;

CREATE TABLE employee.employees (
    id uuid,
    type text NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    patronymic text NOT NULL,
    phone text NOT NULL,
    birth_date date NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    archived_at timestamptz,
    CONSTRAINT pk_employees PRIMARY KEY (id),
    CONSTRAINT check_employees_type CHECK (type IN ('owner', 'manager', 'foreman', 'loader', 'assembler'))
);

CREATE TABLE employee.credentials (
    employee_id uuid,
    email text NOT NULL,
    password_hash text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT pk_credentials PRIMARY KEY (employee_id)
);

ALTER TABLE employee.credentials ADD CONSTRAINT fk_credentials_employee FOREIGN KEY (employee_id) REFERENCES employee.employees (id) ON DELETE CASCADE;
CREATE UNIQUE INDEX uidx_credentials_email ON employee.credentials (email);

CREATE TABLE employee.passports (
    employee_id uuid,
    series text NOT NULL,
    number text NOT NULL,
    issued_by text NOT NULL,
    issued_at date NOT NULL,
    department_code text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT pk_passports PRIMARY KEY (employee_id)
);

ALTER TABLE employee.passports ADD CONSTRAINT fk_passports_employee FOREIGN KEY (employee_id) REFERENCES employee.employees (id) ON DELETE CASCADE;

WITH emp AS (
    INSERT INTO employee.employees (id, type, first_name, last_name, patronymic, phone, birth_date)
    VALUES (gen_random_uuid(), 'owner', 'Михай', '', '', '977 667 50 15', '2026-01-01'::date)
    RETURNING id
)
INSERT INTO employee.credentials (employee_id, email, password_hash)
VALUES ((SELECT id FROM emp), 'demo@demo.ru', '$2a$10$1h6.JneqyaOoLCOnKhecle9U/FwL2n2xQPG59Y5aBGcwOVrQPDDmO');

COMMIT;