-- Create CDR Importer tables

CREATE TABLE `cdr_importer`.`voip_carrier` (
  `id` SMALLINT UNSIGNED AUTO_INCREMENT NOT NULL,
  `name` VARCHAR(50) NOT NULL,
  `cdr_month_year_column` VARCHAR(9) NOT NULL,
  `rate_card_charge_code_column` VARCHAR(9) NOT NULL,
  `rate_card_price_per_minute_column` VARCHAR(9) NOT NULL,
  `rate_card_price_per_call_column` VARCHAR(9) NOT NULL,
  `date_time_added` DATETIME DEFAULT NOW() NOT NULL,
  UNIQUE (`name`),
  PRIMARY KEY(`id`)
)
CHARACTER SET utf8mb4
COLLATE utf8mb4_bin
ENGINE = InnoDB;

CREATE TABLE `cdr_importer`.`yap_cdr_insert_log` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
  `voip_carrier_id` SMALLINT UNSIGNED NOT NULL,
  `cdr_month_year` VARCHAR(7) NOT NULL,
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

-- Create Views

CREATE VIEW `view___yap_cdr_insert_log_detail` AS
SELECT 
  `voip_carrier`.`id` AS 'voip_carrier_id',
  `voip_carrier`.`name` AS 'voip_carrier_name',
  `yap_cdr_insert_log`.`cdr_month_year`,
  DATE_FORMAT(`yap_cdr_insert_log`.`date_time_added`, '%d/%m/%Y %H:%i:%s') AS 'yap_cdr_insert_log_date_time_added'
FROM `voip_carrier`
INNER JOIN `yap_cdr_insert_log`
ON `voip_carrier`.`id` = `yap_cdr_insert_log`.`voip_carrier_id`;

CREATE VIEW `view___voip_carrier` AS
SELECT 
  `voip_carrier`.`id`,
  `voip_carrier`.`name`,
  DATE_FORMAT(`voip_carrier`.`date_time_added`, '%d/%m/%Y %H:%i:%s') AS 'date_time_added'
FROM `voip_carrier`;
