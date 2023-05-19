CREATE TABLE avatars
(
    hash       CHAR(32) PRIMARY KEY NOT NULL COMMENT '哈希',
    raw        VARCHAR(255) UNIQUE COMMENT '原始',
    user_id    BIGINT  DEFAULT NULL COMMENT '用户ID',
    ban        TINYINT(1) DEFAULT '0' COMMENT '禁用',
    checked    TINYINT(1) DEFAULT '0' COMMENT '已检查',
    created_at TIMESTAMP NOT NULL COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL COMMENT '更新时间'
);

CREATE INDEX idx_avatars_user_id ON avatars (user_id);
