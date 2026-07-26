CREATE TABLE `post` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    `post_id` BIGINT UNSIGNED NOT NULL COMMENT '帖子ID',
    `title` VARCHAR(128) NOT NULL COMMENT '帖子标题',
    `author_id` BIGINT UNSIGNED NOT NULL COMMENT '作者用户ID',
    `community_id` BIGINT UNSIGNED NOT NULL COMMENT '所属社区ID',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '帖子状态',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uk_post_id` (`post_id`),
    KEY `idx_author_id` (`author_id`),
    KEY `idx_community_id` (`community_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `post_content` (
    `post_id` BIGINT UNSIGNED PRIMARY KEY COMMENT '帖子ID',
    `content` TEXT NOT NULL COMMENT '帖子内容'
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
