CREATE TABLE avatars
(
    hash       CHAR(32) PRIMARY KEY  NOT NULL,
    raw        VARCHAR(255) UNIQUE,
    user_id    BIGINT  DEFAULT NULL,
    created_at TIMESTAMP             NOT NULL,
    updated_at TIMESTAMP             NOT NULL
);

COMMENT ON TABLE avatars IS '头像';
COMMENT ON COLUMN avatars.hash IS '哈希';
COMMENT ON COLUMN avatars.raw IS '原始';
COMMENT ON COLUMN avatars.user_id IS '用户ID';
COMMENT ON COLUMN avatars.created_at IS '创建时间';
COMMENT ON COLUMN avatars.updated_at IS '更新时间';

CREATE INDEX idx_avatars_user_id ON avatars (user_id);
