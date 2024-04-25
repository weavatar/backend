CREATE TABLE apps
(
    id         TEXT PRIMARY KEY NOT NULL,
    user_id    TEXT             NOT NULL,
    name       TEXT             NOT NULL,
    secret     TEXT             NOT NULL,
    created_at TIMESTAMP(0)     NOT NULL,
    updated_at TIMESTAMP(0)     NOT NULL
);

COMMENT ON TABLE apps IS '应用';
COMMENT ON COLUMN apps.id IS 'ID';
COMMENT ON COLUMN apps.user_id IS '用户ID';
COMMENT ON COLUMN apps.name IS '应用名称';
COMMENT ON COLUMN apps.secret IS '应用密钥';
COMMENT ON COLUMN apps.created_at IS '创建时间';
COMMENT ON COLUMN apps.updated_at IS '更新时间';

CREATE INDEX idx_apps_user_id ON apps (user_id);
