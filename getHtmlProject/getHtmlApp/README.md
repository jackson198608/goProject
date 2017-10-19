# introduction
- this package is design for the situation when you want to get the request body save new file using abuyun proxy

# Installation

    go get github.com/jackson198608/goProject/getHtml 

# Run

- run ask
    getHtml asksave /etc/configask.json all/update  (all: all of data. update:newest 7 days data)
    getHtml ask /etc/configask.json
- run thread 
    getHtml threadsave /etc/configthread.json all/update  (all: all of data. update:newest 7 days data)
    getHtml thread /etc/configthread.json

# database : process 
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