mysqldump: [Warning] Using a password on the command line interface can be insecure.
-- MySQL dump 10.13  Distrib 9.3.0, for macos13.7 (arm64)
--
-- Host: localhost    Database: sale_dev
-- ------------------------------------------------------
-- Server version	9.3.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
mysqldump: Error: 'Access denied; you need (at least one of) the PROCESS privilege(s) for this operation' when trying to dump tablespaces

--
-- Current Database: `sale_dev`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `sale_dev` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `sale_dev`;

--
-- Table structure for table `c_recharges`
--

DROP TABLE IF EXISTS `c_recharges`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `c_recharges` (
  `id` varchar(64) NOT NULL,
  `member_id` varchar(64) DEFAULT NULL,
  `member_name` longtext,
  `member_phone` varchar(32) DEFAULT NULL,
  `center_id` longtext,
  `center_name` longtext,
  `amount` double DEFAULT NULL,
  `points` bigint DEFAULT NULL,
  `payment_method` longtext,
  `operator_id` longtext,
  `operator_name` longtext,
  `remark` longtext,
  `balance_before` bigint DEFAULT NULL,
  `balance_after` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_c_recharges_member_id` (`member_id`),
  KEY `idx_c_recharges_member_phone` (`member_phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `c_recharges`
--

LOCK TABLES `c_recharges` WRITE;
/*!40000 ALTER TABLE `c_recharges` DISABLE KEYS */;
INSERT INTO `c_recharges` VALUES ('1b887a2b-050d-456f-83c1-9830d50af83d','220855543','测试用户','17615860006','center-bj-cy','北京朝阳中心',10,10,'cash','2','13900000001','验证真实操作员',150,160,'2026-04-17 16:34:15.789'),('2463aac3-e90a-41c1-aa52-22bba2675b30','220855543','17615860006','17615860006','center-bj-cy','北京朝阳中心',100,100,'cash','op123','张出纳','',0,100,'2026-04-17 13:51:08.599'),('868dfe70-5460-4c6a-b8b9-99b2b1b4a48f','220855543','测试用户','17615860006','center-bj-cy','北京朝阳中心',50,50,'cash','op123','张出纳','验证操作员姓名',100,150,'2026-04-17 16:29:32.459');
/*!40000 ALTER TABLE `c_recharges` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `card_issue_records`
--

DROP TABLE IF EXISTS `card_issue_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `card_issue_records` (
  `id` varchar(64) NOT NULL,
  `card_no` varchar(32) NOT NULL COMMENT '卡号',
  `user_id` varchar(64) NOT NULL COMMENT '用户ID',
  `user_phone` varchar(32) NOT NULL COMMENT '用户手机号',
  `issue_reason` varchar(64) NOT NULL COMMENT '发放原因:购买套餐包/推荐奖励/其他',
  `issue_type` tinyint NOT NULL COMMENT '1=实体卡（运营绑定）,2=虚拟卡（用户领取）',
  `recharge_center_id` varchar(64) NOT NULL COMMENT '充值中心ID',
  `operator_id` varchar(64) NOT NULL COMMENT '操作员ID',
  `related_user_phone` varchar(32) DEFAULT '' COMMENT '推荐奖励时关联购买人手机号',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_card_no` (`card_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_recharge_center_id` (`recharge_center_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='门店卡发放记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `card_issue_records`
--

LOCK TABLES `card_issue_records` WRITE;
/*!40000 ALTER TABLE `card_issue_records` DISABLE KEYS */;
INSERT INTO `card_issue_records` VALUES ('6de0d7f7-8c35-48b8-9781-4188389d6a29','TJ00000001','','17615860006','购买套餐包',1,'center-bj-cy','1','','','2026-04-17 14:11:47'),('eaae37df-9cf8-4c61-899b-ac3a1290e4e8','TJ00000002','','17615860006','购买套餐包',1,'center-bj-cy','1','','','2026-04-17 16:43:14'),('fd9909fe-bc35-44e4-ade3-1e131d776f50','TJ00000003','','17615860006','购买套餐包',1,'center-bj-cy','1','','','2026-04-17 16:49:35');
/*!40000 ALTER TABLE `card_issue_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `card_transactions`
--

DROP TABLE IF EXISTS `card_transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `card_transactions` (
  `id` varchar(64) NOT NULL,
  `card_no` varchar(32) NOT NULL COMMENT '卡号',
  `type` varchar(32) NOT NULL COMMENT 'issue/consume/freeze/unfreeze/activate/void',
  `amount` int DEFAULT '0' COMMENT '金额（元）',
  `balance_before` int DEFAULT '0' COMMENT '交易前余额（元）',
  `balance_after` int DEFAULT '0' COMMENT '交易后余额（元）',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `operator_id` varchar(64) DEFAULT '' COMMENT '操作员ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_card_no` (`card_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='门店卡交易记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `card_transactions`
--

LOCK TABLES `card_transactions` WRITE;
/*!40000 ALTER TABLE `card_transactions` DISABLE KEYS */;
INSERT INTO `card_transactions` VALUES ('25de68dd-49b2-41fb-b326-7829b5d900ec','TJ00000001','consume',100,1000,900,'','1','2026-04-17 16:50:08'),('47d61c20-5029-473a-a380-cce40dcede04','TJ00000001','issue',0,1000,1000,'发放给用户 17615860006（购买套餐包）','1','2026-04-17 14:11:47'),('55fcc861-b6a6-4d9a-9db2-4257c53d4e85','TJ00000002','issue',0,1000,1000,'发放给用户 17615860006（购买套餐包）','1','2026-04-17 16:43:14'),('9151db8d-b729-4f50-9055-2d5419192b2a','TJ00000003','issue',0,1000,1000,'发放给用户 17615860006（购买套餐包）','1','2026-04-17 16:49:35');
/*!40000 ALTER TABLE `card_transactions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `casbin_casbin_rule`
--

DROP TABLE IF EXISTS `casbin_casbin_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `casbin_casbin_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) DEFAULT NULL,
  `v0` varchar(100) DEFAULT NULL,
  `v1` varchar(100) DEFAULT NULL,
  `v2` varchar(100) DEFAULT NULL,
  `v3` varchar(100) DEFAULT NULL,
  `v4` varchar(100) DEFAULT NULL,
  `v5` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `casbin_casbin_rule`
--

LOCK TABLES `casbin_casbin_rule` WRITE;
/*!40000 ALTER TABLE `casbin_casbin_rule` DISABLE KEYS */;
INSERT INTO `casbin_casbin_rule` VALUES (44,'g','4','center_admin','','','',''),(45,'g','5','operator','','','',''),(11,'p','center_admin','/api/v1/admin/users','*','','',''),(12,'p','center_admin','/api/v1/admin/users/*','*','','',''),(27,'p','center_admin','/api/v1/card/*/freeze','*','','',''),(28,'p','center_admin','/api/v1/card/*/unfreeze','*','','',''),(29,'p','center_admin','/api/v1/card/*/void','*','','',''),(18,'p','center_admin','/api/v1/card/available','GET','','',''),(19,'p','center_admin','/api/v1/card/available/*','GET','','',''),(26,'p','center_admin','/api/v1/card/bind','*','','',''),(25,'p','center_admin','/api/v1/card/center-stats','GET','','',''),(17,'p','center_admin','/api/v1/card/consume','*','','',''),(21,'p','center_admin','/api/v1/card/detail/*','GET','','',''),(23,'p','center_admin','/api/v1/card/inventory-stats','GET','','',''),(20,'p','center_admin','/api/v1/card/list','GET','','',''),(24,'p','center_admin','/api/v1/card/monthly-trend','GET','','',''),(22,'p','center_admin','/api/v1/card/stats','GET','','',''),(16,'p','center_admin','/api/v1/card/verify/*','*','','',''),(5,'p','center_admin','/api/v1/center','GET','','',''),(6,'p','center_admin','/api/v1/center/*','GET','','',''),(3,'p','center_admin','/api/v1/dashboard/*','GET','','',''),(9,'p','center_admin','/api/v1/operator','GET','','',''),(13,'p','center_admin','/api/v1/recharge/c-entry/*','*','','',''),(14,'p','center_admin','/api/v1/recharge/records/*','GET','','',''),(1,'p','center_admin','/api/v1/user/*','*','','',''),(41,'p','operator','/api/v1/card/*/freeze','*','','',''),(42,'p','operator','/api/v1/card/*/unfreeze','*','','',''),(43,'p','operator','/api/v1/card/*/void','*','','',''),(32,'p','operator','/api/v1/card/available','GET','','',''),(33,'p','operator','/api/v1/card/available/*','GET','','',''),(40,'p','operator','/api/v1/card/bind','*','','',''),(39,'p','operator','/api/v1/card/center-stats','GET','','',''),(31,'p','operator','/api/v1/card/consume','*','','',''),(35,'p','operator','/api/v1/card/detail/*','GET','','',''),(37,'p','operator','/api/v1/card/inventory-stats','GET','','',''),(34,'p','operator','/api/v1/card/list','GET','','',''),(38,'p','operator','/api/v1/card/monthly-trend','GET','','',''),(36,'p','operator','/api/v1/card/stats','GET','','',''),(30,'p','operator','/api/v1/card/verify/*','*','','',''),(7,'p','operator','/api/v1/center','GET','','',''),(8,'p','operator','/api/v1/center/*','GET','','',''),(4,'p','operator','/api/v1/dashboard/*','GET','','',''),(10,'p','operator','/api/v1/operator','GET','','',''),(46,'p','operator','/api/v1/recharge/c-entry/*','*',NULL,NULL,NULL),(15,'p','operator','/api/v1/recharge/records/*','GET','','',''),(2,'p','operator','/api/v1/user/*','*','','','');
/*!40000 ALTER TABLE `casbin_casbin_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `recharge_applications`
--

DROP TABLE IF EXISTS `recharge_applications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `recharge_applications` (
  `id` varchar(64) NOT NULL,
  `center_id` varchar(64) DEFAULT NULL,
  `center_name` longtext,
  `amount` double DEFAULT NULL,
  `points` bigint DEFAULT NULL,
  `base_points` bigint DEFAULT NULL,
  `rebate_points` bigint DEFAULT NULL,
  `rebate_rate` bigint DEFAULT NULL,
  `applicant_id` varchar(64) DEFAULT NULL,
  `applicant_name` longtext,
  `transaction_no` longtext,
  `screenshot` longtext,
  `remark` longtext,
  `status` varchar(32) DEFAULT 'pending',
  `approved_by` longtext,
  `approved_at` datetime(3) DEFAULT NULL,
  `approval_remark` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_recharge_applications_center_id` (`center_id`),
  KEY `idx_recharge_applications_applicant_id` (`applicant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `recharge_applications`
--

LOCK TABLES `recharge_applications` WRITE;
/*!40000 ALTER TABLE `recharge_applications` DISABLE KEYS */;
INSERT INTO `recharge_applications` VALUES ('3c767a06-4808-49b6-a6b0-7254c72df23a','1','北京朝阳中心',2000,2020,2000,20,1,'user123','张财务','x1','','','pending','',NULL,'','2026-04-18 10:39:35.888','2026-04-18 10:39:35.888'),('4aa4d301-b537-4450-89f5-388cdbc1f4a7','center-bj-cy','北京朝阳中心',10000,10100,10000,100,1,'user123','张财务','TEST001','','测试充值','approved','admin','2026-04-17 13:46:59.885','测试审批','2026-04-17 13:37:27.738','2026-04-17 13:46:59.885'),('7592b803-afdb-43de-be84-12bbfc37b0b1','3','上海浦东中心',1000,1010,1000,10,1,'user123','张财务','c1','','','pending','',NULL,'','2026-04-17 14:07:42.589','2026-04-17 14:07:42.589'),('f9388969-79be-42f9-a325-a27350fcf8c1','1','北京朝阳中心',3000,3030,3000,30,1,'2','13900000001-1','1','','','pending','',NULL,'','2026-04-18 11:02:35.535','2026-04-18 11:02:35.535');
/*!40000 ALTER TABLE `recharge_applications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `recharge_centers`
--

DROP TABLE IF EXISTS `recharge_centers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `recharge_centers` (
  `id` varchar(64) NOT NULL,
  `name` varchar(128) DEFAULT NULL,
  `code` varchar(64) DEFAULT NULL,
  `province` varchar(32) DEFAULT NULL,
  `city` varchar(32) DEFAULT NULL,
  `district` varchar(32) DEFAULT NULL,
  `address` longtext,
  `manager_id` varchar(64) DEFAULT NULL,
  `phone` longtext,
  `balance` double DEFAULT '0',
  `status` varchar(32) DEFAULT 'active',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_recharge_centers_name` (`name`),
  UNIQUE KEY `idx_recharge_centers_code` (`code`),
  KEY `idx_recharge_centers_manager_id` (`manager_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `recharge_centers`
--

LOCK TABLES `recharge_centers` WRITE;
/*!40000 ALTER TABLE `recharge_centers` DISABLE KEYS */;
INSERT INTO `recharge_centers` VALUES ('center-bj-cy','北京朝阳中心','BJ_CY','','','','北京市朝阳区','','',9940,'active','2026-04-17 09:44:10.275','2026-04-17 16:34:15.805'),('center-bj-hd','北京海淀中心','BJ_HD','','','','北京市海淀区','','',0,'active','2026-04-17 09:44:10.277','2026-04-17 09:44:10.277'),('center-sh-pd','上海浦东中心','SH_PD','','','','上海市浦东新区','','',0,'active','2026-04-17 09:44:10.280','2026-04-17 09:44:10.280');
/*!40000 ALTER TABLE `recharge_centers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `recharge_operators`
--

DROP TABLE IF EXISTS `recharge_operators`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `recharge_operators` (
  `id` varchar(64) NOT NULL,
  `name` longtext,
  `phone` varchar(32) DEFAULT NULL,
  `password` longtext,
  `center_id` varchar(64) DEFAULT NULL,
  `role` longtext,
  `status` varchar(32) DEFAULT 'active',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_recharge_operators_phone` (`phone`),
  KEY `idx_recharge_operators_center_id` (`center_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `recharge_operators`
--

LOCK TABLES `recharge_operators` WRITE;
/*!40000 ALTER TABLE `recharge_operators` DISABLE KEYS */;
INSERT INTO `recharge_operators` VALUES ('op-licn','李出纳','13800138002','$2a$10$vXPjbfC511sMp3zdk1uFzOfxRWmtsZXnNIX7buP4C9Aq6In5YhV5S','center-bj-cy','出纳','active','2026-04-17 09:44:10.284','2026-04-17 09:44:10.284'),('op-zhangcw','张财务','13800138001','$2a$10$vXPjbfC511sMp3zdk1uFzOfxRWmtsZXnNIX7buP4C9Aq6In5YhV5S','center-bj-cy','财务','active','2026-04-17 09:44:10.283','2026-04-17 09:44:10.283');
/*!40000 ALTER TABLE `recharge_operators` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `schema_migrations`
--

DROP TABLE IF EXISTS `schema_migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `schema_migrations` (
  `version` varchar(255) NOT NULL,
  `executed_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `schema_migrations`
--

LOCK TABLES `schema_migrations` WRITE;
/*!40000 ALTER TABLE `schema_migrations` DISABLE KEYS */;
INSERT INTO `schema_migrations` VALUES ('20260415_100000_store_card_redesign.sql','2026-04-17 01:44:17');
/*!40000 ALTER TABLE `schema_migrations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `store_cards`
--

DROP TABLE IF EXISTS `store_cards`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `store_cards` (
  `id` varchar(64) NOT NULL,
  `card_no` varchar(32) NOT NULL COMMENT '卡号 TJ00000001',
  `card_type` tinyint DEFAULT '1' COMMENT '1=实体卡,2=虚拟卡',
  `status` tinyint DEFAULT '1' COMMENT '1=已入库,2=已发放,3=已激活,4=已冻结,5=已过期,6=已作废',
  `balance` int DEFAULT '1000' COMMENT '余额（元），面值固定1000',
  `recharge_center_id` varchar(64) DEFAULT NULL COMMENT '划拨到的充值中心ID',
  `user_id` varchar(64) DEFAULT NULL COMMENT '绑定的用户ID',
  `batch_no` varchar(64) DEFAULT '' COMMENT '批次号',
  `issue_reason` varchar(64) DEFAULT '' COMMENT '发放原因:购买套餐包/推荐奖励/其他',
  `issued_at` datetime DEFAULT NULL COMMENT '发放时间',
  `activated_at` datetime DEFAULT NULL COMMENT '激活时间（首次核销）',
  `expired_at` datetime DEFAULT NULL COMMENT '过期时间（激活日+1年）',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `card_no` (`card_no`),
  KEY `idx_status` (`status`),
  KEY `idx_recharge_center_id` (`recharge_center_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='门店卡';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `store_cards`
--

LOCK TABLES `store_cards` WRITE;
/*!40000 ALTER TABLE `store_cards` DISABLE KEYS */;
INSERT INTO `store_cards` VALUES ('b309d38c-3a0b-11f1-a9a1-195356b58b96','TJ00000001',1,3,900,'center-bj-cy','','B20260417001','购买套餐包','2026-04-17 14:11:47','2026-04-17 16:50:08','2027-04-17 16:50:08','2026-04-17 11:15:34','2026-04-17 16:50:08'),('b30ae812-3a0b-11f1-a9a1-195356b58b96','TJ00000002',1,2,1000,'center-bj-cy','','B20260417001','购买套餐包','2026-04-17 16:43:14',NULL,NULL,'2026-04-17 11:15:34','2026-04-17 16:43:14'),('b30afc12-3a0b-11f1-a9a1-195356b58b96','TJ00000003',1,2,1000,'center-bj-cy','','B20260417001','购买套餐包','2026-04-17 16:49:35',NULL,NULL,'2026-04-17 11:15:34','2026-04-17 16:49:35'),('b30afd8e-3a0b-11f1-a9a1-195356b58b96','TJ00000004',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 16:42:57'),('b30afed8-3a0b-11f1-a9a1-195356b58b96','TJ00000005',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 16:42:57'),('b30b0022-3a0b-11f1-a9a1-195356b58b96','TJ00000006',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 16:42:57'),('b30b018a-3a0b-11f1-a9a1-195356b58b96','TJ00000007',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 16:42:57'),('b30b02ac-3a0b-11f1-a9a1-195356b58b96','TJ00000008',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 16:42:57'),('b30b0c48-3a0b-11f1-a9a1-195356b58b96','TJ00000009',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 16:42:57'),('b30b0d6a-3a0b-11f1-a9a1-195356b58b96','TJ00000010',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 16:42:57'),('b30b0ebe-3a0b-11f1-a9a1-195356b58b96','TJ00000011',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 16:42:57'),('b30b101c-3a0b-11f1-a9a1-195356b58b96','TJ00000012',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:16'),('b30b1134-3a0b-11f1-a9a1-195356b58b96','TJ00000013',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:16'),('b30b2ad4-3a0b-11f1-a9a1-195356b58b96','TJ00000014',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:16'),('b30b2c3c-3a0b-11f1-a9a1-195356b58b96','TJ00000015',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:16'),('b30b4596-3a0b-11f1-a9a1-195356b58b96','TJ00000016',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:16'),('b30b4726-3a0b-11f1-a9a1-195356b58b96','TJ00000017',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:16'),('b30b4852-3a0b-11f1-a9a1-195356b58b96','TJ00000018',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:16'),('b30b4974-3a0b-11f1-a9a1-195356b58b96','TJ00000019',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:16'),('b30b4a8c-3a0b-11f1-a9a1-195356b58b96','TJ00000020',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:16'),('b30b4e88-3a0b-11f1-a9a1-195356b58b96','TJ00000021',1,1,1000,'center-bj-cy',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:16'),('b30b4fd2-3a0b-11f1-a9a1-195356b58b96','TJ00000022',1,1,1000,'center-sh-pd',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:44'),('b30b5aae-3a0b-11f1-a9a1-195356b58b96','TJ00000023',1,1,1000,'center-sh-pd',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:44'),('b30b5c0c-3a0b-11f1-a9a1-195356b58b96','TJ00000024',1,1,1000,'center-sh-pd',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:44'),('b30b5d24-3a0b-11f1-a9a1-195356b58b96','TJ00000025',1,1,1000,'center-sh-pd',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:44'),('b30b5f04-3a0b-11f1-a9a1-195356b58b96','TJ00000026',1,1,1000,'center-sh-pd',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:44'),('b30b659e-3a0b-11f1-a9a1-195356b58b96','TJ00000027',1,1,1000,'center-sh-pd',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:44'),('b30b66c0-3a0b-11f1-a9a1-195356b58b96','TJ00000028',1,1,1000,'center-sh-pd',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:44'),('b30b67ce-3a0b-11f1-a9a1-195356b58b96','TJ00000029',1,1,1000,'center-sh-pd',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:44'),('b30b6aa8-3a0b-11f1-a9a1-195356b58b96','TJ00000030',1,1,1000,'center-sh-pd',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:44'),('b30b6bc0-3a0b-11f1-a9a1-195356b58b96','TJ00000031',1,1,1000,'center-sh-pd',NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-18 08:40:44'),('b30b6cce-3a0b-11f1-a9a1-195356b58b96','TJ00000032',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b6ddc-3a0b-11f1-a9a1-195356b58b96','TJ00000033',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b707a-3a0b-11f1-a9a1-195356b58b96','TJ00000034',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b719c-3a0b-11f1-a9a1-195356b58b96','TJ00000035',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b72a0-3a0b-11f1-a9a1-195356b58b96','TJ00000036',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b73ae-3a0b-11f1-a9a1-195356b58b96','TJ00000037',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b7502-3a0b-11f1-a9a1-195356b58b96','TJ00000038',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b7b4c-3a0b-11f1-a9a1-195356b58b96','TJ00000039',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b7c5a-3a0b-11f1-a9a1-195356b58b96','TJ00000040',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b7d68-3a0b-11f1-a9a1-195356b58b96','TJ00000041',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b7ec6-3a0b-11f1-a9a1-195356b58b96','TJ00000042',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b8114-3a0b-11f1-a9a1-195356b58b96','TJ00000043',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b8218-3a0b-11f1-a9a1-195356b58b96','TJ00000044',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b8330-3a0b-11f1-a9a1-195356b58b96','TJ00000045',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b93ca-3a0b-11f1-a9a1-195356b58b96','TJ00000046',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b9550-3a0b-11f1-a9a1-195356b58b96','TJ00000047',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b965e-3a0b-11f1-a9a1-195356b58b96','TJ00000048',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b97e4-3a0b-11f1-a9a1-195356b58b96','TJ00000049',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30b9cda-3a0b-11f1-a9a1-195356b58b96','TJ00000050',1,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30bd682-3a0b-11f1-a9a1-195356b58b96','TJ00000051',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30bd970-3a0b-11f1-a9a1-195356b58b96','TJ00000052',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30bda7e-3a0b-11f1-a9a1-195356b58b96','TJ00000053',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30be33e-3a0b-11f1-a9a1-195356b58b96','TJ00000054',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30be424-3a0b-11f1-a9a1-195356b58b96','TJ00000055',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30be50a-3a0b-11f1-a9a1-195356b58b96','TJ00000056',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30be5f0-3a0b-11f1-a9a1-195356b58b96','TJ00000057',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30be6cc-3a0b-11f1-a9a1-195356b58b96','TJ00000058',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30bffa4-3a0b-11f1-a9a1-195356b58b96','TJ00000059',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c009e-3a0b-11f1-a9a1-195356b58b96','TJ00000060',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c0184-3a0b-11f1-a9a1-195356b58b96','TJ00000061',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c0260-3a0b-11f1-a9a1-195356b58b96','TJ00000062',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c033c-3a0b-11f1-a9a1-195356b58b96','TJ00000063',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c040e-3a0b-11f1-a9a1-195356b58b96','TJ00000064',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c04ea-3a0b-11f1-a9a1-195356b58b96','TJ00000065',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c05c6-3a0b-11f1-a9a1-195356b58b96','TJ00000066',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c06a2-3a0b-11f1-a9a1-195356b58b96','TJ00000067',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c0774-3a0b-11f1-a9a1-195356b58b96','TJ00000068',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c0850-3a0b-11f1-a9a1-195356b58b96','TJ00000069',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c092c-3a0b-11f1-a9a1-195356b58b96','TJ00000070',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c0a08-3a0b-11f1-a9a1-195356b58b96','TJ00000071',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c0ada-3a0b-11f1-a9a1-195356b58b96','TJ00000072',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c0bb6-3a0b-11f1-a9a1-195356b58b96','TJ00000073',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c0c92-3a0b-11f1-a9a1-195356b58b96','TJ00000074',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c0d64-3a0b-11f1-a9a1-195356b58b96','TJ00000075',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c0e40-3a0b-11f1-a9a1-195356b58b96','TJ00000076',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c0f1c-3a0b-11f1-a9a1-195356b58b96','TJ00000077',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c0ff8-3a0b-11f1-a9a1-195356b58b96','TJ00000078',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c10d4-3a0b-11f1-a9a1-195356b58b96','TJ00000079',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c11a6-3a0b-11f1-a9a1-195356b58b96','TJ00000080',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c12c8-3a0b-11f1-a9a1-195356b58b96','TJ00000081',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c13a4-3a0b-11f1-a9a1-195356b58b96','TJ00000082',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c1480-3a0b-11f1-a9a1-195356b58b96','TJ00000083',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c1552-3a0b-11f1-a9a1-195356b58b96','TJ00000084',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c1638-3a0b-11f1-a9a1-195356b58b96','TJ00000085',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c1714-3a0b-11f1-a9a1-195356b58b96','TJ00000086',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c17f0-3a0b-11f1-a9a1-195356b58b96','TJ00000087',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c18cc-3a0b-11f1-a9a1-195356b58b96','TJ00000088',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c19b2-3a0b-11f1-a9a1-195356b58b96','TJ00000089',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c1a8e-3a0b-11f1-a9a1-195356b58b96','TJ00000090',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c1b60-3a0b-11f1-a9a1-195356b58b96','TJ00000091',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c1c3c-3a0b-11f1-a9a1-195356b58b96','TJ00000092',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c1d18-3a0b-11f1-a9a1-195356b58b96','TJ00000093',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c1dea-3a0b-11f1-a9a1-195356b58b96','TJ00000094',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c1ec6-3a0b-11f1-a9a1-195356b58b96','TJ00000095',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c1f98-3a0b-11f1-a9a1-195356b58b96','TJ00000096',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c2074-3a0b-11f1-a9a1-195356b58b96','TJ00000097',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c21a0-3a0b-11f1-a9a1-195356b58b96','TJ00000098',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c2290-3a0b-11f1-a9a1-195356b58b96','TJ00000099',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34'),('b30c2452-3a0b-11f1-a9a1-195356b58b96','TJ00000100',2,1,1000,NULL,NULL,'B20260417001','',NULL,NULL,NULL,'2026-04-17 11:15:34','2026-04-17 11:15:34');
/*!40000 ALTER TABLE `store_cards` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用户名（登录账号）',
  `phone` varchar(20) NOT NULL COMMENT '手机号',
  `password` varchar(255) NOT NULL COMMENT '密码（bcrypt加密）',
  `name` varchar(50) NOT NULL COMMENT '姓名',
  `role` varchar(20) DEFAULT 'operator' COMMENT '角色：super_admin=超管, admin=管理员, operator=操作员',
  `center_id` varchar(64) DEFAULT NULL COMMENT '所属充值中心ID',
  `center_name` varchar(100) DEFAULT NULL COMMENT '所属充值中心名称（冗余字段）',
  `status` tinyint DEFAULT '1' COMMENT '状态：1=启用, 0=禁用',
  `last_login_at` datetime(3) DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_phone` (`phone`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'','13800000000','$2a$10$vXPjbfC511sMp3zdk1uFzOfxRWmtsZXnNIX7buP4C9Aq6In5YhV5S','系统管理员','hq_admin',NULL,NULL,1,'2026-04-18 11:31:58.000','','2026-04-17 09:44:10.271','2026-04-17 09:44:10.271',NULL),(2,'13900000001','13900000001','$2a$10$zF55Gwd8GZkK/c9MFWAlsetJeRFLPka4EEq0cO1XUX7Qv/CaDOJF2','13900000001-1','super_admin',NULL,NULL,1,'2026-04-18 11:03:49.000','','2026-04-17 10:45:37.294','2026-04-18 09:59:48.739',NULL),(3,'test_center_admin','13700001111','$2a$10$79jb4ByYtlKpHaCZ1xhyCuq1zl2zCWXrdUIBCq81OG.rgvr3XNni2','测试中心管理员','center_admin','center-bj-cy',NULL,1,NULL,NULL,'2026-04-18 10:04:54.580','2026-04-18 10:04:54.580','2026-04-18 10:05:13.000'),(4,'18613802422','18613802422','$2a$10$M7wz1e1YMBd9z.6iusp9N.Y.EdbAfXGKOMCdxZTfLwsMRuHhOqFki','b端充值','center_admin','center-bj-cy',NULL,1,'2026-04-18 10:53:45.000','','2026-04-18 10:11:47.428','2026-04-18 10:23:40.583',NULL),(5,'18613802423','18613802423','$2a$10$puaVYPDNnGmzXf8tb.4Yv.XlgTpoNCP2IpdaoG7HWl5NMCoNphxcC','C端-操作员','operator','center-bj-cy',NULL,1,'2026-04-18 11:47:50.000','','2026-04-18 11:34:17.302','2026-04-18 11:34:17.302',NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'sale_dev'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-04-18 14:55:57
