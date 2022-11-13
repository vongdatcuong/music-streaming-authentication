CREATE TABLE User (
	`user_id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `email` VARCHAR(256) NOT NULL,
    `first_name` VARCHAR(256) NOT NULL,
    `last_name` VARCHAR(256) NOT NULL,
    `status` TINYINT UNSIGNED NOT NULL,
    `password` CHAR(60) NOT NULL,
    `new_song_noti` BOOLEAN DEFAULT FALSE,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,

    PRIMARY KEY (`user_id`),
    UNIQUE (`email`)
);

CREATE TABLE Permission (
	`permission_id` INT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(256) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `status` TINYINT UNSIGNED NOT NULL,

    PRIMARY KEY (`permission_id`),
    UNIQUE (`name`)
);

CREATE TABLE User_Permission (
	`permission_id` INT(20) UNSIGNED NOT NULL,
    `user_id` BIGINT(20) UNSIGNED NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,

    PRIMARY KEY (`permission_id`, `user_id`),
    FOREIGN KEY (`permission_id`) REFERENCES Permission(`permission_id`),
    FOREIGN KEY (`user_id`) REFERENCES User(`user_id`)
);