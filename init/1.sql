CREATE TABLE `Settings` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,  -- แก้ไข: เพิ่มประเภทข้อมูลและ AUTO_INCREMENT
  `created_at` datetime(6) DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` datetime(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),  -- แก้ไข: เพิ่ม ON UPDATE สำหรับ updated_at
  `deleted_at` datetime(6) DEFAULT NULL,  -- แก้ไข: เปลี่ยน DEFAULT เป็น NULL
  `key` varchar(255) DEFAULT NULL,
  `value` varchar(255) DEFAULT NULL
);

GRANT ALL ON *.* TO 'web'@'%' WITH GRANT OPTION;