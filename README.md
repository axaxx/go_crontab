# go_crontab
golang的定时项目(不用停机,热更新),模拟的Linux的crontab 
说明:
1.定时任务规则保存在数据库中(附录有表结构)
2.提供接口来实时增加新任务,add-增加实时任务链,build-编译新扩展
3.代码中有扩展示例 
4.第一版本只实现了诸如如下形式的定时任务 * * * * 或者x x x x 分 时 日 月 语义和Linux相同,不包含周

使用:配置数据库
用户名密码自己改,改完之后在config.go文件调整
配置库示例
CREATE TABLE IF NOT EXISTS `crontab_plugins_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(256) NOT NULL COMMENT '扩展名',
  `plugin` varchar(256) NOT NULL COMMENT '扩展标志',
  `plugin_path` varchar(1024) NOT NULL COMMENT '扩展路径',
  `status` varchar(64) NOT NULL COMMENT '数据状态',
  `exec_time` varchar(1024) NOT NULL COMMENT '执行时间',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='配置表' AUTO_INCREMENT=2 ;

INSERT INTO `crontab_plugins_config` (`id`, `name`, `plugin`, `plugin_path`, `status`, `exec_time`, `create_time`, `update_time`) VALUES
(1, 'ceshi', 'test', '/plugins/test', '1', '1 * * *', '2021-09-20 11:34:27', '2021-09-20 16:40:51');


项目初始化
go mod init go_crontab
go mod tidy

即可
