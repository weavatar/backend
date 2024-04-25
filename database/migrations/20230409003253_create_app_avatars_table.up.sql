CREATE TABLE app_avatars
(
    id            TEXT PRIMARY KEY NOT NULL,
    app_id        TEXT             NOT NULL,
    avatar_sha256 TEXT             NOT NULL,
    created_at    TIMESTAMP(0)     NOT NULL,
    updated_at    TIMESTAMP(0)     NOT NULL
);

COMMENT ON TABLE app_avatars IS '应用头像';
COMMENT ON COLUMN app_avatars.id IS 'ID';
COMMENT ON COLUMN app_avatars.app_id IS '应用ID';
COMMENT ON COLUMN app_avatars.avatar_sha256 IS '头像SHA256';
COMMENT ON COLUMN app_avatars.created_at IS '创建时间';
COMMENT ON COLUMN app_avatars.updated_at IS '更新时间';

CREATE INDEX idx_app_avatars_app_id ON app_avatars (app_id);
CREATE INDEX idx_app_avatars_avatar_hash ON app_avatars (avatar_sha256);
