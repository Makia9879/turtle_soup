-- 撤销turtle_soup用户的所有权限
REVOKE ALL PRIVILEGES ON `turtle_soup`.* FROM 'turtle_soup'@'%';

-- 删除turtle_soup数据库
DROP DATABASE IF EXISTS `turtle_soup`;

-- 删除turtle_soup用户
DROP USER IF EXISTS 'turtle_soup'@'%';

-- 刷新权限
FLUSH PRIVILEGES;