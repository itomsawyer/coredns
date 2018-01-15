-- MySQL dump 10.13  Distrib 5.7.17, for macos10.12 (x86_64)
--
-- Host: localhost    Database: iwg
-- ------------------------------------------------------
-- Server version	5.6.35

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `_locker`
--

CREATE SCHEMA IF NOT EXISTS `iwg` DEFAULT CHARACTER SET utf8 ;
USE `iwg` ;

DROP TABLE IF EXISTS `_locker`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `_locker` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL DEFAULT '',
  `clientset_id` int(11) DEFAULT NULL,
  `domain_pool_id` int(11) DEFAULT NULL,
  `netlink_id` int(11) DEFAULT NULL,
  `viewer_id` int(11) DEFAULT NULL,
  `routeset_id` int(11) DEFAULT NULL,
  `domainlink_id` int(11) DEFAULT NULL,
  `netlinkset_id` int(11) DEFAULT NULL,
  `route_id` int(11) DEFAULT NULL,
  `outlink_id` int(11) DEFAULT NULL,
  `policy_id` int(11) DEFAULT NULL,
  `policy_detail_id` int(11) DEFAULT NULL,
  `ldns_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `clientset_id` (`clientset_id`),
  KEY `domain_pool_id` (`domain_pool_id`),
  KEY `netlink_id` (`netlink_id`),
  KEY `viewer_id` (`viewer_id`),
  KEY `routeset_id` (`routeset_id`),
  KEY `domainlink_id` (`domainlink_id`),
  KEY `netlinkset_id` (`netlinkset_id`),
  KEY `route_id` (`route_id`),
  KEY `outlink_id` (`outlink_id`),
  KEY `policy_id` (`policy_id`),
  KEY `policy_detail_id` (`policy_detail_id`),
  KEY `ldns_id` (`ldns_id`),
  CONSTRAINT `_locker_ibfk_1` FOREIGN KEY (`clientset_id`) REFERENCES `clientset` (`id`),
  CONSTRAINT `_locker_ibfk_10` FOREIGN KEY (`policy_id`) REFERENCES `policy` (`id`),
  CONSTRAINT `_locker_ibfk_11` FOREIGN KEY (`policy_detail_id`) REFERENCES `policy_detail` (`id`),
  CONSTRAINT `_locker_ibfk_12` FOREIGN KEY (`ldns_id`) REFERENCES `ldns` (`id`),
  CONSTRAINT `_locker_ibfk_2` FOREIGN KEY (`domain_pool_id`) REFERENCES `domain_pool` (`id`),
  CONSTRAINT `_locker_ibfk_3` FOREIGN KEY (`netlink_id`) REFERENCES `netlink` (`id`),
  CONSTRAINT `_locker_ibfk_4` FOREIGN KEY (`viewer_id`) REFERENCES `viewer` (`id`),
  CONSTRAINT `_locker_ibfk_5` FOREIGN KEY (`routeset_id`) REFERENCES `routeset` (`id`),
  CONSTRAINT `_locker_ibfk_6` FOREIGN KEY (`domainlink_id`) REFERENCES `domainlink` (`id`),
  CONSTRAINT `_locker_ibfk_7` FOREIGN KEY (`netlinkset_id`) REFERENCES `netlinkset` (`id`),
  CONSTRAINT `_locker_ibfk_8` FOREIGN KEY (`route_id`) REFERENCES `route` (`id`),
  CONSTRAINT `_locker_ibfk_9` FOREIGN KEY (`outlink_id`) REFERENCES `outlink` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='Internal lock for default data, DO NOT MODIFY unless you known what you are doing';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `_locker`
--

LOCK TABLES `_locker` WRITE;
/*!40000 ALTER TABLE `_locker` DISABLE KEYS */;
INSERT INTO `_locker` VALUES (1,'global default data locker',1,1,1,1,1,1,1,1,1,1,1,1);
/*!40000 ALTER TABLE `_locker` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `base_route_view`
--

DROP TABLE IF EXISTS `base_route_view`;
/*!50001 DROP VIEW IF EXISTS `base_route_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `base_route_view` AS SELECT 
 1 AS `routeset_id`,
 1 AS `routeset_name`,
 1 AS `netlinkset_id`,
 1 AS `netlinkset_name`,
 1 AS `route_id`,
 1 AS `route_priority`,
 1 AS `route_score`,
 1 AS `outlink_id`,
 1 AS `outlink_name`,
 1 AS `outlink_addr`,
 1 AS `outlink_typ`,
 1 AS `outlink_enable`,
 1 AS `outlink_unavailable`,
 1 AS `route_enable`,
 1 AS `route_unavailable`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `best_route_view`
--

