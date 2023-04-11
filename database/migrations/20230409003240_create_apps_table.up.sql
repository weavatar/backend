CREATE TABLE apps
(
    id         BIGINT PRIMARY KEY NOT NULL,
    user_id    BIGINT             NOT NULL,
    name       VARCHAR(255)       NOT NULL,
    Secret     VARCHAR(255)       NOT NULL,
    created_at TIMESTAMP(3)       NOT NULL,
    updated_at TIMESTAMP(3)       NOT NULL
);

COMMENT ON COLUMN apps.id IS 'ID';
COMMENT ON COLUMN apps.user_id IS '用户ID';
COMMENT ON COLUMN apps.name IS '应用名称';
COMMENT ON COLUMN apps.Secret IS '应用密钥';
COMMENT ON COLUMN apps.created_at IS '创建时间';
COMMENT ON COLUMN apps.updated_at IS '更新时间';

CREATE INDEX idx_apps_user_id ON apps (user_id);
