/*
Navicat MySQL Data Transfer

Source Server         : ossfile
Source Server Version : 50726
Source Host           : localhost:3306
Source Database       : oss_filemeta

Target Server Type    : MYSQL
Target Server Version : 50726
File Encoding         : 65001

Date: 2024-06-12 10:46:45
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `file_meta`
-- ----------------------------
DROP TABLE IF EXISTS `file_meta`;
CREATE TABLE `file_meta` (
  `fid` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) DEFAULT NULL,
  `dir` bigint(20) DEFAULT NULL,
  `hash` longtext,
  `name` longtext,
  `size` bigint(20) DEFAULT NULL,
  `content_type` longtext,
  `upload_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  `is_del` tinyint(1) DEFAULT NULL,
  `is_creating` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`fid`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of file_meta
-- ----------------------------

-- ----------------------------
-- Table structure for `oss_meta`
-- ----------------------------
DROP TABLE IF EXISTS `oss_meta`;
CREATE TABLE `oss_meta` (
  `hash` varchar(191) NOT NULL,
  `size` bigint(20) DEFAULT NULL,
  `is_del` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of oss_meta
-- ----------------------------

-- ----------------------------
-- Table structure for `oss_user`
-- ----------------------------
DROP TABLE IF EXISTS `oss_user`;
CREATE TABLE `oss_user` (
  `uid` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(191) DEFAULT NULL,
  `password` longtext,
  `nickname` longtext,
  `email` longtext,
  `avatar` longtext,
  `isAdmin` bigint(20) DEFAULT '0',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `uni_oss_user_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of oss_user
-- ----------------------------

-- ----------------------------
-- Table structure for `resource_usage`
-- ----------------------------
DROP TABLE IF EXISTS `resource_usage`;
CREATE TABLE `resource_usage` (
  `rid` bigint(20) NOT NULL AUTO_INCREMENT,
  `addr` varchar(64) DEFAULT NULL,
  `node_type` bigint(20) DEFAULT NULL,
  `cpu_usage` double DEFAULT NULL,
  `memory_usage` double DEFAULT NULL,
  `disk_usage` double DEFAULT NULL,
  `network_usage` double DEFAULT NULL,
  `create_time` datetime(3) DEFAULT NULL,
  `cpu_current_usage` double DEFAULT NULL,
  `cpu_max_usage` double DEFAULT NULL,
  `memory_current_usage` bigint(20) unsigned DEFAULT NULL,
  `memory_total` bigint(20) unsigned DEFAULT NULL,
  `disk_current_usage` bigint(20) unsigned DEFAULT NULL,
  `disk_total` bigint(20) unsigned DEFAULT NULL,
  `network_sent` bigint(20) unsigned DEFAULT NULL,
  `network_recv` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`rid`),
  KEY `lastestRecord` (`addr`,`create_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of resource_usage
-- ----------------------------
