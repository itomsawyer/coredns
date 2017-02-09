-- MySQL dump 10.13  Distrib 5.5.50, for debian-linux-gnu (i686)
--
-- Host: localhost    Database: igw
-- ------------------------------------------------------
-- Server version	5.5.50-0+deb8u1

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
  `info` text NOT NULL,
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
-- Temporary table structure for view `domain_view`
--

DROP TABLE IF EXISTS `domain_view`;
/*!50001 DROP VIEW IF EXISTS `domain_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `domain_view` (
  `domain_id` tinyint NOT NULL,
  `domain` tinyint NOT NULL,
  `domain_pool_id` tinyint NOT NULL,
  `pool_name` tinyint NOT NULL
) ENGINE=MyISAM */;
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
-- Table structure for table `gw`
--

DROP TABLE IF EXISTS `gw`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gw` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `addr` varchar(255) NOT NULL,
  `typ` varchar(32) NOT NULL DEFAULT 'normal',
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `unavailable` int(16) NOT NULL DEFAULT '0' COMMENT 'if other than zero, gw is unavailable, each bit indicate different reason',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='network gateway, aka outlink';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gw`
--

LOCK TABLES `gw` WRITE;
/*!40000 ALTER TABLE `gw` DISABLE KEYS */;
/*!40000 ALTER TABLE `gw` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gwlink`
--

DROP TABLE IF EXISTS `gwlink`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gwlink` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `gw_id` int(64) NOT NULL,
  `netlinkset_id` int(64) NOT NULL,
  `gwlinkset_id` int(64) NOT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `priority` int(11) NOT NULL DEFAULT '0',
  `score` int(11) DEFAULT NULL COMMENT 'netlink performance index',
  `unavailable` int(16) NOT NULL DEFAULT '0' COMMENT 'if other than zero, gwlink is unavailable, each bit indicate different reason',
  PRIMARY KEY (`id`),
  UNIQUE KEY `netlinkset_id` (`netlinkset_id`,`gw_id`,`gwlinkset_id`),
  KEY `gw_id` (`gw_id`),
  KEY `gwlinkset_id` (`gwlinkset_id`),
  CONSTRAINT `gwlink_ibfk_1` FOREIGN KEY (`gw_id`) REFERENCES `gw` (`id`),
  CONSTRAINT `gwlink_ibfk_2` FOREIGN KEY (`netlinkset_id`) REFERENCES `netlinkset` (`id`),
  CONSTRAINT `gwlink_ibfk_3` FOREIGN KEY (`gwlinkset_id`) REFERENCES `gwlinkset` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Performance of using gateway to serve paricular netlink';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gwlink`
--

LOCK TABLES `gwlink` WRITE;
/*!40000 ALTER TABLE `gwlink` DISABLE KEYS */;
/*!40000 ALTER TABLE `gwlink` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary table structure for view `gwlink_view`
--

DROP TABLE IF EXISTS `gwlink_view`;
/*!50001 DROP VIEW IF EXISTS `gwlink_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `gwlink_view` (
  `gwlinkset_id` tinyint NOT NULL,
  `gwlinkset_name` tinyint NOT NULL,
  `netlinkset_id` tinyint NOT NULL,
  `netlinkset_name` tinyint NOT NULL,
  `gwlink_id` tinyint NOT NULL,
  `gwlink_priority` tinyint NOT NULL,
  `gwlink_score` tinyint NOT NULL,
  `gw_id` tinyint NOT NULL,
  `gw_name` tinyint NOT NULL,
  `gw_addr` tinyint NOT NULL,
  `gw_typ` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `gwlinkset`
--

DROP TABLE IF EXISTS `gwlinkset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gwlinkset` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `info` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gwlinkset`
--

LOCK TABLES `gwlinkset` WRITE;
/*!40000 ALTER TABLE `gwlinkset` DISABLE KEYS */;
/*!40000 ALTER TABLE `gwlinkset` ENABLE KEYS */;
UNLOCK TABLES;

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
  `priority` int(11) NOT NULL DEFAULT '0',
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
  `info` text NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='client ipnet set';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ipset`
--

LOCK TABLES `ipset` WRITE;
/*!40000 ALTER TABLE `ipset` DISABLE KEYS */;
/*!40000 ALTER TABLE `ipset` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary table structure for view `ipset_view`
--

DROP TABLE IF EXISTS `ipset_view`;
/*!50001 DROP VIEW IF EXISTS `ipset_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `ipset_view` (
  `ipnet_id` tinyint NOT NULL,
  `ip_start` tinyint NOT NULL,
  `ip_end` tinyint NOT NULL,
  `ipnet` tinyint NOT NULL,
  `mask` tinyint NOT NULL,
  `priority` tinyint NOT NULL,
  `ipset_id` tinyint NOT NULL,
  `ipset_name` tinyint NOT NULL
) ENGINE=MyISAM */;
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
  `priority` int(11) NOT NULL DEFAULT '0',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='netlink (isp + province or CP) of a target ip';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `netlink`
