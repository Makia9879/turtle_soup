-- 创建本地登录用户turtle_soup
CREATE USER 'turtle_soup'@'localhost' IDENTIFIED BY 'WSijURlVdgWOvvuxADdhM5sXtT3J';

-- 创建数据库turtle_soup
CREATE DATABASE IF NOT EXISTS `turtle_soup` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 授予turtle_soup用户对turtle_soup数据库的全部权限
GRANT ALL PRIVILEGES ON `turtle_soup`.* TO 'turtle_soup'@'localhost';

-- 刷新权限
FLUSH PRIVILEGES;

USE turtle_soup;