CREATE TABLE app_info (
    info_key   VARCHAR(64)  NOT NULL PRIMARY KEY,
    info_value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO app_info (info_key, info_value) VALUES ('app_name', 'Alumni Portal');
