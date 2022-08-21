
CREATE TABLE `b_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID,自增主键',
  `mobile` varchar(128) NOT NULL DEFAULT '' COMMENT '手机号',
  `username` varchar(128) NOT NULL DEFAULT '' COMMENT '用户名',
  `create_time` datetime NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '注册时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_mobile` (`mobile`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE `b_user_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID,自增主键',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '用户ID',
  `mobile` varchar(128) NOT NULL DEFAULT '' COMMENT '手机号',
  `username` varchar(128) NOT NULL DEFAULT '' COMMENT '用户名',
  `create_time` datetime NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '注册时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='用户日志表';

