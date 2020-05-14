DROP TABLE IF EXISTS chanxian.t_http_reply_rule;
CREATE TABLE chanxian.t_http_reply_rule(
    `Frule_id` int(11) NOT NULL COMMENT '规则ID',
    `Frequest_uri` varchar(250) NOT NULL COMMENT '请求的uri',
    `Fresponse` TEXT COMMENT '响应的内容',
    `Fcheck` TEXT COMMENT '根据请求报文返回相应内容',
    `Fremark` varchar(250) NOT NULL DEFAULT '' COMMENT '规则备注',
    `Fcreate_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `Fupdate_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`Frule_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='http回复规则表';

