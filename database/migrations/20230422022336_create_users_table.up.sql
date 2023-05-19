CREATE TABLE users
(
    id         BIGINT PRIMARY KEY NOT NULL COMMENT 'ID',
    open_id    CHAR(32) UNIQUE NOT NULL COMMENT 'OpenID',
    union_id   CHAR(32) UNIQUE NOT NULL COMMENT 'UnionID',
    nickname   VARCHAR(255) NOT NULL COMMENT '昵称',
    is_admin   TINYINT(1) DEFAULT '0' COMMENT '是否是管理员',
    real_name  TINYINT(1) DEFAULT '0' COMMENT '是否实名认证',
    created_at TIMESTAMP NOT NULL COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL COMMENT '更新时间',
    deleted_at TIMESTAMP DEFAULT NULL COMMENT '删除时间'
);

CREATE INDEX idx_users_nickname ON users (nickname);
