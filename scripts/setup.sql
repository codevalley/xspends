CREATE DATABASE IF NOT EXISTS xspends;

USE xspends;

CREATE TABLE IF NOT EXISTS `users` (
    `user_id` BIGINT NOT NULL,
    `username` VARCHAR(255) NOT NULL UNIQUE,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `currency` VARCHAR(10) DEFAULT 'USD',
    `scope_id` BIGINT NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`)
);


CREATE TABLE IF NOT EXISTS `scopes` (
    `scope_id` BIGINT NOT NULL,
    `type` VARCHAR(64) NOT NULL, 
    PRIMARY KEY (`scope_id`)
);

CREATE TABLE IF NOT EXISTS `user_scopes` (
    `user_id` BIGINT NOT NULL,
    `scope_id` BIGINT NOT NULL,
    `role` VARCHAR(64) DEFAULT 'view', 
    FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`),
    FOREIGN KEY (`scope_id`) REFERENCES `scopes`(`scope_id`),
    PRIMARY KEY (`user_id`, `scope_id`)
);

CREATE TABLE IF NOT EXISTS `user_groups` (
    `group_id` BIGINT NOT NULL,
    `owner_id` BIGINT NOT NULL,
    `scope_id` BIGINT NOT NULL,
    `group_name` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `icon` VARCHAR(255),
    `status` VARCHAR(64) NOT NULL DEFAULT 'active', 
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`group_id`),
    FOREIGN KEY (`owner_id`) REFERENCES `users`(`user_id`),
    FOREIGN KEY (`scope_id`) REFERENCES `scopes`(`scope_id`)
);

CREATE TABLE IF NOT EXISTS `categories` (
    `category_id` BIGINT NOT NULL,
    `user_id` BIGINT NOT NULL,
    `scope_id` BIGINT NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `icon` VARCHAR(255),
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`),
    PRIMARY KEY (`category_id`),
    UNIQUE (`user_id`, `name`),
    FOREIGN KEY (`scope_id`) REFERENCES `scopes`(`scope_id`)
);

CREATE TABLE IF NOT EXISTS `sources` (
    `source_id` BIGINT NOT NULL,
    `user_id` BIGINT NOT NULL,
    `scope_id` BIGINT NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `type` VARCHAR(64) NOT NULL,  
    `balance` DECIMAL(10, 2) DEFAULT 0.00,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`),
    PRIMARY KEY (`source_id`),
    FOREIGN KEY (`scope_id`) REFERENCES `scopes`(`scope_id`),
    UNIQUE (`user_id`, `name`)
);

CREATE TABLE IF NOT EXISTS `tags` (
    `tag_id` BIGINT NOT NULL,
    `user_id` BIGINT NOT NULL,
    `scope_id` BIGINT NOT NULL,
    `name` VARCHAR(255) NOT NULL UNIQUE,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`),
    PRIMARY KEY (`tag_id`),
    FOREIGN KEY (`scope_id`) REFERENCES `scopes`(`scope_id`),
    UNIQUE (`user_id`, `name`)
);

CREATE TABLE IF NOT EXISTS `transactions` (
    `transaction_id` BIGINT NOT NULL,
    `user_id` BIGINT NOT NULL,
    `type` VARCHAR(255) NOT NULL DEFAULT 'SAVINGS',
    `scope_id` BIGINT NOT NULL,
    `source_id` BIGINT,
    `amount` DECIMAL(10, 2) NOT NULL,
    `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `category_id` BIGINT,
    `description` TEXT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`),
    FOREIGN KEY (`source_id`) REFERENCES `sources`(`source_id`),
    FOREIGN KEY (`scope_id`) REFERENCES `scopes`(`scope_id`),
    FOREIGN KEY (`category_id`) REFERENCES `categories`(`category_id`),
    PRIMARY KEY (`transaction_id`)
);

CREATE TABLE IF NOT EXISTS `transaction_tags` (
    `transaction_id` BIGINT NOT NULL,
    `tag_id` BIGINT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`transaction_id`) REFERENCES `transactions`(`transaction_id`),
    FOREIGN KEY (`tag_id`) REFERENCES `tags`(`tag_id`),
    PRIMARY KEY (`transaction_id`, `tag_id`)
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_transactions_userid ON transactions(user_id);
CREATE INDEX idx_categories_userid ON categories(user_id);
CREATE INDEX idx_sources_userid ON sources(user_id);
CREATE INDEX idx_tags_userid ON tags(user_id);