--

LOCK TABLES `netlink` WRITE;
/*!40000 ALTER TABLE `netlink` DISABLE KEYS */;
/*!40000 ALTER TABLE `netlink` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary table structure for view `netlink_view`
--

DROP TABLE IF EXISTS `netlink_view`;
/*!50001 DROP VIEW IF EXISTS `netlink_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `netlink_view` (
  `iptable_id` tinyint NOT NULL,
  `ip_start` tinyint NOT NULL,
  `ip_end` tinyint NOT NULL,
  `ipnet` tinyint NOT NULL,
  `mask` tinyint NOT NULL,
  `priority` tinyint NOT NULL,
  `netlink_id` tinyint NOT NULL,
  `isp` tinyint NOT NULL,
  `region` tinyint NOT NULL,
  `typ` tinyint NOT NULL
) ENGINE=MyISAM */;
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
-- Temporary table structure for view `netlinkset_view`
--

DROP TABLE IF EXISTS `netlinkset_view`;
/*!50001 DROP VIEW IF EXISTS `netlinkset_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `netlinkset_view` (
  `domain_pool_id` tinyint NOT NULL,
  `netlink_id` tinyint NOT NULL,
  `netlinkset_id` tinyint NOT NULL,
  `domain_pool_name` tinyint NOT NULL,
  `isp` tinyint NOT NULL,
  `region` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

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
  `priority` int(11) NOT NULL DEFAULT '0',
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
-- Temporary table structure for view `policy_view`
--

DROP TABLE IF EXISTS `policy_view`;
/*!50001 DROP VIEW IF EXISTS `policy_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `policy_view` (
  `ipset_id` tinyint NOT NULL,
  `domain_pool_id` tinyint NOT NULL,
  `policy_id` tinyint NOT NULL,
  `policy_name` tinyint NOT NULL,
  `policy_sequence` tinyint NOT NULL,
  `priority` tinyint NOT NULL,
  `weight` tinyint NOT NULL,
  `op` tinyint NOT NULL,
  `op_typ` tinyint NOT NULL,
  `ldns_id` tinyint NOT NULL,
  `name` tinyint NOT NULL,
  `addr` tinyint NOT NULL,
  `typ` tinyint NOT NULL,
  `rrset_id` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `route_view`
--

DROP TABLE IF EXISTS `route_view`;
/*!50001 DROP VIEW IF EXISTS `route_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `route_view` (
  `ipset_id` tinyint NOT NULL,
  `ipset_name` tinyint NOT NULL,
  `domain_pool_id` tinyint NOT NULL,
  `domain_pool_name` tinyint NOT NULL,
  `gwlinkset_id` tinyint NOT NULL,
  `gwlinkset_name` tinyint NOT NULL,
  `netlinkset_id` tinyint NOT NULL,
  `netlinkset_name` tinyint NOT NULL,
  `gwlink_id` tinyint NOT NULL,
  `gwlink_priority` tinyint NOT NULL,
  `gwlink_score` tinyint NOT NULL,
  `gw_id` tinyint NOT NULL,
  `gw_name` tinyint NOT NULL,
  `gw_addr` tinyint NOT NULL,
  `gw_typ` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

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
-- Temporary table structure for view `the_view`
--

DROP TABLE IF EXISTS `the_view`;
/*!50001 DROP VIEW IF EXISTS `the_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `the_view` (
  `ipset_id` tinyint NOT NULL,
  `ipset_name` tinyint NOT NULL,
  `domain_pool_id` tinyint NOT NULL,
  `domain_pool_name` tinyint NOT NULL,
  `gwlinkset_id` tinyint NOT NULL,
  `gwlinkset_name` tinyint NOT NULL
) ENGINE=MyISAM */;
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
  `gwlinkset_id` int(64) NOT NULL,
  `policy_id` int(64) NOT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `ipset_id` (`ipset_id`),
  KEY `domain_pool_id` (`domain_pool_id`),
  KEY `gwlinkset_id` (`gwlinkset_id`),
  KEY `policy_id` (`policy_id`),
  CONSTRAINT `viewer_ibfk_1` FOREIGN KEY (`ipset_id`) REFERENCES `ipset` (`id`),
  CONSTRAINT `viewer_ibfk_2` FOREIGN KEY (`domain_pool_id`) REFERENCES `domain_pool` (`id`),
  CONSTRAINT `viewer_ibfk_3` FOREIGN KEY (`gwlinkset_id`) REFERENCES `gwlinkset` (`id`),
  CONSTRAINT `viewer_ibfk_4` FOREIGN KEY (`policy_id`) REFERENCES `policy` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='map of <ipset , domain_pool> -> <policy, gwlinkset>';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `viewer`
--

