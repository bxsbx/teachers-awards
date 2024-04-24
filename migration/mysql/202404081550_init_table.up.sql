/*
 Navicat Premium Data Transfer

 Source Server         : TeachersAwards
 Source Server Type    : MySQL
 Source Server Version : 50732
 Source Host           : rm-wz976878216hxs1yl0o.mysql.rds.aliyuncs.com:3306
 Source Schema         : teachers_awards

 Target Server Type    : MySQL
 Target Server Version : 50732
 File Encoding         : 65001

 Date: 08/04/2024 15:51:33
*/

ALTER DATABASE `teachers_awards` CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_general_ci';

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for activity
-- ----------------------------

CREATE TABLE `activity`  (
  `activity_id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `activity_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '活动名称',
  `year` int(4) UNSIGNED NOT NULL COMMENT '年份',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '申报须知',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '活动文件url',
  `review_num` tinyint(4) NOT NULL COMMENT '审核人数',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime NOT NULL COMMENT '结束时间',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`activity_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 79 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '名师评优活动' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for activity_one_indicator
-- ----------------------------

CREATE TABLE `activity_one_indicator`  (
  `activity_id` int(11) UNSIGNED NOT NULL COMMENT '活动id',
  `one_indicator_id` int(11) UNSIGNED NOT NULL COMMENT '一级指标id',
  `one_indicator_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '一级指标名称',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '评分标准说明',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`activity_id`, `one_indicator_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '活动下的指标' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for activity_two_indicator
-- ----------------------------

CREATE TABLE `activity_two_indicator`  (
  `activity_id` int(11) UNSIGNED NOT NULL COMMENT '活动id',
  `two_indicator_id` int(11) UNSIGNED NOT NULL COMMENT '二级指标id',
  `two_indicator_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '二级指标名称',
  `score` int(11) NOT NULL COMMENT '分值',
  `one_indicator_id` int(11) NOT NULL COMMENT '所属一级指标',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`activity_id`, `two_indicator_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for announcement
-- ----------------------------

CREATE TABLE `announcement`  (
  `announcement_id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '标题',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容',
  `annex` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '附件链接',
  `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户id',
  `user_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名称',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`announcement_id`) USING BTREE,
  INDEX `index_user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 55 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '公告' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for expert_auth_indicator
-- ----------------------------

CREATE TABLE `expert_auth_indicator`  (
  `user_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户id',
  `two_indicator_id` int(11) NOT NULL COMMENT '二级指标id',
  PRIMARY KEY (`user_id`, `two_indicator_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '专家授权指标' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for judges_verify
-- ----------------------------

CREATE TABLE `judges_verify`  (
  `user_indicator_pass_id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_activity_indicator_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户活动申报二级指标id',
  `user_activity_id` bigint(20) NOT NULL COMMENT '用户活动申报id',
  `judges_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '评委id',
  `judges_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '评委姓名',
  `judges_type` tinyint(2) UNSIGNED NOT NULL COMMENT '评委类型，1：学校，2：专家，3：教育局',
  `is_pass` tinyint(1) UNSIGNED NOT NULL COMMENT '0：未通过，1：通过',
  `score` int(11) NOT NULL COMMENT '得分',
  `opinion` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '审核意见',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `delete_time` datetime NULL DEFAULT NULL COMMENT '删除时间',
  `delete_at` bigint(20) UNSIGNED NOT NULL COMMENT '删除时间戳（毫秒）',
  PRIMARY KEY (`user_indicator_pass_id`) USING BTREE,
  UNIQUE INDEX `index_unique`(`user_activity_indicator_id`, `judges_id`, `judges_type`, `delete_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 266 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '评委审核' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for notice
-- ----------------------------

CREATE TABLE `notice`  (
  `notice_id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `notice_type` tinyint(2) UNSIGNED NOT NULL COMMENT '通知类型，1：系统通知',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '通知内容',
  `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '接收者id',
  `is_read` tinyint(1) UNSIGNED NOT NULL COMMENT '0：未读，1：已读',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`notice_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '消息通知' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for one_indicator
-- ----------------------------

CREATE TABLE `one_indicator`  (
  `one_indicator_id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `one_indicator_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '一级指标名称',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '评分标准说明',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`one_indicator_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 122 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '一级指标' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for operation_record
-- ----------------------------

CREATE TABLE `operation_record`  (
  `operation_id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '操作id',
  `relational_id` bigint(20) UNSIGNED NOT NULL COMMENT '关联的id',
  `relational_type` tinyint(2) UNSIGNED NOT NULL COMMENT '关联表，1：user_activity_indicator',
  `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户id',
  `user_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户姓名',
  `operation_role` tinyint(2) UNSIGNED NOT NULL COMMENT '角色，1：学校，2：专家，3：教育局，4：教师',
  `operation_type` tinyint(2) UNSIGNED NOT NULL COMMENT '操作类型，1：添加，2：修改，3：删除',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '操作说明',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`operation_id`) USING BTREE,
  INDEX `index_relational`(`relational_id`, `relational_type`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 278 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for read
-- ----------------------------

CREATE TABLE `read`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `read_id` int(11) UNSIGNED NOT NULL COMMENT '已读id',
  `read_type` tinyint(2) UNSIGNED NOT NULL COMMENT '1：公告',
  `read_user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '已读用户id',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_index`(`read_id`, `read_type`, `read_user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 70 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '记录用户的公告已读' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for two_indicator
-- ----------------------------

CREATE TABLE `two_indicator`  (
  `two_indicator_id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `two_indicator_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '二级指标名称',
  `score` int(11) NOT NULL COMMENT '分值',
  `one_indicator_id` int(11) NOT NULL COMMENT '所属一级指标',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`two_indicator_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 132 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '二级指标' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_activity
-- ----------------------------

CREATE TABLE `user_activity`  (
  `user_activity_id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `activity_id` int(11) UNSIGNED NOT NULL COMMENT '活动id',
  `user_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户id',
  `user_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户姓名',
  `user_sex` tinyint(1) UNSIGNED NOT NULL COMMENT '1：男，2：女',
  `birthday` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '出生日期',
  `identity_card` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '身份证号',
  `phone` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '手机号',
  `subject_code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '科目code',
  `school_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '学校id',
  `school_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '学校名称',
  `declare_type` tinyint(2) UNSIGNED NOT NULL COMMENT '1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研',
  `final_score` tinyint(4) UNSIGNED NOT NULL COMMENT '最终得分（各项通过的审核）',
  `rank` int(11) UNSIGNED NOT NULL COMMENT '排名，如果不为0：表示为评定结果',
  `rank_prize` tinyint(2) NOT NULL COMMENT '0：无，1：一等奖，2：二等奖，3：三等奖',
  `prize` int(11) UNSIGNED NOT NULL COMMENT '奖金',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime NULL DEFAULT NULL COMMENT '删除时间',
  `delete_at` bigint(20) UNSIGNED NOT NULL COMMENT '删除时间戳（毫秒）',
  PRIMARY KEY (`user_activity_id`) USING BTREE,
  UNIQUE INDEX `index_user_activity`(`user_id`, `activity_id`, `delete_at`) USING BTREE,
  INDEX `index_activity`(`activity_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 38 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户申报活动' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_activity_indicator
-- ----------------------------

CREATE TABLE `user_activity_indicator`  (
  `user_activity_indicator_id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_activity_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户申报活动的id',
  `two_indicator_id` int(11) UNSIGNED NOT NULL COMMENT '二级指标id',
  `award_date` date NOT NULL COMMENT '获奖日期',
  `certificate_type` tinyint(2) NOT NULL COMMENT '1.证书；需要填写证书有效期，2.证明；不需要有效期',
  `certificate_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '证书url',
  `certificate_start_date` date NULL DEFAULT NULL COMMENT '证书有效期——开始时间',
  `certificate_end_date` date NULL DEFAULT NULL COMMENT '证书有效期——结束时间',
  `status` tinyint(2) NOT NULL COMMENT '审核状态说明：\r\n1.已提交：教师提交申报、进入审核流程，均显示已提交状态\r\n2.已通过：仅当教育局通过最终审核，显示已通过，无法进行修改、撤销操作\r\n3.学校未通过：学校审核不通过，审核流程不进入专家审核阶段\r\n4.专家未通过：无论专家审核是否通过，均进入教育局审核阶段\r\n5.教育局未通过：该项不通过\r\n6.结果已评定，未审批',
  `finish_review_num` tinyint(4) UNSIGNED NOT NULL COMMENT '当前已审核人数',
  `review_process` tinyint(2) NOT NULL COMMENT '当前审核进程，1：学校，2：专家，3：教育局，4：结束',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime NULL DEFAULT NULL COMMENT '删除时间',
  `delete_at` bigint(20) UNSIGNED NOT NULL COMMENT '删除时间戳（毫秒）',
  PRIMARY KEY (`user_activity_indicator_id`) USING BTREE,
  INDEX `index_user_activity_id`(`user_activity_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 100 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户申报的活动下的指标' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_info
-- ----------------------------

CREATE TABLE `user_info`  (
  `user_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户id',
  `user_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名称',
  `user_sex` tinyint(1) UNSIGNED NOT NULL COMMENT '1：男，2：女',
  `birthday` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '出生日期',
  `identity_card` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '身份证号',
  `phone` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '手机号',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '头像',
  `subject_code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '科目code',
  `school_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '学校id',
  `school_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '学校名称',
  `declare_type` tinyint(2) UNSIGNED NOT NULL COMMENT '1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研',
  `role` int(11) NOT NULL COMMENT '角色，1：学校，2：专家，4：教育局，8：教师，16：超级管理员，多个则相加',
  `export_auth` tinyint(1) NOT NULL COMMENT '专家是否授权 1：未授权 2：已授权',
  `auth_day` date NULL DEFAULT NULL COMMENT '授权日期',
  `year` int(11) NOT NULL COMMENT '年份',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '记录用户当前基本信息' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
