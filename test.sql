DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS phone;

CREATE TABLE IF NOT EXISTS `users`(
  `id` varchar(16) PRIMARY KEY,
  `email` varchar(256) UNIQUE NOT NULL,
  `name` varchar(255) NOT NULL,
  `age` integer,
  `password` varchar(32) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `phone`(
  `id` varchar(16) PRIMARY KEY,
  `phone` varchar(12) NOT NULL,
  `description` text,
  `user_id` varchar(16) NOT NULL,
  `is_fax` bool NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT `fk_phone_users` FOREIGN KEY(`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE RESTRICT
);

CREATE INDEX `users_index_0` ON `users` (`name`);

CREATE INDEX `phone_index_1` ON `phone` (`phone`);


