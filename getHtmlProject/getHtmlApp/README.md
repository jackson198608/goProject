database : process 
execute record
CREATE TABLE `execute_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `process_name` varchar(250) DEFAULT NULL COMMENT '进程名',
  `data_source` varchar(50) DEFAULT NULL COMMENT '数据来源',
  `platform` varchar(50) DEFAULT NULL COMMENT '展示平台',
  `created` datetime NOT NULL COMMENT '创建时间',
  `lastdate` datetime NOT NULL COMMENT '最新执行时间',
  PRIMARY KEY (`id`),
  KEY `lastdate` (`lastdate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;