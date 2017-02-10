-- MySQL dump 10.13  Distrib 5.6.35, for Linux (i686)
--
-- Host: localhost    Database: igw
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
-- Temporary view structure for view `base_route_view`
--

DROP TABLE IF EXISTS `base_route_view`;
/*!50001 DROP VIEW IF EXISTS `base_route_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `base_route_view` AS SELECT 
 1 AS `ipset_id`,
 1 AS `ipset_name`,
 1 AS `domain_pool_id`,
 1 AS `domain_pool_name`,
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
-- Table structure for table `domain`
--

DROP TABLE IF EXISTS `domain`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `domain` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `domain` varchar(255) NOT NULL,
  `domain_pool_id` int(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `domain` (`domain`,`domain_pool_id`),
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
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `info` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='serve domains set';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `domain_pool`
--

LOCK TABLES `domain_pool` WRITE;
/*!40000 ALTER TABLE `domain_pool` DISABLE KEYS */;
INSERT INTO `domain_pool` VALUES (1,'global','Base domain pool for all of domains which are not specifically configured');
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
 1 AS `pool_name`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `domainlink`
--

DROP TABLE IF EXISTS `domainlink`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `domainlink` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `domain_pool_id` int(64) NOT NULL,
  `netlink_id` int(64) NOT NULL,
  `netlinkset_id` int(64) NOT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `domain_pool_id` (`domain_pool_id`,`netlink_id`),
  KEY `netlink_id` (`netlink_id`),
  KEY `netlinkset_id` (`netlinkset_id`),
  CONSTRAINT `domainlink_ibfk_1` FOREIGN KEY (`domain_pool_id`) REFERENCES `domain_pool` (`id`),
  CONSTRAINT `domainlink_ibfk_2` FOREIGN KEY (`netlink_id`) REFERENCES `netlink` (`id`),
  CONSTRAINT `domainlink_ibfk_3` FOREIGN KEY (`netlinkset_id`) REFERENCES `netlinkset` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='bind domain_pool and netlink to a netlinkset';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `domainlink`
--

