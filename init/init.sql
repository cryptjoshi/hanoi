
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+07:00";

CREATE TABLE `Settings` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `key` varchar(255) DEFAULT NULL,
  `value` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `Settings` (`id`, `created_at`, `updated_at`, `deleted_at`, `key`, `value`) VALUES
(1, '2024-11-09 00:22:03.286', '2024-11-09 00:22:03.286', NULL, 'baseCurrency', 'USD'),
(2, '2024-11-09 00:22:03.286', '2024-11-09 00:22:03.286', NULL, 'customerCurrency', 'THB'),
(3, '2024-11-09 00:22:03.286', '2024-11-09 00:22:03.286', NULL, 'baseRate', '1'),
(4, '2024-11-09 00:22:03.286', '2024-11-09 00:22:03.286', NULL, 'customerRate', '34.03'),


ALTER TABLE `Settings`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_Settings_deleted_at` (`deleted_at`);

ALTER TABLE `Settings`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;
  
COMMIT;
