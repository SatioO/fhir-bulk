DROP TABLE `auth_server`;
DROP TABLE `fhir_app`;

CREATE TABLE `fhir_app` (
  `id` VARCHAR(50) NOT NULL,
  `base_url` VARCHAR(255) NOT NULL,
  `token` VARCHAR(1000) NULL,
  `status` ENUM('active', 'inactive') NOT NULL DEFAULT 'active',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(`id`)
);
CREATE TABLE `auth_server` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `token_url` VARCHAR(255) NULL,
  `client_id` VARCHAR(255) NOT NULL,
  `client_secret` VARCHAR(255) NOT NULL,
  `app_id` VARCHAR(50) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(`id`),
  FOREIGN KEY(`app_id`) REFERENCES `fhir_app`(`id`)
);

SHOW TABLES;

-- INSERT INTO `fhir_client` (`id`, `base_url`, `token`, `status`) VALUES ('CERNER', 'https://fhir-ehr-code.cerner.com/r4/ec2458f2-1e24-41c8-b71b-0e701af7583d', '', 'active');
-- INSERT INTO `fhir_client` (`id`, `base_url`, `token`, `status`) VALUES ('EPIC', 'https://fhir.epic.com/interconnect-fhir-oauth/api/FHIR/DSTU2', '', 'inactive');
-- SELECT * FROM `fhir_client`
-- SELECT * FROM `auth_server`
-- DELETE FROM `fhir_client` WHERE id='CERNER';