LOCK TABLES `domainlink` WRITE;
/*!40000 ALTER TABLE `domainlink` DISABLE KEYS */;
/*!40000 ALTER TABLE `domainlink` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `dst_route_view`
--

DROP TABLE IF EXISTS `dst_route_view`;
/*!50001 DROP VIEW IF EXISTS `dst_route_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `dst_route_view` AS SELECT 
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
-- Table structure for table `filter`
--

DROP TABLE IF EXISTS `filter`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `filter` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `src_ip_start` varchar(40) DEFAULT NULL,
  `src_ip_end` varchar(40) DEFAULT NULL,
  `ipset_id` int(64) DEFAULT NULL,
  `domain_id` int(64) DEFAULT NULL,
  `dst_ip` varchar(40) DEFAULT NULL,
  `outlink_id` int(64) DEFAULT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `priority` int(11) NOT NULL DEFAULT '20',
  PRIMARY KEY (`id`),
  KEY `ipset_id` (`ipset_id`),
  KEY `domain_id` (`domain_id`),
  KEY `outlink_id` (`outlink_id`),
  CONSTRAINT `filter_ibfk_1` FOREIGN KEY (`ipset_id`) REFERENCES `ipset` (`id`) ON DELETE CASCADE,
  CONSTRAINT `filter_ibfk_2` FOREIGN KEY (`domain_id`) REFERENCES `domain` (`id`) ON DELETE CASCADE,
  CONSTRAINT `filter_ibfk_3` FOREIGN KEY (`outlink_id`) REFERENCES `outlink` (`id`) ON DELETE CASCADE
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
 1 AS `src_ip_end`,
 1 AS `ipset_id`,
 1 AS `domain_id`,
 1 AS `dst_ip`,
 1 AS `outlink_id`,
 1 AS `enable`,
 1 AS `priority`,
 1 AS `domain`,
 1 AS `outlink_name`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `ipnet`
--

DROP TABLE IF EXISTS `ipnet`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ipnet` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `ip_start` varchar(40) NOT NULL,
  `ip_end` varchar(40) NOT NULL,
  `ipnet` varchar(40) NOT NULL,
  `mask` int(8) NOT NULL,
  `priority` int(11) NOT NULL DEFAULT '20',
  `ipset_id` int(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_start` (`ip_start`,`ip_end`,`priority`),
  KEY `ipset_id` (`ipset_id`),
  CONSTRAINT `ipnet_ibfk_1` FOREIGN KEY (`ipset_id`) REFERENCES `ipset` (`id`) ON DELETE CASCADE
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
-- Table structure for table `ipset`
--

DROP TABLE IF EXISTS `ipset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ipset` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `info` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='client ipnet set';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ipset`
--

LOCK TABLES `ipset` WRITE;
/*!40000 ALTER TABLE `ipset` DISABLE KEYS */;
INSERT INTO `ipset` VALUES (1,'unknown','src ipnet that igw has no idea where it belongs to');
/*!40000 ALTER TABLE `ipset` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `ipset_view`
--

DROP TABLE IF EXISTS `ipset_view`;
/*!50001 DROP VIEW IF EXISTS `ipset_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `ipset_view` AS SELECT 
 1 AS `ipnet_id`,
 1 AS `ip_start`,
 1 AS `ip_end`,
 1 AS `ipnet`,
 1 AS `mask`,
 1 AS `priority`,
 1 AS `ipset_id`,
 1 AS `ipset_name`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `iptable`
--

DROP TABLE IF EXISTS `iptable`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `iptable` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `ip_start` varchar(40) NOT NULL,
  `ip_end` varchar(40) NOT NULL,
  `ipnet` varchar(40) NOT NULL,
  `mask` int(8) NOT NULL,
  `priority` int(11) NOT NULL DEFAULT '20',
  `netlink_id` int(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_start` (`ip_start`,`ip_end`,`priority`),
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
-- Table structure for table `ldns`
--

DROP TABLE IF EXISTS `ldns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ldns` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `addr` varchar(255) NOT NULL,
  `typ` varchar(16) NOT NULL DEFAULT 'A',
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `unavailable` int(16) NOT NULL DEFAULT '0' COMMENT 'if other than zero, ldns is unavailable with each bit indicate different reason',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='upstream ldns info';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ldns`
--

LOCK TABLES `ldns` WRITE;
/*!40000 ALTER TABLE `ldns` DISABLE KEYS */;
/*!40000 ALTER TABLE `ldns` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `netlink`
--

DROP TABLE IF EXISTS `netlink`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `netlink` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `isp` varchar(255) NOT NULL,
  `region` varchar(255) NOT NULL,
  `typ` varchar(32) NOT NULL DEFAULT 'normal',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='netlink (isp + province or CP) of a target ip';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `netlink`
--

LOCK TABLES `netlink` WRITE;
/*!40000 ALTER TABLE `netlink` DISABLE KEYS */;
INSERT INTO `netlink` VALUES (1,'unknown','unknown','normal');
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
 1 AS `ip_end`,
 1 AS `ipnet`,
 1 AS `mask`,
 1 AS `priority`,
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
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='netlink set ';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `netlinkset`
--

