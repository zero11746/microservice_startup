DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
                        `id` bigint NOT NULL COMMENT '用户唯一ID（雪花算法生成）',
                        `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
                        `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
                        `password` varchar(100) DEFAULT NULL COMMENT '加密存储的密码（bcrypt算法）',
                        `school_id` int unsigned DEFAULT NULL COMMENT '所属学校',
                        `major_id` int unsigned DEFAULT NULL COMMENT '所属专业ID',
                        `admission_grade` int DEFAULT NULL COMMENT '入学年级（如2021、2022）',
                        `avatar_url` varchar(255) DEFAULT 'default_avatar.png' COMMENT '头像URL',
                        `experience` int NOT NULL DEFAULT 0 COMMENT '经验值',
                        `status` tinyint NOT NULL DEFAULT 1 COMMENT '账号状态（1-正常，0-封禁，2-注销）',
                        `logoff_time` int unsigned DEFAULT 0 COMMENT '注销时间',
                        `banned_time` int unsigned DEFAULT 0 COMMENT '封禁时间',
                        `unbanned_time` int unsigned DEFAULT 0 COMMENT '解封时间',
                        `created_at` int unsigned COMMENT '注册时间（UNIX时间戳）',
                        `updated_at` int unsigned COMMENT '更新时间（UNIX时间戳）',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `uk_phone_status_logoff_time` (`phone`, `status`, `logoff_time`),
                        KEY `idx_school_id` (`school_id`) COMMENT '按学校查询用户的索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户表';

DROP TABLE IF EXISTS `school`;
CREATE TABLE `school` (
                          `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '学校ID',
                          `name` varchar(100) NOT NULL COMMENT '学校全称',
                          `logo_url` varchar(255) NOT NULL DEFAULT '' COMMENT '学校Logo图片URL',
                          `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态（1：启用，0：未启用）',
                          `created_at` int unsigned COMMENT '创建时间（UNIX时间戳）',
                          `updated_at` int unsigned COMMENT '更新时间（UNIX时间戳）',
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `uk_name` (`name`) COMMENT '学校名称唯一'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '学校表';

DROP TABLE IF EXISTS `major`;
CREATE TABLE `major` (
                         `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '专业ID',
                         `name` varchar(100) NOT NULL COMMENT '专业名称',
                         `description` varchar(255) DEFAULT NULL COMMENT '专业描述',
                         `created_at` int unsigned COMMENT '创建时间（UNIX时间戳）',
                         `updated_at` int unsigned COMMENT '更新时间（UNIX时间戳）',
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `uk_name` (`name`) COMMENT '专业名称全局唯一'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '专业表';

DROP TABLE IF EXISTS `school_major`;
CREATE TABLE `school_major` (
                                `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '关联ID',
                                `school_id` int unsigned NOT NULL COMMENT '学校ID',
                                `major_id` int unsigned NOT NULL COMMENT '专业ID',
                                `created_at` int unsigned COMMENT '创建时间（UNIX时间戳）',
                                PRIMARY KEY (`id`),
                                UNIQUE KEY `uk_school_major` (`school_id`, `major_id`) COMMENT '学校与专业的组合唯一',
                                KEY `idx_school_id` (`school_id`),
                                KEY `idx_major_id` (`major_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '学校-专业关联表';