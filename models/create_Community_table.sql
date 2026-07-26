DROP TABLE IF EXISTS `community`;

CREATE TABLE `community`(
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `community_id` BIGINT UNSIGNED NOT NULL,
    `community_name` VARCHAR(64) NOT NULL,
    `description` VARCHAR(256),
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uk_community_id` (`community_id`),
    UNIQUE KEY `uk_community_name` (`community_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;