LOCK TABLES `netlinkset` WRITE;
/*!40000 ALTER TABLE `netlinkset` DISABLE KEYS */;
/*!40000 ALTER TABLE `netlinkset` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `netlinkset_view`
--

DROP TABLE IF EXISTS `netlinkset_view`;
/*!50001 DROP VIEW IF EXISTS `netlinkset_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `netlinkset_view` AS SELECT 
 1 AS `domain_pool_id`,
 1 AS `netlink_id`,
 1 AS `netlinkset_id`,
 1 AS `domain_pool_name`,
 1 AS `isp`,
 1 AS `region`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `outlink`
--

DROP TABLE IF EXISTS `outlink`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `outlink` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `addr` varchar(255) NOT NULL,
  `typ` varchar(32) NOT NULL DEFAULT 'normal',
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `unavailable` int(16) NOT NULL DEFAULT '0' COMMENT 'if other than zero, outlink is unavailable, each bit indicate different reason',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='network gateway, aka outlink';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `outlink`
--

LOCK TABLES `outlink` WRITE;
/*!40000 ALTER TABLE `outlink` DISABLE KEYS */;
/*!40000 ALTER TABLE `outlink` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `policy`
--

DROP TABLE IF EXISTS `policy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `policy` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='policy index of choose ldns upstream forwarder';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `policy`
--

LOCK TABLES `policy` WRITE;
/*!40000 ALTER TABLE `policy` DISABLE KEYS */;
/*!40000 ALTER TABLE `policy` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `policy_detail`
--

DROP TABLE IF EXISTS `policy_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `policy_detail` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `policy_id` int(64) NOT NULL,
  `policy_sequence` int(11) NOT NULL DEFAULT '0',
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `priority` int(11) NOT NULL DEFAULT '20',
  `weight` int(11) NOT NULL DEFAULT '100',
  `op` varchar(255) NOT NULL DEFAULT 'and',
  `op_typ` varchar(64) NOT NULL DEFAULT 'builtin',
  `ldns_id` int(64) NOT NULL,
  `rrset_id` int(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `ldns_id` (`ldns_id`),
  KEY `rrset_id` (`rrset_id`),
  KEY `policy_id` (`policy_id`),
  CONSTRAINT `policy_detail_ibfk_1` FOREIGN KEY (`ldns_id`) REFERENCES `ldns` (`id`),
  CONSTRAINT `policy_detail_ibfk_2` FOREIGN KEY (`rrset_id`) REFERENCES `rrset` (`id`),
  CONSTRAINT `policy_detail_ibfk_3` FOREIGN KEY (`policy_id`) REFERENCES `policy` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='policy detail of choose ldns upstream forwarder';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `policy_detail`
--

LOCK TABLES `policy_detail` WRITE;
/*!40000 ALTER TABLE `policy_detail` DISABLE KEYS */;
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
 1 AS `ipset_id`,
 1 AS `domain_pool_id`,
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
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `outlink_id` int(64) NOT NULL,
  `netlinkset_id` int(64) NOT NULL,
  `routeset_id` int(64) NOT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `priority` int(11) NOT NULL DEFAULT '20',
  `score` int(11) DEFAULT NULL COMMENT 'netlink performance index',
  `unavailable` int(16) NOT NULL DEFAULT '0' COMMENT 'if other than zero, route is unavailable, each bit indicate different reason',
  PRIMARY KEY (`id`),
  UNIQUE KEY `netlinkset_id` (`netlinkset_id`,`outlink_id`,`routeset_id`),
  KEY `outlink_id` (`outlink_id`),
  KEY `routeset_id` (`routeset_id`),
  CONSTRAINT `route_ibfk_1` FOREIGN KEY (`outlink_id`) REFERENCES `outlink` (`id`),
  CONSTRAINT `route_ibfk_2` FOREIGN KEY (`netlinkset_id`) REFERENCES `netlinkset` (`id`),
  CONSTRAINT `route_ibfk_3` FOREIGN KEY (`routeset_id`) REFERENCES `routeset` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Performance of using gateway to serve paricular netlink';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `route`
--

LOCK TABLES `route` WRITE;
/*!40000 ALTER TABLE `route` DISABLE KEYS */;
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
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `info` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `routeset`
--

