DROP INDEX idx_users_nickname;

CREATE TABLE users_temp
(
    id         BIGINT PRIMARY KEY         NOT NULL,
    open_id    CHAR(32) UNIQUE            NOT NULL,
    union_id   CHAR(32) UNIQUE            NOT NULL,
    nickname   VARCHAR(255)               NOT NULL,
    avatar     VARCHAR(255) DEFAULT NULL,
    is_admin   BOOLEAN      DEFAULT FALSE NOT NULL,
    real_name  BOOLEAN      DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP                  NOT NULL,
    updated_at TIMESTAMP                  NOT NULL,
    deleted_at TIMESTAMP    DEFAULT NULL
);

COMMENT ON COLUMN users_temp.id IS 'ID';
COMMENT ON COLUMN users_temp.open_id IS 'OpenID';
COMMENT ON COLUMN users_temp.union_id IS 'UnionID';
COMMENT ON COLUMN users_temp.nickname IS '昵称';
COMMENT ON COLUMN users_temp.avatar IS '头像';
COMMENT ON COLUMN users_temp.is_admin IS '是否是管理员';
COMMENT ON COLUMN users_temp.real_name IS '是否实名认证';
COMMENT ON COLUMN users_temp.created_at IS '创建时间';
COMMENT ON COLUMN users_temp.updated_at IS '更新时间';
COMMENT ON COLUMN users_temp.deleted_at IS '删除时间';

CREATE INDEX idx_users_nickname ON users_temp (nickname);

INSERT INTO "users_temp" ("id", "open_id", "union_id", "nickname", "is_admin", "real_name", "created_at", "updated_at",
                          "deleted_at")
SELECT "id",
       "open_id",
       "union_id",
       "nickname",
       "is_admin",
       "real_name",
       "created_at",
       "updated_at",
       "deleted_at"
FROM "users";

DROP TABLE "users";

ALTER TABLE "users_temp"
    RENAME TO "users";

UPDATE users
SET avatar = 'https://weavatar.com/avatar/?d=mp'
WHERE avatar IS NULL;

ALTER TABLE users
    ALTER COLUMN avatar SET NOT NULL;
ALTER TABLE users
    ALTER COLUMN avatar DROP DEFAULT;
