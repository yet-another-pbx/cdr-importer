-- Create CDR Importer tables

CREATE TABLE `cdr_importer`.`voip_carrier` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  UNIQUE (`name`),
  PRIMARY KEY(`id`)
)
CHARACTER SET utf8mb4
COLLATE utf8mb4_bin
ENGINE = InnoDB;

CREATE TABLE `cdr_importer`.`yap_cdr_insert_log` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
  `voip_carrier_id` BIGINT UNSIGNED NOT NULL,
  `cdr_month_year` VARCHAR(255) NOT NULL,
  `date_time_added` DATETIME DEFAULT NOW() NOT NULL,
  PRIMARY KEY(`id`)
)
CHARACTER SET utf8mb4
COLLATE utf8mb4_bin
ENGINE = InnoDB;
