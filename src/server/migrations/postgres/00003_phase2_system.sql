-- +goose Up
CREATE TABLE sys_menus (
    id BIGINT PRIMARY KEY,
    parent_id BIGINT NOT NULL DEFAULT 0,
    type VARCHAR(16) NOT NULL,
    name VARCHAR(100) NOT NULL,
    path VARCHAR(255) NOT NULL DEFAULT '',
    component VARCHAR(255) NOT NULL DEFAULT '',
    icon VARCHAR(100) NOT NULL DEFAULT '',
    sort INTEGER NOT NULL DEFAULT 0,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    permission VARCHAR(150) NULL UNIQUE,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (parent_id, name)
);
CREATE INDEX idx_sys_menus_parent_id ON sys_menus (parent_id);

CREATE TABLE sys_role_menus (
    role_id BIGINT NOT NULL REFERENCES sys_roles(id) ON DELETE CASCADE,
    menu_id BIGINT NOT NULL REFERENCES sys_menus(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, menu_id)
);

CREATE TABLE sys_dictionaries (
    id BIGINT PRIMARY KEY,
    code VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sys_dictionary_items (
    id BIGINT PRIMARY KEY,
    dictionary_id BIGINT NOT NULL REFERENCES sys_dictionaries(id) ON DELETE CASCADE,
    label VARCHAR(100) NOT NULL,
    value VARCHAR(255) NOT NULL,
    sort INTEGER NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (dictionary_id, value)
);

CREATE TABLE sys_operation_logs (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    request_id VARCHAR(64) NOT NULL,
    method VARCHAR(10) NOT NULL,
    path VARCHAR(255) NOT NULL,
    status_code INTEGER NOT NULL,
    duration_ms BIGINT NOT NULL,
    ip VARCHAR(64) NOT NULL DEFAULT '',
    user_agent VARCHAR(500) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_sys_operation_logs_user_id ON sys_operation_logs (user_id);
CREATE INDEX idx_sys_operation_logs_created_at ON sys_operation_logs (created_at);

-- +goose Down
DROP TABLE sys_operation_logs;
DROP TABLE sys_dictionary_items;
DROP TABLE sys_dictionaries;
DROP TABLE sys_role_menus;
DROP TABLE sys_menus;
