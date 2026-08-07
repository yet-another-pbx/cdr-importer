-- Create CDR Importer tables

CREATE TABLE `cdr_importer`.`voip_carrier` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `date_time_added` DATETIME DEFAULT NOW() NOT NULL,
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

----------------------------------------------------------------------------------------------------

-- Add index

ALTER TABLE `yap_cdr_insert_log`
ADD INDEX `index___yap_cdr_insert_log__voip_carrier_id` (`voip_carrier_id`);

----------------------------------------------------------------------------------------------------

-- Create foreign key constraint

ALTER TABLE `yap_cdr_insert_log`
ADD CONSTRAINT fk___yap_cdr_insert_log___voip_carrier
FOREIGN KEY (`voip_carrier_id`)
REFERENCES `voip_carrier` (`id`)
ON UPDATE CASCADE;

----------------------------------------------------------------------------------------------------

-- Create View

CREATE VIEW `view___yap_cdr_insert_log_detail` AS
SELECT 
  `voip_carrier`.`id` AS 'voip_carrier_id',
  `voip_carrier`.`name` AS 'voip_carrier_name',
  `yap_cdr_insert_log`.`cdr_month_year`,
  `yap_cdr_insert_log`.`date_time_added` AS 'yap_cdr_insert_log_date_time_added'
FROM `voip_carrier`
INNER JOIN `yap_cdr_insert_log`
ON `voip_carrier`.`id` = `yap_cdr_insert_log`.`voip_carrier_id`;
