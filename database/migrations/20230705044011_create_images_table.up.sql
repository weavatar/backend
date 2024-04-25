CREATE TABLE images
(
    hash       TEXT PRIMARY KEY      NOT NULL,
    ban        BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP(0)          NOT NULL,
    updated_at TIMESTAMP(0)          NOT NULL
);

COMMENT ON TABLE images IS '图片';
COMMENT ON COLUMN images.hash IS '哈希';
COMMENT ON COLUMN images.ban IS '禁用';
COMMENT ON COLUMN images.created_at IS '创建时间';
COMMENT ON COLUMN images.updated_at IS '更新时间';
