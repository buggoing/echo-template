CREATE DATABASE IF NOT EXISTS echo_template;

use echo_template;
-- 用户相关的表
CREATE TABLE IF NOT EXISTS user_account (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT, 
    username VARCHAR(32) NOT NULL COLLATE utf8mb4_general_ci COMMENT "用户名",
    passwd VARCHAR(32) NOT NULL COMMENT "用户密码",
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT "创建时间戳",
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "更新时间戳"
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE utf8_general_ci;

INSERT INTO user_account SET username = "username", passwd = "passwd";