LOCK TABLES `routeset` WRITE;
/*!40000 ALTER TABLE `routeset` DISABLE KEYS */;
/*!40000 ALTER TABLE `routeset` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rr`
--

DROP TABLE IF EXISTS `rr`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rr` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `rrtype` int(16) unsigned NOT NULL,
  `rrdata` varchar(255) NOT NULL,
  `ttl` int(32) unsigned NOT NULL DEFAULT '300',
  `rrset_id` int(64) DEFAULT NULL,
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
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '',
  `ttl` int(32) unsigned NOT NULL DEFAULT '300',
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
 1 AS `ipset_id`,
 1 AS `ipset_name`,
 1 AS `domain_pool_id`,
 1 AS `domain_pool_name`,
 1 AS `routeset_id`,
 1 AS `routeset_name`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `the_view`
--

DROP TABLE IF EXISTS `the_view`;
/*!50001 DROP VIEW IF EXISTS `the_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `the_view` AS SELECT 
 1 AS `ipset_id`,
 1 AS `ipset_name`,
 1 AS `domain_pool_id`,
 1 AS `domain_pool_name`,
 1 AS `routeset_id`,
 1 AS `routeset_name`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `viewer`
--

DROP TABLE IF EXISTS `viewer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `viewer` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `ipset_id` int(64) NOT NULL,
  `domain_pool_id` int(64) NOT NULL,
  `routeset_id` int(64) NOT NULL,
  `policy_id` int(64) NOT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `ipset_id` (`ipset_id`),
  KEY `domain_pool_id` (`domain_pool_id`),
  KEY `routeset_id` (`routeset_id`),
  KEY `policy_id` (`policy_id`),
  CONSTRAINT `viewer_ibfk_1` FOREIGN KEY (`ipset_id`) REFERENCES `ipset` (`id`),
  CONSTRAINT `viewer_ibfk_2` FOREIGN KEY (`domain_pool_id`) REFERENCES `domain_pool` (`id`),
  CONSTRAINT `viewer_ibfk_3` FOREIGN KEY (`routeset_id`) REFERENCES `routeset` (`id`),
  CONSTRAINT `viewer_ibfk_4` FOREIGN KEY (`policy_id`) REFERENCES `policy` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='map of <ipset , domain_pool> -> <policy, routeset>';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `viewer`
--

LOCK TABLES `viewer` WRITE;
/*!40000 ALTER TABLE `viewer` DISABLE KEYS */;
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
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `base_route_view` AS select `viewer`.`ipset_id` AS `ipset_id`,`ipset`.`name` AS `ipset_name`,`viewer`.`domain_pool_id` AS `domain_pool_id`,`domain_pool`.`name` AS `domain_pool_name`,`viewer`.`routeset_id` AS `routeset_id`,`routeset`.`name` AS `routeset_name`,`route`.`netlinkset_id` AS `netlinkset_id`,`netlinkset`.`name` AS `netlinkset_name`,`route`.`id` AS `route_id`,min(`route`.`priority`) AS `route_priority`,max(`route`.`score`) AS `route_score`,`route`.`outlink_id` AS `outlink_id`,`outlink`.`name` AS `outlink_name`,`outlink`.`addr` AS `outlink_addr`,`outlink`.`typ` AS `outlink_typ` from ((((((`ipset` join `domain_pool`) join `viewer`) join `routeset`) join `netlinkset`) join `route`) join `outlink`) where ((`ipset`.`id` = `viewer`.`ipset_id`) and (`domain_pool`.`id` = `viewer`.`domain_pool_id`) and (`viewer`.`routeset_id` = `routeset`.`id`) and (`netlinkset`.`id` = `route`.`netlinkset_id`) and (`outlink`.`id` = `route`.`outlink_id`) and (`route`.`routeset_id` = `routeset`.`id`) and (`outlink`.`enable` = 1) and (`route`.`enable` = 1) and (`outlink`.`unavailable` = 0) and (`route`.`unavailable` = 0) and (`viewer`.`enable` = 1)) group by `viewer`.`routeset_id` */;
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
/*!50001 VIEW `domain_view` AS select `domain`.`id` AS `domain_id`,`domain`.`domain` AS `domain`,`domain`.`domain_pool_id` AS `domain_pool_id`,`domain_pool`.`name` AS `pool_name` from (`domain` join `domain_pool` on((`domain`.`domain_pool_id` = `domain_pool`.`id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `dst_route_view`
--

/*!50001 DROP VIEW IF EXISTS `dst_route_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `dst_route_view` AS select `route`.`routeset_id` AS `routeset_id`,`routeset`.`name` AS `routeset_name`,`route`.`netlinkset_id` AS `netlinkset_id`,`netlinkset`.`name` AS `netlinkset_name`,`route`.`id` AS `route_id`,min(`route`.`priority`) AS `route_priority`,max(`route`.`score`) AS `route_score`,`route`.`outlink_id` AS `outlink_id`,`outlink`.`name` AS `outlink_name`,`outlink`.`addr` AS `outlink_addr`,`outlink`.`typ` AS `outlink_typ` from (((`netlinkset` join `outlink`) join `route`) join `routeset`) where ((`netlinkset`.`id` = `route`.`netlinkset_id`) and (`outlink`.`id` = `route`.`outlink_id`) and (`route`.`routeset_id` = `routeset`.`id`) and (`outlink`.`enable` = 1) and (`route`.`enable` = 1) and (`outlink`.`unavailable` = 0) and (`route`.`unavailable` = 0)) group by `route`.`routeset_id` */;
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
/*!50001 VIEW `filter_view` AS select `filter`.`id` AS `id`,`filter`.`src_ip_start` AS `src_ip_start`,`filter`.`src_ip_end` AS `src_ip_end`,`filter`.`ipset_id` AS `ipset_id`,`filter`.`domain_id` AS `domain_id`,`filter`.`dst_ip` AS `dst_ip`,`filter`.`outlink_id` AS `outlink_id`,`filter`.`enable` AS `enable`,`filter`.`priority` AS `priority`,`domain`.`domain` AS `domain`,`outlink`.`name` AS `outlink_name` from (((`domain` join `filter`) join `outlink`) join `ipset`) where ((`filter`.`ipset_id` = `ipset`.`id`) and (`filter`.`outlink_id` = `outlink`.`id`) and (`filter`.`domain_id` = `domain`.`id`) and (`filter`.`enable` = 1)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `ipset_view`
--

/*!50001 DROP VIEW IF EXISTS `ipset_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `ipset_view` AS select `ipnet`.`id` AS `ipnet_id`,`ipnet`.`ip_start` AS `ip_start`,`ipnet`.`ip_end` AS `ip_end`,`ipnet`.`ipnet` AS `ipnet`,`ipnet`.`mask` AS `mask`,`ipnet`.`priority` AS `priority`,`ipnet`.`ipset_id` AS `ipset_id`,`ipset`.`name` AS `ipset_name` from (`ipnet` join `ipset` on((`ipset`.`id` = `ipnet`.`ipset_id`))) */;
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
/*!50001 VIEW `netlink_view` AS select `iptable`.`id` AS `iptable_id`,`iptable`.`ip_start` AS `ip_start`,`iptable`.`ip_end` AS `ip_end`,`iptable`.`ipnet` AS `ipnet`,`iptable`.`mask` AS `mask`,`iptable`.`priority` AS `priority`,`iptable`.`netlink_id` AS `netlink_id`,`netlink`.`isp` AS `isp`,`netlink`.`region` AS `region`,`netlink`.`typ` AS `typ` from (`netlink` join `iptable` on((`iptable`.`netlink_id` = `netlink`.`id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `netlinkset_view`
--

/*!50001 DROP VIEW IF EXISTS `netlinkset_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `netlinkset_view` AS select `domainlink`.`domain_pool_id` AS `domain_pool_id`,`domainlink`.`netlink_id` AS `netlink_id`,`domainlink`.`netlinkset_id` AS `netlinkset_id`,`domain_pool`.`name` AS `domain_pool_name`,`netlink`.`isp` AS `isp`,`netlink`.`region` AS `region` from ((`domain_pool` join `netlink`) join `domainlink`) where ((`domain_pool`.`id` = `domainlink`.`domain_pool_id`) and (`netlink`.`id` = `domainlink`.`netlink_id`) and (`domainlink`.`enable` = 1)) */;
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
/*!50001 VIEW `policy_view` AS select `viewer`.`ipset_id` AS `ipset_id`,`viewer`.`domain_pool_id` AS `domain_pool_id`,`viewer`.`policy_id` AS `policy_id`,`policy`.`name` AS `policy_name`,`policy_detail`.`policy_sequence` AS `policy_sequence`,`policy_detail`.`priority` AS `priority`,`policy_detail`.`weight` AS `weight`,`policy_detail`.`op` AS `op`,`policy_detail`.`op_typ` AS `op_typ`,`policy_detail`.`ldns_id` AS `ldns_id`,`ldns`.`name` AS `name`,`ldns`.`addr` AS `addr`,`ldns`.`typ` AS `typ`,`policy_detail`.`rrset_id` AS `rrset_id` from (((`viewer` join `policy`) join `policy_detail`) join `ldns`) where ((`policy`.`id` = `policy_detail`.`policy_id`) and (`ldns`.`id` = `policy_detail`.`ldns_id`) and (`viewer`.`policy_id` = `policy`.`id`) and (`ldns`.`unavailable` = 0) and (`ldns`.`enable` = 1) and (`viewer`.`enable` = 1) and (`policy_detail`.`enable` = 1)) */;
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
/*!50001 VIEW `route_view` AS select `route`.`routeset_id` AS `routeset_id`,`routeset`.`name` AS `routeset_name`,`route`.`netlinkset_id` AS `netlinkset_id`,`netlinkset`.`name` AS `netlinkset_name`,`route`.`id` AS `route_id`,min(`route`.`priority`) AS `route_priority`,max(`route`.`score`) AS `route_score`,`route`.`outlink_id` AS `outlink_id`,`outlink`.`name` AS `outlink_name`,`outlink`.`addr` AS `outlink_addr`,`outlink`.`typ` AS `outlink_typ` from (((`netlinkset` join `outlink`) join `route`) join `routeset`) where ((`netlinkset`.`id` = `route`.`netlinkset_id`) and (`outlink`.`id` = `route`.`outlink_id`) and (`route`.`routeset_id` = `routeset`.`id`) and (`outlink`.`enable` = 1) and (`route`.`enable` = 1) and (`outlink`.`unavailable` = 0) and (`route`.`unavailable` = 0)) group by `route`.`routeset_id` */;
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
/*!50001 VIEW `src_view` AS select `viewer`.`ipset_id` AS `ipset_id`,`ipset`.`name` AS `ipset_name`,`viewer`.`domain_pool_id` AS `domain_pool_id`,`domain_pool`.`name` AS `domain_pool_name`,`viewer`.`routeset_id` AS `routeset_id`,`routeset`.`name` AS `routeset_name` from (((`ipset` join `domain_pool`) join `viewer`) join `routeset`) where ((`ipset`.`id` = `viewer`.`ipset_id`) and (`domain_pool`.`id` = `viewer`.`domain_pool_id`) and (`routeset`.`id` = `viewer`.`routeset_id`) and (`viewer`.`enable` = 1)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `the_view`
--

/*!50001 DROP VIEW IF EXISTS `the_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `the_view` AS select `viewer`.`ipset_id` AS `ipset_id`,`ipset`.`name` AS `ipset_name`,`viewer`.`domain_pool_id` AS `domain_pool_id`,`domain_pool`.`name` AS `domain_pool_name`,`viewer`.`routeset_id` AS `routeset_id`,`routeset`.`name` AS `routeset_name` from (((`ipset` join `domain_pool`) join `viewer`) join `routeset`) where ((`ipset`.`id` = `viewer`.`ipset_id`) and (`domain_pool`.`id` = `viewer`.`domain_pool_id`) and (`routeset`.`id` = `viewer`.`routeset_id`) and (`viewer`.`enable` = 1)) */;
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

-- Dump completed on 2017-02-10 18:07:06
