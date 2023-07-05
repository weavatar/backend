CREATE TABLE images
(
    hash       CHAR(32) PRIMARY KEY NOT NULL COMMENT '哈希',
    ban        TINYINT(1)      DEFAULT '0' COMMENT '禁用',
    created_at TIMESTAMP            NOT NULL COMMENT '创建时间',
    updated_at TIMESTAMP            NOT NULL COMMENT '更新时间'
);
