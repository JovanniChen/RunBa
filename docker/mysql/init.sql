-- RunBa 数据库初始化脚本

CREATE TABLE `forges` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `image_url` VARCHAR(500) DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `swordsmiths` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `image_url` VARCHAR(500) DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `certificates` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `number` VARCHAR(255) NOT NULL,
  `overall_length` VARCHAR(255) DEFAULT NULL,
  `blade_length` VARCHAR(255) DEFAULT NULL,
  `handle_length` VARCHAR(255) DEFAULT NULL,
  `blade_material` VARCHAR(255) DEFAULT NULL,
  `blade_thickness` VARCHAR(255) DEFAULT NULL,
  `kissan_length` VARCHAR(255) DEFAULT NULL,
  `hamon` VARCHAR(255) DEFAULT NULL,
  `mekugi` TINYINT DEFAULT NULL,
  `sakihaba` VARCHAR(255) DEFAULT NULL,
  `motohaba` VARCHAR(255) DEFAULT NULL,
  `saya_material` VARCHAR(255) DEFAULT NULL,
  `tsuba_material` VARCHAR(255) DEFAULT NULL,
  `habaki_material` VARCHAR(255) DEFAULT NULL,
  `ito_sageo_material` VARCHAR(255) DEFAULT NULL,
  `forge_id` BIGINT UNSIGNED NOT NULL,
  `swordsmith_id` BIGINT UNSIGNED NOT NULL,
  `date_completed` VARCHAR(255) DEFAULT NULL,
  `no` VARCHAR(255) DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_certificates_forge_id` (`forge_id`),
  KEY `idx_certificates_swordsmith_id` (`swordsmith_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `username` VARCHAR(50) NOT NULL UNIQUE,
  `password` VARCHAR(255) NOT NULL,
  `nickname` VARCHAR(50) DEFAULT '',
  `status` TINYINT NOT NULL DEFAULT 1,
  `last_login` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME NULL,
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 默认管理员（账号：admin，密码：ChangeMe123!）
INSERT INTO `users`(`username`, `password`, `nickname`, `status`, `created_at`, `updated_at`)
VALUES('admin', '$2a$10$Q7zEbOkd728rwWywoJzxuOTaPoWrhlWdnPP8gemPxNFzNLYIsZ.1m', 'admin', 1, NOW(), NOW());
