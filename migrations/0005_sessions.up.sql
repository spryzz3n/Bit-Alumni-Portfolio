CREATE TABLE sessions (
    token_hash BINARY(32)      NOT NULL PRIMARY KEY,
    user_id    BIGINT UNSIGNED NOT NULL,
    expires_at TIMESTAMP       NOT NULL,
    created_at TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_sessions_user (user_id),
    KEY idx_sessions_expires (expires_at),
    CONSTRAINT fk_session_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
