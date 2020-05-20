SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';


-- -----------------------------------------------------
-- Schema foody_auth
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `foody_auth` DEFAULT CHARACTER SET utf8 ;
USE `foody_auth` ;


-- -----------------------------------------------------
-- Table `foody_auth`.`client`
-- -----------------------------------------------------
CREATE TABLE `client` (
  `id` varchar(80) NOT NULL,
  `name` varchar(255) NOT NULL,
  `secret` varchar(80) NOT NULL,
  `type` enum('public','confidential') NOT NULL DEFAULT 'confidential',
  `grant_type` varchar(80) NOT NULL,
  `user_role` enum('customer') NOT NULL,
  `access_token_lifetime` bigint(20) DEFAULT NULL,
  `refresh_token_lifetime` bigint(20) DEFAULT NULL,
  `status` enum('active', 'inactive') NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_client_secret` (`secret`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4;


-- -----------------------------------------------------
-- Table `foody_auth`.`refresh_token`
-- -----------------------------------------------------
CREATE TABLE `refresh_token` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `token` varchar(256) NOT NULL,
  `expiry_date` datetime NOT NULL,
  `client_id` varchar(80) NOT NULL,
  `user_id` varchar(80) NOT NULL,
  `user_role` enum('customer') NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_token` (`token`),
  UNIQUE KEY `unique_client_user` (`client_id`,`user_id`),
  CONSTRAINT `fk_client` FOREIGN KEY (`client_id`) REFERENCES `client` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;