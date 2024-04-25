CREATE TABLE avatars
(
    sha256     TEXT PRIMARY KEY NOT NULL,
    md5        TEXT UNIQUE      NOT NULL,
    raw        TEXT UNIQUE      NOT NULL,
    user_id    TEXT             NOT NULL,
    created_at TIMESTAMP(0)     NOT NULL,
    updated_at TIMESTAMP(0)     NOT NULL
);

COMMENT ON TABLE avatars IS '头像';
COMMENT ON COLUMN avatars.sha256 IS 'SHA256';
COMMENT ON COLUMN avatars.md5 IS 'MD5';
COMMENT ON COLUMN avatars.raw IS '原始';
COMMENT ON COLUMN avatars.user_id IS '用户ID';
COMMENT ON COLUMN avatars.created_at IS '创建时间';
COMMENT ON COLUMN avatars.updated_at IS '更新时间';

CREATE INDEX idx_avatars_user_id ON avatars (user_id);
