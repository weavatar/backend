CREATE TABLE avatars
(
    hash       CHAR(32) PRIMARY KEY NOT NULL,
    raw        VARCHAR(255) UNIQUE,
    user_id    BIGINT  DEFAULT NULL,
    ban        BOOLEAN DEFAULT '0',
    checked    BOOLEAN DEFAULT '0',
    created_at TIMESTAMP(3)         NOT NULL,
    updated_at TIMESTAMP(3)         NOT NULL
);

COMMENT ON COLUMN avatars.hash IS '哈希';
COMMENT ON COLUMN avatars.raw IS '原始';
COMMENT ON COLUMN avatars.user_id IS '用户ID';
COMMENT ON COLUMN avatars.ban IS '禁用';
COMMENT ON COLUMN avatars.checked IS '已检查';
COMMENT ON COLUMN avatars.created_at IS '创建时间';
COMMENT ON COLUMN avatars.updated_at IS '更新时间';

CREATE INDEX idx_avatars_user_id ON avatars (user_id);
