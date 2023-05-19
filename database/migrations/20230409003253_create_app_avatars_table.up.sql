CREATE TABLE app_avatars
(
    id          BIGINT PRIMARY KEY NOT NULL COMMENT 'ID',
    app_id      BIGINT NOT NULL COMMENT '应用ID',
    avatar_hash CHAR(32) NOT NULL COMMENT '头像哈希',
    ban         TINYINT(1) DEFAULT '0' COMMENT '禁用',
    checked     TINYINT(1) DEFAULT '0' COMMENT '已检查',
    created_at  TIMESTAMP NOT NULL COMMENT '创建时间',
    updated_at  TIMESTAMP NOT NULL COMMENT '更新时间'
);

CREATE INDEX idx_app_avatars_app_id ON app_avatars (app_id);
CREATE INDEX idx_app_avatars_avatar_hash ON app_avatars (avatar_hash);
