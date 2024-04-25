CREATE TABLE users
(
    id         TEXT PRIMARY KEY           NOT NULL,
    open_id    TEXT UNIQUE                NOT NULL,
    union_id   TEXT UNIQUE                NOT NULL,
    nickname   TEXT                       NOT NULL,
    avatar     TEXT                       NOT NULL,
    is_admin   BOOLEAN      DEFAULT FALSE NOT NULL,
    real_name  BOOLEAN      DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP(0)               NOT NULL,
    updated_at TIMESTAMP(0)               NOT NULL,
    deleted_at TIMESTAMP(0) DEFAULT NULL
);

COMMENT ON TABLE users IS '用户';
COMMENT ON COLUMN users.id IS 'ID';
COMMENT ON COLUMN users.open_id IS 'OpenID';
COMMENT ON COLUMN users.union_id IS 'UnionID';
COMMENT ON COLUMN users.nickname IS '昵称';
COMMENT ON COLUMN users.avatar IS '头像';
COMMENT ON COLUMN users.is_admin IS '是否是管理员';
COMMENT ON COLUMN users.real_name IS '是否实名认证';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';
COMMENT ON COLUMN users.deleted_at IS '删除时间';

CREATE INDEX idx_users_nickname ON users (nickname);
