CREATE TABLE auth.sessions
(
    id           uuid,
    account_id   uuid        NOT NULL,
    refresh_hash text        NOT NULL,
    user_agent   text        NOT NULL,
    expires_at   timestamptz NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT pk_sessions PRIMARY KEY (id),
    CONSTRAINT fk_sessions_account FOREIGN KEY (account_id) REFERENCES auth.accounts (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX uidx_session_refresh_hash ON auth.sessions (refresh_hash);
CREATE INDEX idx_sessions_expires_at ON auth.sessions (expires_at);
