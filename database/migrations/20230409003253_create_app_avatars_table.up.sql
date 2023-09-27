CREATE TABLE app_avatars
(
    id          BIGINT PRIMARY KEY NOT NULL,
    app_id      BIGINT             NOT NULL,
    avatar_hash CHAR(32)           NOT NULL,
    ban         BOOLEAN DEFAULT '0',
    checked     BOOLEAN DEFAULT '0',
    created_at  TIMESTAMP(3)       NOT NULL,
    updated_at  TIMESTAMP(3)       NOT NULL
);

COMMENT ON COLUMN app_avatars.id IS 'ID';
COMMENT ON COLUMN app_avatars.app_id IS '应用ID';
COMMENT ON COLUMN app_avatars.avatar_hash IS '头像哈希';
COMMENT ON COLUMN app_avatars.ban IS '禁用';
COMMENT ON COLUMN app_avatars.checked IS '已检查';
COMMENT ON COLUMN app_avatars.created_at IS '创建时间';
COMMENT ON COLUMN app_avatars.updated_at IS '更新时间';

CREATE INDEX idx_app_avatars_app_id ON app_avatars (app_id);
CREATE INDEX idx_app_avatars_avatar_hash ON app_avatars (avatar_hash);