DROP TABLE IF EXISTS `best_route_view`;
/*!50001 DROP VIEW IF EXISTS `best_route_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `best_route_view` AS SELECT 
 1 AS `routeset_id`,
 1 AS `netlinkset_id`,
 1 AS `route_priority`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `clientset`
--

DROP TABLE IF EXISTS `clientset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `clientset` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(127) NOT NULL,
  `info` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='client ipnet set';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `clientset`
--

LOCK TABLES `clientset` WRITE;
/*!40000 ALTER TABLE `clientset` DISABLE KEYS */;
INSERT INTO `clientset` VALUES (1,'default','global default');
/*!40000 ALTER TABLE `clientset` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `clientset_view`
--

DROP TABLE IF EXISTS `clientset_view`;
/*!50001 DROP VIEW IF EXISTS `clientset_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `clientset_view` AS SELECT 
 1 AS `ipnet_id`,
 1 AS `ip_start`,
 1 AS `ip_start_int`,
 1 AS `ip_end`,
 1 AS `ip_end_int`,
 1 AS `ipnet`,
 1 AS `mask`,
 1 AS `clientset_id`,
 1 AS `clientset_name`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `clientset_wl_view`
--

DROP TABLE IF EXISTS `clientset_wl_view`;
/*!50001 DROP VIEW IF EXISTS `clientset_wl_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `clientset_wl_view` AS SELECT 
 1 AS `ipnet_wl_id`,
 1 AS `ip_start`,
 1 AS `ip_start_int`,
 1 AS `ip_end`,
 1 AS `ip_end_int`,
 1 AS `ipnet`,
 1 AS `mask`,
 1 AS `clientset_id`,
 1 AS `clientset_name`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `dns_forward_view`
--

DROP TABLE IF EXISTS `dns_forward_view`;
/*!50001 DROP VIEW IF EXISTS `dns_forward_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `dns_forward_view` AS SELECT 
 1 AS `zone_id`,
 1 AS `dm`,
 1 AS `forward_typ`,
 1 AS `ldns_name`,
 1 AS `ldns_addr`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `dns_forward_zone`
--

DROP TABLE IF EXISTS `dns_forward_zone`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dns_forward_zone` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `dm` varchar(255) DEFAULT NULL,
  `typ` varchar(16) NOT NULL DEFAULT 'only',
  PRIMARY KEY (`id`),
  UNIQUE KEY `dm_UNIQUE` (`dm`)
) ENGINE=InnoDB AUTO_INCREMENT=220 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dns_forward_zone`
--

LOCK TABLES `dns_forward_zone` WRITE;
/*!40000 ALTER TABLE `dns_forward_zone` DISABLE KEYS */;
/*!40000 ALTER TABLE `dns_forward_zone` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dns_forwarders`
--

DROP TABLE IF EXISTS `dns_forwarders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dns_forwarders` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `zone_id` int(11) DEFAULT NULL,
  `ldns_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `zone_id_idx` (`zone_id`),
  KEY `ldns_id_idx` (`ldns_id`),
  CONSTRAINT `ldns_id` FOREIGN KEY (`ldns_id`) REFERENCES `ldns` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `zone_id` FOREIGN KEY (`zone_id`) REFERENCES `dns_forward_zone` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=284 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dns_forwarders`
--

LOCK TABLES `dns_forwarders` WRITE;
/*!40000 ALTER TABLE `dns_forwarders` DISABLE KEYS */;
/*!40000 ALTER TABLE `dns_forwarders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `domain`
--

DROP TABLE IF EXISTS `domain`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `domain` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `domain` varchar(255) NOT NULL,
  `domain_pool_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `domain` (`domain`),
  KEY `domain_pool_id` (`domain_pool_id`),
  CONSTRAINT `domain_ibfk_1` FOREIGN KEY (`domain_pool_id`) REFERENCES `domain_pool` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='serve domains';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `domain`
--

LOCK TABLES `domain` WRITE;
/*!40000 ALTER TABLE `domain` DISABLE KEYS */;
/*!40000 ALTER TABLE `domain` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `domain_pool`
--

DROP TABLE IF EXISTS `domain_pool`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `domain_pool` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(127) NOT NULL,
  `info` varchar(255) NOT NULL DEFAULT '',
  `typ` varchar(16) NOT NULL DEFAULT 'normal',
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `unavailable` smallint(5) unsigned NOT NULL DEFAULT '0',
  `domain_monitor` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='serve domains set';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `domain_pool`
--

LOCK TABLES `domain_pool` WRITE;
/*!40000 ALTER TABLE `domain_pool` DISABLE KEYS */;
INSERT INTO `domain_pool` VALUES (1,'default','global default','normal',1,0,0);
/*!40000 ALTER TABLE `domain_pool` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `domain_view`
--

DROP TABLE IF EXISTS `domain_view`;
/*!50001 DROP VIEW IF EXISTS `domain_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `domain_view` AS SELECT 
 1 AS `domain_id`,
 1 AS `domain`,
 1 AS `domain_pool_id`,
 1 AS `pool_name`,
 1 AS `domain_monitor`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `domainlink`
--

DROP TABLE IF EXISTS `domainlink`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `domainlink` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `domain_pool_id` int(11) NOT NULL,
  `netlink_id` int(11) NOT NULL,
  `netlinkset_id` int(11) NOT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `domain_pool_id` (`domain_pool_id`,`netlink_id`),
  KEY `netlink_id` (`netlink_id`),
  KEY `netlinkset_id` (`netlinkset_id`),
  CONSTRAINT `domainlink_ibfk_1` FOREIGN KEY (`domain_pool_id`) REFERENCES `domain_pool` (`id`),
  CONSTRAINT `domainlink_ibfk_2` FOREIGN KEY (`netlink_id`) REFERENCES `netlink` (`id`),
  CONSTRAINT `domainlink_ibfk_3` FOREIGN KEY (`netlinkset_id`) REFERENCES `netlinkset` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='bind domain_pool and netlink to a netlinkset';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `domainlink`
--

LOCK TABLES `domainlink` WRITE;
/*!40000 ALTER TABLE `domainlink` DISABLE KEYS */;
INSERT INTO `domainlink` VALUES (1,1,1,1,1);
/*!40000 ALTER TABLE `domainlink` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `dst_view`
--

DROP TABLE IF EXISTS `dst_view`;
/*!50001 DROP VIEW IF EXISTS `dst_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `dst_view` AS SELECT 
 1 AS `domain_pool_id`,
 1 AS `netlink_id`,
 1 AS `netlinkset_id`,
 1 AS `domain_pool_name`,
 1 AS `isp`,
 1 AS `region`,
 1 AS `netlinkset_name`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `filter`
--

DROP TABLE IF EXISTS `filter`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `filter` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `src_ip_start` varchar(40) DEFAULT NULL,
  `src_ip_end` varchar(40) DEFAULT NULL,
  `clientset_id` int(11) DEFAULT NULL,
  `domain_id` int(11) DEFAULT NULL,
  `domain_pool_id` int(11) DEFAULT NULL,
  `dst_ip_start` varchar(40) DEFAULT NULL,
  `dst_ip_end` varchar(40) DEFAULT NULL,
  `netlink_id` int(11) DEFAULT NULL,
  `target_ip` varchar(40) DEFAULT NULL,
  `outlink_id` int(11) DEFAULT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `clientset_id` (`clientset_id`),
  KEY `netlink_id` (`netlink_id`),
  KEY `domain_id` (`domain_id`),
  KEY `domain_pool_id` (`domain_pool_id`),
  KEY `outlink_id` (`outlink_id`),
  CONSTRAINT `filter_ibfk_1` FOREIGN KEY (`clientset_id`) REFERENCES `clientset` (`id`) ON DELETE CASCADE,
  CONSTRAINT `filter_ibfk_2` FOREIGN KEY (`netlink_id`) REFERENCES `netlink` (`id`) ON DELETE CASCADE,
  CONSTRAINT `filter_ibfk_3` FOREIGN KEY (`domain_id`) REFERENCES `domain` (`id`) ON DELETE CASCADE,
  CONSTRAINT `filter_ibfk_4` FOREIGN KEY (`domain_pool_id`) REFERENCES `domain_pool` (`id`) ON DELETE CASCADE,
  CONSTRAINT `filter_ibfk_5` FOREIGN KEY (`outlink_id`) REFERENCES `outlink` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='custom route strategy like iptables';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `filter`
--

LOCK TABLES `filter` WRITE;
/*!40000 ALTER TABLE `filter` DISABLE KEYS */;
/*!40000 ALTER TABLE `filter` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `filter_view`
--

DROP TABLE IF EXISTS `filter_view`;
/*!50001 DROP VIEW IF EXISTS `filter_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `filter_view` AS SELECT 
 1 AS `id`,
 1 AS `src_ip_start`,
 1 AS `src_ip_start_int`,
 1 AS `src_ip_end`,
 1 AS `src_ip_end_int`,
 1 AS `clientset_id`,
 1 AS `clientset_name`,
 1 AS `domain_id`,
 1 AS `domain`,
 1 AS `domain_pool_id`,
 1 AS `domain_pool_name`,
 1 AS `dst_ip_start`,
 1 AS `dst_ip_start_int`,
 1 AS `dst_ip_end`,
 1 AS `dst_ip_end_int`,
 1 AS `netlink_id`,
 1 AS `netlink_name`,
 1 AS `target_ip`,
 1 AS `target_ip_int`,
 1 AS `outlink_id`,
 1 AS `outlink_name`,
 1 AS `outlink_addr`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `ipnet`
--

DROP TABLE IF EXISTS `ipnet`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ipnet` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip_start` varchar(40) NOT NULL,
  `ip_end` varchar(40) NOT NULL,
  `ipnet` varchar(40) NOT NULL,
  `mask` tinyint(3) unsigned NOT NULL,
  `clientset_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_start` (`ip_start`,`ip_end`),
  UNIQUE KEY `ipnet` (`ipnet`,`mask`),
  KEY `clientset_id` (`clientset_id`),
  CONSTRAINT `ipnet_ibfk_1` FOREIGN KEY (`clientset_id`) REFERENCES `clientset` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='client ipnet';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ipnet`
--

LOCK TABLES `ipnet` WRITE;
/*!40000 ALTER TABLE `ipnet` DISABLE KEYS */;
/*!40000 ALTER TABLE `ipnet` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ipnet_wl`
--

DROP TABLE IF EXISTS `ipnet_wl`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ipnet_wl` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip_start` varchar(40) NOT NULL,
  `ip_end` varchar(40) NOT NULL,
  `ipnet` varchar(40) NOT NULL,
  `mask` tinyint(3) unsigned NOT NULL,
  `clientset_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_start` (`ip_start`,`ip_end`),
  UNIQUE KEY `ipnet` (`ipnet`,`mask`),
  KEY `clientset_id` (`clientset_id`),
  CONSTRAINT `ipnet_wl_ibfk_1` FOREIGN KEY (`clientset_id`) REFERENCES `clientset` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='client ipnet whitelist';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ipnet_wl`
--

LOCK TABLES `ipnet_wl` WRITE;
/*!40000 ALTER TABLE `ipnet_wl` DISABLE KEYS */;
/*!40000 ALTER TABLE `ipnet_wl` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `iptable`
--

DROP TABLE IF EXISTS `iptable`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `iptable` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip_start` varchar(40) NOT NULL,
  `ip_end` varchar(40) NOT NULL,
  `ipnet` varchar(40) NOT NULL,
  `mask` tinyint(3) unsigned NOT NULL,
  `netlink_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_start` (`ip_start`,`ip_end`),
  UNIQUE KEY `ipnet` (`ipnet`,`mask`),
  KEY `netlink_id` (`netlink_id`),
  CONSTRAINT `iptable_ibfk_1` FOREIGN KEY (`netlink_id`) REFERENCES `netlink` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='IP to netlink';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `iptable`
--

LOCK TABLES `iptable` WRITE;
/*!40000 ALTER TABLE `iptable` DISABLE KEYS */;
/*!40000 ALTER TABLE `iptable` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `iptable_wl`
--

DROP TABLE IF EXISTS `iptable_wl`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `iptable_wl` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip_start` varchar(40) NOT NULL,
  `ip_end` varchar(40) NOT NULL,
  `ipnet` varchar(40) NOT NULL,
  `mask` tinyint(3) unsigned NOT NULL,
  `netlink_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_start` (`ip_start`,`ip_end`),
  UNIQUE KEY `ipnet` (`ipnet`,`mask`),
  KEY `netlink_id` (`netlink_id`),
  CONSTRAINT `iptable_wl_ibfk_1` FOREIGN KEY (`netlink_id`) REFERENCES `netlink` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='IP to netlink whitelist';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `iptable_wl`
--

LOCK TABLES `iptable_wl` WRITE;
/*!40000 ALTER TABLE `iptable_wl` DISABLE KEYS */;
/*!40000 ALTER TABLE `iptable_wl` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ldns`
--

DROP TABLE IF EXISTS `ldns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ldns` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(127) NOT NULL,
  `addr` varchar(40) NOT NULL,
  `typ` varchar(32) NOT NULL DEFAULT 'upstream',
  `checkdm` varchar(64) NOT NULL DEFAULT 'baidu.com',
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `unavailable` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT 'if other than zero, ldns is unavailable with each bit indicate different reason',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='upstream ldns info';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ldns`
--

LOCK TABLES `ldns` WRITE;
/*!40000 ALTER TABLE `ldns` DISABLE KEYS */;
INSERT INTO `ldns` VALUES (1,'default','223.5.5.5','upstream','baidu.com',1,0);
INSERT INTO `ldns` VALUES (2,'ali_public_dns_2','223.6.6.6','upstream','baidu.com',1,0);
/*!40000 ALTER TABLE `ldns` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `natlink`
--

DROP TABLE IF EXISTS `natlink`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `natlink` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `outlink_id` int(11) DEFAULT NULL,
  `natserver_id` int(11) DEFAULT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `gw` varchar(126) NOT NULL,
  `addr` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `outlink_id` (`outlink_id`),
  KEY `natserver_id` (`natserver_id`),
  CONSTRAINT `natlink_ibfk_1` FOREIGN KEY (`outlink_id`) REFERENCES `outlink` (`id`),
  CONSTRAINT `natlink_ibfk_2` FOREIGN KEY (`natserver_id`) REFERENCES `natserver` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='outlink on specific nat server';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `natlink`
--

LOCK TABLES `natlink` WRITE;
/*!40000 ALTER TABLE `natlink` DISABLE KEYS */;
/*!40000 ALTER TABLE `natlink` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `natserver`
--

DROP TABLE IF EXISTS `natserver`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `natserver` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(127) NOT NULL,
  `addr` varchar(127) NOT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `unavailable` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT 'if other than zero, outlink is unavailable, each bit indicate different reason',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='nat server';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `natserver`
--

LOCK TABLES `natserver` WRITE;
/*!40000 ALTER TABLE `natserver` DISABLE KEYS */;
/*!40000 ALTER TABLE `natserver` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `netlink`
--

DROP TABLE IF EXISTS `netlink`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `netlink` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `isp` varchar(127) NOT NULL,
  `region` varchar(127) NOT NULL DEFAULT '',
  `typ` varchar(32) NOT NULL DEFAULT 'normal',
  PRIMARY KEY (`id`),
  UNIQUE KEY `isp` (`isp`,`region`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='netlink (isp + province or CP) of a target ip';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `netlink`
--

LOCK TABLES `netlink` WRITE;
/*!40000 ALTER TABLE `netlink` DISABLE KEYS */;
INSERT INTO `netlink` VALUES (1,'default','global default','normal');
/*!40000 ALTER TABLE `netlink` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `netlink_view`
--

DROP TABLE IF EXISTS `netlink_view`;
/*!50001 DROP VIEW IF EXISTS `netlink_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `netlink_view` AS SELECT 
 1 AS `iptable_id`,
 1 AS `ip_start`,
 1 AS `ip_start_int`,
 1 AS `ip_end`,
 1 AS `ip_end_int`,
 1 AS `ipnet`,
 1 AS `mask`,
 1 AS `netlink_id`,
 1 AS `isp`,
 1 AS `region`,
 1 AS `typ`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `netlink_wl_view`
--

DROP TABLE IF EXISTS `netlink_wl_view`;
/*!50001 DROP VIEW IF EXISTS `netlink_wl_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `netlink_wl_view` AS SELECT 
 1 AS `iptable_wl_id`,
 1 AS `ip_start`,
 1 AS `ip_start_int`,
 1 AS `ip_end`,
 1 AS `ip_end_int`,
 1 AS `ipnet`,
 1 AS `mask`,
 1 AS `netlink_id`,
 1 AS `isp`,
 1 AS `region`,
 1 AS `typ`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `netlinkset`
--

DROP TABLE IF EXISTS `netlinkset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `netlinkset` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(127) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='netlink set ';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `netlinkset`
--

LOCK TABLES `netlinkset` WRITE;
/*!40000 ALTER TABLE `netlinkset` DISABLE KEYS */;
INSERT INTO `netlinkset` VALUES (1,'default');
/*!40000 ALTER TABLE `netlinkset` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `oplog`
--

DROP TABLE IF EXISTS `oplog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `oplog` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `opr` varchar(16) DEFAULT NULL,
  `action` varchar(6) DEFAULT NULL,
  `tbl` varchar(16) DEFAULT NULL,
  `row_id` int(11) DEFAULT NULL,
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Log of operations to other tables in this db';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `oplog`
--

LOCK TABLES `oplog` WRITE;
/*!40000 ALTER TABLE `oplog` DISABLE KEYS */;
/*!40000 ALTER TABLE `oplog` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `outlink`
--

DROP TABLE IF EXISTS `outlink`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `outlink` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(127) NOT NULL,
  `addr` varchar(255) NOT NULL,
  `typ` varchar(32) NOT NULL DEFAULT 'normal',
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `unavailable` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT 'if other than zero, outlink is unavailable, each bit indicate different reason',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='network gateway, aka outlink';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `outlink`
--

LOCK TABLES `outlink` WRITE;
/*!40000 ALTER TABLE `outlink` DISABLE KEYS */;
INSERT INTO `outlink` VALUES (1,'default','0.0.0.0','normal',1,0);
/*!40000 ALTER TABLE `outlink` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `outlink_view`
--

DROP TABLE IF EXISTS `outlink_view`;
/*!50001 DROP VIEW IF EXISTS `outlink_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `outlink_view` AS SELECT 
 1 AS `outlink_id`,
 1 AS `outlink_addr`,
 1 AS `natlink_addr`,
 1 AS `natserver_id`,
 1 AS `nat_name`,
 1 AS `natlink_gw`,
 1 AS `natlink_status`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `policy`
--

DROP TABLE IF EXISTS `policy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `policy` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(127) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='policy index of choose ldns upstream forwarder';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `policy`
--

LOCK TABLES `policy` WRITE;
/*!40000 ALTER TABLE `policy` DISABLE KEYS */;
INSERT INTO `policy` VALUES (1,'default');
/*!40000 ALTER TABLE `policy` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `policy_detail`
--

DROP TABLE IF EXISTS `policy_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `policy_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `policy_id` int(11) NOT NULL,
  `policy_sequence` int(11) NOT NULL DEFAULT '0',
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `priority` smallint(6) NOT NULL DEFAULT '20',
  `weight` smallint(6) NOT NULL DEFAULT '100',
  `op` varchar(127) NOT NULL DEFAULT 'and',
  `op_typ` varchar(32) NOT NULL DEFAULT 'builtin',
  `ldns_id` int(11) NOT NULL,
  `rrset_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `policy_id` (`policy_id`,`ldns_id`),
  KEY `ldns_id` (`ldns_id`),
  KEY `rrset_id` (`rrset_id`),
  CONSTRAINT `policy_detail_ibfk_1` FOREIGN KEY (`ldns_id`) REFERENCES `ldns` (`id`),
  CONSTRAINT `policy_detail_ibfk_2` FOREIGN KEY (`rrset_id`) REFERENCES `rrset` (`id`),
  CONSTRAINT `policy_detail_ibfk_3` FOREIGN KEY (`policy_id`) REFERENCES `policy` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='policy detail of choose ldns upstream forwarder';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `policy_detail`
--

LOCK TABLES `policy_detail` WRITE;
/*!40000 ALTER TABLE `policy_detail` DISABLE KEYS */;
INSERT INTO `policy_detail` VALUES (1,1,0,1,20,100,'and','builtin',1,NULL);
/*!40000 ALTER TABLE `policy_detail` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `policy_view`
--

DROP TABLE IF EXISTS `policy_view`;
/*!50001 DROP VIEW IF EXISTS `policy_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `policy_view` AS SELECT 
 1 AS `policy_id`,
 1 AS `policy_name`,
 1 AS `policy_sequence`,
 1 AS `priority`,
 1 AS `weight`,
 1 AS `op`,
 1 AS `op_typ`,
 1 AS `ldns_id`,
 1 AS `name`,
 1 AS `addr`,
 1 AS `typ`,
 1 AS `rrset_id`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `route`
--

DROP TABLE IF EXISTS `route`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `route` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `routeset_id` int(11) NOT NULL,
  `netlinkset_id` int(11) NOT NULL,
  `outlink_id` int(11) NOT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `priority` smallint(6) NOT NULL DEFAULT '20',
  `score` smallint(6) NOT NULL DEFAULT '50' COMMENT 'netlink performance index',
  `unavailable` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT 'if other than zero, route is unavailable, each bit indicate different reason',
  PRIMARY KEY (`id`),
  UNIQUE KEY `netlinkset_id` (`netlinkset_id`,`outlink_id`,`routeset_id`),
  UNIQUE KEY `netlinkset_id_2` (`netlinkset_id`,`outlink_id`,`priority`),
  KEY `outlink_id` (`outlink_id`),
  KEY `routeset_id` (`routeset_id`),
  CONSTRAINT `route_ibfk_1` FOREIGN KEY (`outlink_id`) REFERENCES `outlink` (`id`),
  CONSTRAINT `route_ibfk_2` FOREIGN KEY (`netlinkset_id`) REFERENCES `netlinkset` (`id`),
  CONSTRAINT `route_ibfk_3` FOREIGN KEY (`routeset_id`) REFERENCES `routeset` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='Performance of using gateway to serve paricular netlink';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `route`
--

LOCK TABLES `route` WRITE;
/*!40000 ALTER TABLE `route` DISABLE KEYS */;
INSERT INTO `route` VALUES (1,1,1,1,1,20,50,0);
/*!40000 ALTER TABLE `route` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `route_view`
--

DROP TABLE IF EXISTS `route_view`;
/*!50001 DROP VIEW IF EXISTS `route_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `route_view` AS SELECT 
 1 AS `routeset_id`,
 1 AS `routeset_name`,
 1 AS `netlinkset_id`,
 1 AS `netlinkset_name`,
 1 AS `route_id`,
 1 AS `route_priority`,
 1 AS `route_score`,
 1 AS `outlink_id`,
 1 AS `outlink_name`,
 1 AS `outlink_addr`,
 1 AS `outlink_typ`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `routeset`
--

DROP TABLE IF EXISTS `routeset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `routeset` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(127) NOT NULL,
  `info` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `routeset`
--

LOCK TABLES `routeset` WRITE;
/*!40000 ALTER TABLE `routeset` DISABLE KEYS */;
INSERT INTO `routeset` VALUES (1,'default','global default');
/*!40000 ALTER TABLE `routeset` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rr`
--

DROP TABLE IF EXISTS `rr`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rr` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `rrtype` int(11) NOT NULL,
  `rrdata` varchar(255) NOT NULL,
  `ttl` int(10) unsigned NOT NULL DEFAULT '300',
  `rrset_id` int(11) DEFAULT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `rrset_id` (`rrset_id`),
  CONSTRAINT `rr_ibfk_1` FOREIGN KEY (`rrset_id`) REFERENCES `rrset` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='dns rr(resource record)';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rr`
--

LOCK TABLES `rr` WRITE;
/*!40000 ALTER TABLE `rr` DISABLE KEYS */;
/*!40000 ALTER TABLE `rr` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rrset`
--

DROP TABLE IF EXISTS `rrset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rrset` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '',
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='dns rrset(resource record)';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rrset`
--

LOCK TABLES `rrset` WRITE;
/*!40000 ALTER TABLE `rrset` DISABLE KEYS */;
/*!40000 ALTER TABLE `rrset` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `src_view`
--

DROP TABLE IF EXISTS `src_view`;
/*!50001 DROP VIEW IF EXISTS `src_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `src_view` AS SELECT 
 1 AS `clientset_id`,
 1 AS `clientset_name`,
 1 AS `domain_pool_id`,
 1 AS `domain_pool_name`,
 1 AS `routeset_id`,
 1 AS `routeset_name`,
 1 AS `policy_id`,
 1 AS `policy_name`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `viewer`
--

DROP TABLE IF EXISTS `viewer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `viewer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `clientset_id` int(11) NOT NULL,
  `domain_pool_id` int(11) NOT NULL,
  `routeset_id` int(11) NOT NULL,
  `policy_id` int(11) NOT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `clientset_id` (`clientset_id`,`domain_pool_id`),
  KEY `domain_pool_id` (`domain_pool_id`),
  KEY `routeset_id` (`routeset_id`),
  KEY `policy_id` (`policy_id`),
  CONSTRAINT `viewer_ibfk_1` FOREIGN KEY (`clientset_id`) REFERENCES `clientset` (`id`),
  CONSTRAINT `viewer_ibfk_2` FOREIGN KEY (`domain_pool_id`) REFERENCES `domain_pool` (`id`),
  CONSTRAINT `viewer_ibfk_3` FOREIGN KEY (`routeset_id`) REFERENCES `routeset` (`id`),
  CONSTRAINT `viewer_ibfk_4` FOREIGN KEY (`policy_id`) REFERENCES `policy` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='map of <clientset , domain_pool> -> <policy, routeset>';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `viewer`
--

LOCK TABLES `viewer` WRITE;
/*!40000 ALTER TABLE `viewer` DISABLE KEYS */;
INSERT INTO `viewer` VALUES (1,1,1,1,1,1);
/*!40000 ALTER TABLE `viewer` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Final view structure for view `base_route_view`
--

/*!50001 DROP VIEW IF EXISTS `base_route_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `base_route_view` AS select `route`.`routeset_id` AS `routeset_id`,`routeset`.`name` AS `routeset_name`,`route`.`netlinkset_id` AS `netlinkset_id`,`netlinkset`.`name` AS `netlinkset_name`,`route`.`id` AS `route_id`,`route`.`priority` AS `route_priority`,`route`.`score` AS `route_score`,`route`.`outlink_id` AS `outlink_id`,`outlink`.`name` AS `outlink_name`,`outlink`.`addr` AS `outlink_addr`,`outlink`.`typ` AS `outlink_typ`,`outlink`.`enable` AS `outlink_enable`,`outlink`.`unavailable` AS `outlink_unavailable`,`route`.`enable` AS `route_enable`,`route`.`unavailable` AS `route_unavailable` from (((`netlinkset` join `outlink`) join `route`) join `routeset`) where ((`netlinkset`.`id` = `route`.`netlinkset_id`) and (`outlink`.`id` = `route`.`outlink_id`) and (`route`.`routeset_id` = `routeset`.`id`)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `best_route_view`
--

/*!50001 DROP VIEW IF EXISTS `best_route_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `best_route_view` AS select `base_route_view`.`routeset_id` AS `routeset_id`,`base_route_view`.`netlinkset_id` AS `netlinkset_id`,min(`base_route_view`.`route_priority`) AS `route_priority` from `base_route_view` where ((`base_route_view`.`outlink_enable` = 1) and (`base_route_view`.`route_enable` = 1) and (`base_route_view`.`route_unavailable` = 0) and (`base_route_view`.`outlink_unavailable` = 0) and (`base_route_view`.`route_score` <> 0)) group by `base_route_view`.`routeset_id`,`base_route_view`.`netlinkset_id` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `clientset_view`
--

/*!50001 DROP VIEW IF EXISTS `clientset_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `clientset_view` AS select `ipnet`.`id` AS `ipnet_id`,`ipnet`.`ip_start` AS `ip_start`,inet_aton(`ipnet`.`ip_start`) AS `ip_start_int`,`ipnet`.`ip_end` AS `ip_end`,inet_aton(`ipnet`.`ip_end`) AS `ip_end_int`,`ipnet`.`ipnet` AS `ipnet`,`ipnet`.`mask` AS `mask`,`ipnet`.`clientset_id` AS `clientset_id`,`clientset`.`name` AS `clientset_name` from (`ipnet` join `clientset` on((`clientset`.`id` = `ipnet`.`clientset_id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `clientset_wl_view`
--

/*!50001 DROP VIEW IF EXISTS `clientset_wl_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `clientset_wl_view` AS select `ipnet_wl`.`id` AS `ipnet_wl_id`,`ipnet_wl`.`ip_start` AS `ip_start`,inet_aton(`ipnet_wl`.`ip_start`) AS `ip_start_int`,`ipnet_wl`.`ip_end` AS `ip_end`,inet_aton(`ipnet_wl`.`ip_end`) AS `ip_end_int`,`ipnet_wl`.`ipnet` AS `ipnet`,`ipnet_wl`.`mask` AS `mask`,`ipnet_wl`.`clientset_id` AS `clientset_id`,`clientset`.`name` AS `clientset_name` from (`ipnet_wl` join `clientset` on((`clientset`.`id` = `ipnet_wl`.`clientset_id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `dns_forward_view`
--

/*!50001 DROP VIEW IF EXISTS `dns_forward_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `dns_forward_view` AS select `dns_forward_zone`.`id` AS `zone_id`,`dns_forward_zone`.`dm` AS `dm`,`dns_forward_zone`.`typ` AS `forward_typ`,`ldns`.`name` AS `ldns_name`,`ldns`.`addr` AS `ldns_addr` from ((`dns_forward_zone` join `dns_forwarders`) join `ldns`) where ((`ldns`.`id` = `dns_forwarders`.`ldns_id`) and (`dns_forward_zone`.`id` = `dns_forwarders`.`zone_id`)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `domain_view`
--

/*!50001 DROP VIEW IF EXISTS `domain_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `domain_view` AS select `domain`.`id` AS `domain_id`,`domain`.`domain` AS `domain`,`domain`.`domain_pool_id` AS `domain_pool_id`,`domain_pool`.`name` AS `pool_name`,`domain_pool`.`domain_monitor` AS `domain_monitor` from (`domain` join `domain_pool` on((`domain`.`domain_pool_id` = `domain_pool`.`id`))) where ((`domain_pool`.`enable` = 1) and (`domain_pool`.`unavailable` = 0)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `dst_view`
--

/*!50001 DROP VIEW IF EXISTS `dst_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `dst_view` AS select `domainlink`.`domain_pool_id` AS `domain_pool_id`,`domainlink`.`netlink_id` AS `netlink_id`,`domainlink`.`netlinkset_id` AS `netlinkset_id`,`domain_pool`.`name` AS `domain_pool_name`,`netlink`.`isp` AS `isp`,`netlink`.`region` AS `region`,`netlinkset`.`name` AS `netlinkset_name` from (((`domain_pool` join `netlink`) join `domainlink`) join `netlinkset`) where ((`domain_pool`.`id` = `domainlink`.`domain_pool_id`) and (`netlink`.`id` = `domainlink`.`netlink_id`) and (`domain_pool`.`enable` = 1) and (`domain_pool`.`unavailable` = 0) and (`domainlink`.`enable` = 1) and (`netlinkset`.`id` = `domainlink`.`netlinkset_id`)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `filter_view`
--

/*!50001 DROP VIEW IF EXISTS `filter_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `filter_view` AS select `filter`.`id` AS `id`,`filter`.`src_ip_start` AS `src_ip_start`,inet_aton(`filter`.`src_ip_start`) AS `src_ip_start_int`,`filter`.`src_ip_end` AS `src_ip_end`,inet_aton(`filter`.`src_ip_end`) AS `src_ip_end_int`,`filter`.`clientset_id` AS `clientset_id`,`clientset`.`name` AS `clientset_name`,`filter`.`domain_id` AS `domain_id`,`domain`.`domain` AS `domain`,`filter`.`domain_pool_id` AS `domain_pool_id`,`domain_pool`.`name` AS `domain_pool_name`,`filter`.`dst_ip_start` AS `dst_ip_start`,inet_aton(`filter`.`dst_ip_start`) AS `dst_ip_start_int`,`filter`.`dst_ip_end` AS `dst_ip_end`,inet_aton(`filter`.`dst_ip_end`) AS `dst_ip_end_int`,`filter`.`netlink_id` AS `netlink_id`,concat_ws(':',`netlink`.`isp`,`netlink`.`region`) AS `netlink_name`,`filter`.`target_ip` AS `target_ip`,inet_aton(`filter`.`target_ip`) AS `target_ip_int`,`filter`.`outlink_id` AS `outlink_id`,`outlink`.`name` AS `outlink_name`,`outlink`.`addr` AS `outlink_addr` from (((((`filter` join `domain`) join `domain_pool`) join `clientset`) join `netlink`) join `outlink`) where ((`filter`.`netlink_id` = `netlink`.`id`) and (`filter`.`clientset_id` = `clientset`.`id`) and (`filter`.`outlink_id` = `outlink`.`id`) and (`filter`.`domain_pool_id` = `domain_pool`.`id`) and (`filter`.`domain_id` = `domain`.`id`) and (`filter`.`enable` = 1)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `netlink_view`
--

/*!50001 DROP VIEW IF EXISTS `netlink_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `netlink_view` AS select `iptable`.`id` AS `iptable_id`,`iptable`.`ip_start` AS `ip_start`,inet_aton(`iptable`.`ip_start`) AS `ip_start_int`,`iptable`.`ip_end` AS `ip_end`,inet_aton(`iptable`.`ip_end`) AS `ip_end_int`,`iptable`.`ipnet` AS `ipnet`,`iptable`.`mask` AS `mask`,`iptable`.`netlink_id` AS `netlink_id`,`netlink`.`isp` AS `isp`,`netlink`.`region` AS `region`,`netlink`.`typ` AS `typ` from (`netlink` join `iptable` on((`iptable`.`netlink_id` = `netlink`.`id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `netlink_wl_view`
--

/*!50001 DROP VIEW IF EXISTS `netlink_wl_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `netlink_wl_view` AS select `iptable_wl`.`id` AS `iptable_wl_id`,`iptable_wl`.`ip_start` AS `ip_start`,inet_aton(`iptable_wl`.`ip_start`) AS `ip_start_int`,`iptable_wl`.`ip_end` AS `ip_end`,inet_aton(`iptable_wl`.`ip_end`) AS `ip_end_int`,`iptable_wl`.`ipnet` AS `ipnet`,`iptable_wl`.`mask` AS `mask`,`iptable_wl`.`netlink_id` AS `netlink_id`,`netlink`.`isp` AS `isp`,`netlink`.`region` AS `region`,`netlink`.`typ` AS `typ` from (`netlink` join `iptable_wl` on((`iptable_wl`.`netlink_id` = `netlink`.`id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `outlink_view`
--

/*!50001 DROP VIEW IF EXISTS `outlink_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `outlink_view` AS select `natlink`.`outlink_id` AS `outlink_id`,`outlink`.`addr` AS `outlink_addr`,`natlink`.`addr` AS `natlink_addr`,`natlink`.`natserver_id` AS `natserver_id`,`natserver`.`name` AS `nat_name`,`natlink`.`gw` AS `natlink_gw`,`natlink`.`status` AS `natlink_status` from ((`outlink` join `natlink`) join `natserver`) where ((`natlink`.`natserver_id` = `natserver`.`id`) and (`natlink`.`outlink_id` = `outlink`.`id`) and (`natserver`.`enable` = 1) and (`natserver`.`unavailable` = 0)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `policy_view`
--

/*!50001 DROP VIEW IF EXISTS `policy_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `policy_view` AS select `policy`.`id` AS `policy_id`,`policy`.`name` AS `policy_name`,`policy_detail`.`policy_sequence` AS `policy_sequence`,`policy_detail`.`priority` AS `priority`,`policy_detail`.`weight` AS `weight`,`policy_detail`.`op` AS `op`,`policy_detail`.`op_typ` AS `op_typ`,`policy_detail`.`ldns_id` AS `ldns_id`,`ldns`.`name` AS `name`,`ldns`.`addr` AS `addr`,`ldns`.`typ` AS `typ`,`policy_detail`.`rrset_id` AS `rrset_id` from ((`policy` join `policy_detail`) join `ldns`) where ((`policy`.`id` = `policy_detail`.`policy_id`) and (`ldns`.`id` = `policy_detail`.`ldns_id`) and (`ldns`.`unavailable` = 0) and (`ldns`.`enable` = 1) and (`policy_detail`.`enable` = 1)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `route_view`
--

/*!50001 DROP VIEW IF EXISTS `route_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `route_view` AS select `a`.`routeset_id` AS `routeset_id`,`a`.`routeset_name` AS `routeset_name`,`a`.`netlinkset_id` AS `netlinkset_id`,`a`.`netlinkset_name` AS `netlinkset_name`,`a`.`route_id` AS `route_id`,`a`.`route_priority` AS `route_priority`,`a`.`route_score` AS `route_score`,`a`.`outlink_id` AS `outlink_id`,`a`.`outlink_name` AS `outlink_name`,`a`.`outlink_addr` AS `outlink_addr`,`a`.`outlink_typ` AS `outlink_typ` from (`base_route_view` `a` join `best_route_view` `b`) where ((`a`.`routeset_id` = `b`.`routeset_id`) and (`a`.`netlinkset_id` = `b`.`netlinkset_id`) and (`a`.`route_priority` = `b`.`route_priority`)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `src_view`
--

/*!50001 DROP VIEW IF EXISTS `src_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `src_view` AS select `viewer`.`clientset_id` AS `clientset_id`,`clientset`.`name` AS `clientset_name`,`viewer`.`domain_pool_id` AS `domain_pool_id`,`domain_pool`.`name` AS `domain_pool_name`,`viewer`.`routeset_id` AS `routeset_id`,`routeset`.`name` AS `routeset_name`,`viewer`.`policy_id` AS `policy_id`,`policy`.`name` AS `policy_name` from ((((`clientset` join `domain_pool`) join `viewer`) join `routeset`) join `policy`) where ((`clientset`.`id` = `viewer`.`clientset_id`) and (`domain_pool`.`id` = `viewer`.`domain_pool_id`) and (`routeset`.`id` = `viewer`.`routeset_id`) and (`policy`.`id` = `viewer`.`policy_id`) and (`domain_pool`.`enable` = 1) and (`domain_pool`.`unavailable` = 0) and (`viewer`.`enable` = 1)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-12-21 14:10:35
