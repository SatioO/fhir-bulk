-- DROP DATABASE `fhir`;

DROP TABLE `auth_servers`;
DROP TABLE `fhir_resources`;
DROP TABLE `fhir_jobs`;
DROP TABLE `fhir_apps`;

CREATE TABLE `fhir_apps` (
  `id` VARCHAR(50) NOT NULL,
  `base_url` VARCHAR(255) NOT NULL,
  `token` VARCHAR(1000) NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  PRIMARY KEY(`id`)
);
CREATE TABLE `auth_servers` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `token_url` VARCHAR(255) NULL,
  `client_id` VARCHAR(255) NOT NULL,
  `client_secret` VARCHAR(255) NOT NULL,
  `scopes` VARCHAR(1000) NOT NULL,
  `app_id` VARCHAR(50) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  PRIMARY KEY(`id`),
  FOREIGN KEY(`app_id`) REFERENCES `fhir_apps`(`id`)
);
CREATE TABLE `fhir_jobs` (
  `id` VARCHAR(255) NOT NULL,
  `status` VARCHAR(255) NOT NULL,
  `app_id` VARCHAR(50) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  PRIMARY KEY(`id`),
  FOREIGN KEY(`app_id`) REFERENCES `fhir_apps`(`id`)
);
CREATE TABLE `fhir_resources` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `app_id` VARCHAR(50) NOT NULL,
  `job_id` VARCHAR(255) NOT NULL,
  `resource_id` VARCHAR(255) NOT NULL,
  `type` VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  PRIMARY KEY(`id`),
  FOREIGN KEY(`app_id`) REFERENCES `fhir_apps`(`id`),
  FOREIGN KEY(`job_id`) REFERENCES `fhir_jobs`(`id`)
);

SHOW TABLES;

-- INSERT INTO `fhir_jobs` (`id`, `status`, `app_id`) VALUES ('11ef-660b-04e7e2ec-9ba0-36bf4ca24292', 'submitted', 'cerner');

-- INSERT INTO `fhir_apps` (`id`, `base_url`, `token`, `status`) VALUES ('CERNER', 'https://fhir-ehr-code.cerner.com/r4/ec2458f2-1e24-41c8-b71b-0e701af7583d', '', 'active');
-- INSERT INTO `fhir_apps` (`id`, `base_url`, `token`, `status`) VALUES ('EPIC', 'https://fhir.epic.com/interconnect-fhir-oauth/api/FHIR/DSTU2', '', 'inactive');
-- SELECT * FROM `fhir_apps`
-- SELECT * FROM `auth_servers`
-- DELETE FROM `fhir_apps` WHERE id='CERNER';