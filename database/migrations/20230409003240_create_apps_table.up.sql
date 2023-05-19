CREATE TABLE apps
(
    id         BIGINT PRIMARY KEY NOT NULL COMMENT 'ID',
    user_id    BIGINT NOT NULL COMMENT '用户ID',
    name       VARCHAR(255) NOT NULL COMMENT '应用名称',
    Secret     VARCHAR(255) NOT NULL COMMENT '应用密钥',
    created_at TIMESTAMP NOT NULL COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL COMMENT '更新时间'
);

CREATE INDEX idx_apps_user_id ON apps (user_id);