LOCK TABLES `viewer` WRITE;
/*!40000 ALTER TABLE `viewer` DISABLE KEYS */;
/*!40000 ALTER TABLE `viewer` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Final view structure for view `domain_view`
--

/*!50001 DROP TABLE IF EXISTS `domain_view`*/;
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
-- Final view structure for view `gwlink_view`
--

/*!50001 DROP TABLE IF EXISTS `gwlink_view`*/;
/*!50001 DROP VIEW IF EXISTS `gwlink_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `gwlink_view` AS select `gwlink`.`gwlinkset_id` AS `gwlinkset_id`,`gwlinkset`.`name` AS `gwlinkset_name`,`gwlink`.`netlinkset_id` AS `netlinkset_id`,`netlinkset`.`name` AS `netlinkset_name`,`gwlink`.`id` AS `gwlink_id`,min(`gwlink`.`priority`) AS `gwlink_priority`,max(`gwlink`.`score`) AS `gwlink_score`,`gwlink`.`gw_id` AS `gw_id`,`gw`.`name` AS `gw_name`,`gw`.`addr` AS `gw_addr`,`gw`.`typ` AS `gw_typ` from (((`netlinkset` join `gw`) join `gwlink`) join `gwlinkset`) where ((`netlinkset`.`id` = `gwlink`.`netlinkset_id`) and (`gw`.`id` = `gwlink`.`gw_id`) and (`gwlink`.`gwlinkset_id` = `gwlinkset`.`id`) and (`gw`.`enable` = 1) and (`gwlink`.`enable` = 1) and (`gw`.`unavailable` = 0) and (`gwlink`.`unavailable` = 0)) group by `gwlink`.`gwlinkset_id` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `ipset_view`
--

/*!50001 DROP TABLE IF EXISTS `ipset_view`*/;
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

/*!50001 DROP TABLE IF EXISTS `netlink_view`*/;
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

/*!50001 DROP TABLE IF EXISTS `netlinkset_view`*/;
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

/*!50001 DROP TABLE IF EXISTS `policy_view`*/;
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

/*!50001 DROP TABLE IF EXISTS `route_view`*/;
/*!50001 DROP VIEW IF EXISTS `route_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `route_view` AS select `viewer`.`ipset_id` AS `ipset_id`,`ipset`.`name` AS `ipset_name`,`viewer`.`domain_pool_id` AS `domain_pool_id`,`domain_pool`.`name` AS `domain_pool_name`,`viewer`.`gwlinkset_id` AS `gwlinkset_id`,`gwlinkset`.`name` AS `gwlinkset_name`,`gwlink`.`netlinkset_id` AS `netlinkset_id`,`netlinkset`.`name` AS `netlinkset_name`,`gwlink`.`id` AS `gwlink_id`,min(`gwlink`.`priority`) AS `gwlink_priority`,max(`gwlink`.`score`) AS `gwlink_score`,`gwlink`.`gw_id` AS `gw_id`,`gw`.`name` AS `gw_name`,`gw`.`addr` AS `gw_addr`,`gw`.`typ` AS `gw_typ` from ((((((`ipset` join `domain_pool`) join `viewer`) join `gwlinkset`) join `netlinkset`) join `gwlink`) join `gw`) where ((`ipset`.`id` = `viewer`.`ipset_id`) and (`domain_pool`.`id` = `viewer`.`domain_pool_id`) and (`viewer`.`gwlinkset_id` = `gwlinkset`.`id`) and (`netlinkset`.`id` = `gwlink`.`netlinkset_id`) and (`gw`.`id` = `gwlink`.`gw_id`) and (`gwlink`.`gwlinkset_id` = `gwlinkset`.`id`) and (`gw`.`enable` = 1) and (`gwlink`.`enable` = 1) and (`gw`.`unavailable` = 0) and (`gwlink`.`unavailable` = 0) and (`viewer`.`enable` = 1)) group by `viewer`.`gwlinkset_id` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `the_view`
--

/*!50001 DROP TABLE IF EXISTS `the_view`*/;
/*!50001 DROP VIEW IF EXISTS `the_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=MERGE */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `the_view` AS select `viewer`.`ipset_id` AS `ipset_id`,`ipset`.`name` AS `ipset_name`,`viewer`.`domain_pool_id` AS `domain_pool_id`,`domain_pool`.`name` AS `domain_pool_name`,`viewer`.`gwlinkset_id` AS `gwlinkset_id`,`gwlinkset`.`name` AS `gwlinkset_name` from (((`ipset` join `domain_pool`) join `viewer`) join `gwlinkset`) where ((`ipset`.`id` = `viewer`.`ipset_id`) and (`domain_pool`.`id` = `viewer`.`domain_pool_id`) and (`gwlinkset`.`id` = `viewer`.`gwlinkset_id`) and (`viewer`.`enable` = 1)) */;
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

-- Dump completed on 2017-02-09 17:33:57
