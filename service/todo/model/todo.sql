CREATE TABLE `todo` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `title` varchar(50)  NOT NULL DEFAULT '',
    `content` varchar(300)  NOT NULL DEFAULT '',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `done` boolean NOT NULL DEFAULT FALSE,
    `user_id` bigint NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `title_unique` (`title`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;