BEGIN;

CREATE TABLE auth.sessions (
    id uuid,
    employee_id uuid NOT NULL,
    refresh_hash text NOT NULL,
    user_agent text NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT pk_sessions PRIMARY KEY (id)
);

CREATE UNIQUE INDEX uidx_session_refresh_hash ON auth.sessions (refresh_hash);
CREATE INDEX idx_sessions_employee_id ON auth.sessions (employee_id);
CREATE INDEX idx_sessions_expires_at  ON auth.sessions (expires_at);

COMMIT;
