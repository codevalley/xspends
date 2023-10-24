CREATE DATABASE IF NOT EXISTS xspends;

USE xspends;

CREATE TABLE IF NOT EXISTS `users` (
    `id` INT NOT NULL,
    `username` VARCHAR(255) NOT NULL UNIQUE,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `currency` VARCHAR(10) DEFAULT 'USD',
    `password` VARCHAR(255) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `categories` (
    `id` INT NOT NULL,
    `user_id` INT NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `icon` VARCHAR(255),
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `sources` (
    `id` BIGINT NOT NULL,
    `user_id` BIGINT NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `type` VARCHAR(255) NOT NULL, 
    `balance` DECIMAL(10, 2) DEFAULT 0.00,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `tags` (
    `id` INT NOT NULL,
    `user_id` INT NOT NULL,
    `name` VARCHAR(255) NOT NULL UNIQUE,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `transactions` (
    `id` INT NOT NULL,
    `user_id` INT NOT NULL,
    `type` VARCHAR(255) NOT NULL DEFAULT 'SAVINGS',
    `source_id` INT,
    `amount` DECIMAL(10, 2) NOT NULL,
    `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `category_id` INT,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
    FOREIGN KEY (`source_id`) REFERENCES `sources`(`id`),
    FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`),
    PRIMARY KEY (`id`)
);

-- Create a transaction_tags junction table to store the many-to-many relationship between transactions and tags
CREATE TABLE IF NOT EXISTS `transaction_tags` (
    `transaction_id` INT NOT NULL,
    `tag_id` INT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`transaction_id`) REFERENCES `transactions`(`id`),
    FOREIGN KEY (`tag_id`) REFERENCES `tags`(`id`),
    PRIMARY KEY (`transaction_id`, `tag_id`)
);

-- Indexes for optimized user-specific lookups
CREATE INDEX idx_categories_userid ON categories(user_id);
CREATE INDEX idx_sources_userid ON sources(user_id);
CREATE INDEX idx_tags_userid ON tags(user_id);
