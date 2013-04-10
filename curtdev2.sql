-- MySQL dump 10.13  Distrib 5.6.10, for osx10.7 (i386)
--
-- Host: curtsql.cloudapp.net    Database: CurtDev2
-- ------------------------------------------------------
-- Server version	5.6.10-log

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
-- Table structure for table `AcesType`
--

DROP TABLE IF EXISTS `AcesType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `AcesType` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ApiAccess`
--

DROP TABLE IF EXISTS `ApiAccess`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ApiAccess` (
  `id` varchar(64) NOT NULL,
  `key_id` varchar(64) NOT NULL,
  `module_id` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  KEY `FK__ApiAccess__key_i__628F3CBE` (`key_id`),
  KEY `FK__ApiAccess__modul__638360F7` (`module_id`),
  CONSTRAINT `FK__ApiAccess__key_i__628F3CBE` FOREIGN KEY (`key_id`) REFERENCES `ApiKey` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__ApiAccess__modul__638360F7` FOREIGN KEY (`module_id`) REFERENCES `ApiModules` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ApiKey`
--

DROP TABLE IF EXISTS `ApiKey`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ApiKey` (
  `id` varchar(64) NOT NULL,
  `api_key` varchar(64) NOT NULL,
  `type_id` varchar(64) NOT NULL,
  `user_id` varchar(64) NOT NULL,
  `date_added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  KEY `FK__ApiKey__type_id__5AEE1AF6` (`type_id`),
  KEY `FK__ApiKey__user_id__5BE23F2F` (`user_id`),
  CONSTRAINT `FK__ApiKey__type_id__5AEE1AF6` FOREIGN KEY (`type_id`) REFERENCES `ApiKeyType` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__ApiKey__user_id__5BE23F2F` FOREIGN KEY (`user_id`) REFERENCES `CustomerUser` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`%`*/ /*!50003 trigger before_update_api_key
before update on ApiKey
for each row
set new.api_key = uuid() */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Table structure for table `ApiKeyType`
--

DROP TABLE IF EXISTS `ApiKeyType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ApiKeyType` (
  `id` varchar(64) NOT NULL,
  `type` varchar(500) DEFAULT NULL,
  `date_added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ApiModules`
--

DROP TABLE IF EXISTS `ApiModules`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ApiModules` (
  `id` varchar(64) NOT NULL,
  `name` varchar(500) DEFAULT NULL,
  `access_level` varchar(64) DEFAULT NULL,
  `date_added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `AuthAccess`
--

DROP TABLE IF EXISTS `AuthAccess`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `AuthAccess` (
  `id` varchar(64) NOT NULL,
  `userID` varchar(64) NOT NULL,
  `AreaID` varchar(64) NOT NULL,
  `dateAdded` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `AuthAreas`
--

DROP TABLE IF EXISTS `AuthAreas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `AuthAreas` (
  `id` varchar(64) NOT NULL,
  `path` varchar(50) NOT NULL,
  `DomainID` varchar(64) NOT NULL,
  `name` varchar(50) NOT NULL,
  `parentAreaID` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `AuthDomains`
--

DROP TABLE IF EXISTS `AuthDomains`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `AuthDomains` (
  `id` varchar(64) NOT NULL,
  `url` varchar(50) NOT NULL,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `AuthorizedTracking`
--

DROP TABLE IF EXISTS `AuthorizedTracking`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `AuthorizedTracking` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `property` varchar(500) NOT NULL,
  `view_count` int(11) NOT NULL,
  `authorized_id` varchar(500) DEFAULT NULL,
  `date_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `date_modified` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Authors`
--

DROP TABLE IF EXISTS `Authors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Authors` (
  `authorID` int(11) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(510) DEFAULT NULL,
  `last_name` varchar(510) DEFAULT NULL,
  `email` varchar(510) DEFAULT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`authorID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Banners`
--

DROP TABLE IF EXISTS `Banners`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Banners` (
  `bannerID` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) DEFAULT NULL,
  `link` varchar(400) DEFAULT NULL,
  `starts` datetime DEFAULT NULL,
  `ends` datetime DEFAULT NULL,
  `path` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`bannerID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `BaseVehicle`
--

DROP TABLE IF EXISTS `BaseVehicle`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `BaseVehicle` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `AAIABaseVehicleID` int(11) DEFAULT NULL,
  `YearID` int(11) NOT NULL,
  `MakeID` int(11) NOT NULL,
  `ModelID` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_BaseVehicle_IX` (`AAIABaseVehicleID`),
  KEY `FK__BaseVehic__MakeI__75B5891D` (`MakeID`),
  KEY `FK__BaseVehic__Model__76A9AD56` (`ModelID`),
  KEY `FK_BaseVehicle_Year` (`YearID`),
  CONSTRAINT `FK_BaseVehicle_Year` FOREIGN KEY (`YearID`) REFERENCES `vcdb_Year` (`YearID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__BaseVehic__MakeI__75B5891D` FOREIGN KEY (`MakeID`) REFERENCES `vcdb_Make` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__BaseVehic__Model__76A9AD56` FOREIGN KEY (`ModelID`) REFERENCES `vcdb_Model` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=22334 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `BlogCategories`
--

DROP TABLE IF EXISTS `BlogCategories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `BlogCategories` (
  `blogCategoryID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `slug` varchar(255) DEFAULT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`blogCategoryID`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `BlogPost_BlogCategory`
--

DROP TABLE IF EXISTS `BlogPost_BlogCategory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `BlogPost_BlogCategory` (
  `postCategoryID` int(11) NOT NULL AUTO_INCREMENT,
  `blogPostID` int(11) NOT NULL,
  `blogCategoryID` int(11) NOT NULL,
  PRIMARY KEY (`postCategoryID`),
  KEY `FK__BlogPost___blogC__57DD0BE4` (`blogCategoryID`),
  KEY `FK__BlogPost___blogP__58D1301D` (`blogPostID`),
  CONSTRAINT `FK__BlogPost___blogC__57DD0BE4` FOREIGN KEY (`blogCategoryID`) REFERENCES `BlogCategories` (`blogCategoryID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__BlogPost___blogP__58D1301D` FOREIGN KEY (`blogPostID`) REFERENCES `BlogPosts` (`blogPostID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `BlogPosts`
--

DROP TABLE IF EXISTS `BlogPosts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `BlogPosts` (
  `blogPostID` int(11) NOT NULL AUTO_INCREMENT,
  `post_title` varchar(500) NOT NULL,
  `slug` varchar(500) NOT NULL,
  `post_text` longtext,
  `publishedDate` datetime DEFAULT NULL,
  `createdDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `lastModified` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `userID` int(11) NOT NULL,
  `meta_title` varchar(510) DEFAULT NULL,
  `meta_description` varchar(510) DEFAULT NULL,
  `keywords` longtext,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`blogPostID`),
  KEY `BlogPostAuthorID` (`userID`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Brand`
--

DROP TABLE IF EXISTS `Brand`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Brand` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `code` varchar(255) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `BrandPart`
--

DROP TABLE IF EXISTS `BrandPart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `BrandPart` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `brandID` int(11) NOT NULL,
  `brandPartID` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  `interchangeType` char(1) NOT NULL,
  `dateAdded` datetime NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `BusinessClass`
--

DROP TABLE IF EXISTS `BusinessClass`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `BusinessClass` (
  `BusinessClassID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `sort` int(11) NOT NULL,
  `showOnWebsite` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`BusinessClassID`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Cabelas`
--

DROP TABLE IF EXISTS `Cabelas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Cabelas` (
  `cabelasID` int(11) NOT NULL AUTO_INCREMENT,
  `priceCode` int(11) DEFAULT NULL,
  `cabelasPart` varchar(50) NOT NULL,
  PRIMARY KEY (`cabelasID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `CartIntegration`
--

DROP TABLE IF EXISTS `CartIntegration`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `CartIntegration` (
  `referenceID` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `custPartID` int(11) NOT NULL,
  `custID` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`referenceID`)
) ENGINE=InnoDB AUTO_INCREMENT=65097 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `CatPart`
--

DROP TABLE IF EXISTS `CatPart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `CatPart` (
  `catPartID` int(11) NOT NULL AUTO_INCREMENT,
  `catID` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  PRIMARY KEY (`catPartID`),
  KEY `IX_CatPart_Cat_Part` (`catID`,`partID`),
  KEY `FK__CatPart__partID__54945AAA` (`partID`),
  KEY `cat_idx` (`catID`),
  KEY `part_idx` (`partID`),
  CONSTRAINT `FK__CatPart__catID__55887EE3` FOREIGN KEY (`catID`) REFERENCES `Categories` (`catID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__CatPart__partID__54945AAA` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=5026 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Categories`
--

DROP TABLE IF EXISTS `Categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Categories` (
  `catID` int(11) NOT NULL AUTO_INCREMENT,
  `dateAdded` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `parentID` int(11) NOT NULL,
  `catTitle` varchar(100) DEFAULT NULL,
  `shortDesc` varchar(255) DEFAULT NULL,
  `longDesc` longtext,
  `image` varchar(255) DEFAULT NULL,
  `isLifestyle` int(11) NOT NULL,
  `codeID` int(11) NOT NULL DEFAULT '0',
  `sort` int(11) NOT NULL DEFAULT '1',
  `vehicleSpecific` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`catID`),
  KEY `IX_Categories_ParentID` (`parentID`),
  KEY `IX_Categories_Sort` (`sort`),
  KEY `idx` (`catID`),
  KEY `title_idx` (`catTitle`)
) ENGINE=InnoDB AUTO_INCREMENT=277 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Class`
--

DROP TABLE IF EXISTS `Class`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Class` (
  `classID` int(11) NOT NULL AUTO_INCREMENT,
  `class` varchar(255) DEFAULT NULL,
  `image` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`classID`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ColorCode`
--

DROP TABLE IF EXISTS `ColorCode`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ColorCode` (
  `codeID` int(11) NOT NULL,
  `code` varchar(100) DEFAULT NULL,
  `font` varchar(100) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Comments`
--

DROP TABLE IF EXISTS `Comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Comments` (
  `commentID` int(11) NOT NULL AUTO_INCREMENT,
  `blogPostID` int(11) NOT NULL,
  `name` varchar(510) NOT NULL,
  `email` varchar(510) DEFAULT NULL,
  `comment_text` longtext,
  `createdDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `approved` tinyint(1) NOT NULL DEFAULT '0',
  `active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`commentID`),
  KEY `FK__Comments__blogPo__56E8E7AB` (`blogPostID`),
  CONSTRAINT `FK__Comments__blogPo__56E8E7AB` FOREIGN KEY (`blogPostID`) REFERENCES `BlogPosts` (`blogPostID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Company`
--

DROP TABLE IF EXISTS `Company`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Company` (
  `companyID` int(11) NOT NULL AUTO_INCREMENT,
  `company_image` varchar(500) DEFAULT NULL,
  `logo_image` varchar(500) DEFAULT NULL,
  `tagline` varchar(255) DEFAULT NULL,
  `youtube_link` varchar(255) DEFAULT NULL,
  `facebook_link` varchar(255) DEFAULT NULL,
  `twitter_link` varchar(255) DEFAULT NULL,
  `contact_email` varchar(300) DEFAULT NULL,
  `name` varchar(500) DEFAULT NULL,
  `homepage_lookup` int(11) NOT NULL DEFAULT '0',
  `adwords` varchar(500) DEFAULT NULL,
  `merchant_provider` varchar(100) DEFAULT NULL,
  `merchant_id` varchar(200) DEFAULT NULL,
  `analytics_id` varchar(200) NOT NULL DEFAULT '',
  `testimonial_submission` varchar(100) NOT NULL DEFAULT 'Closed',
  `moderate_blog` tinyint(1) NOT NULL DEFAULT '0',
  `stylesheet` varchar(200) NOT NULL DEFAULT 'light_layout.css',
  PRIMARY KEY (`companyID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ConfigAttribute`
--

DROP TABLE IF EXISTS `ConfigAttribute`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ConfigAttribute` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `ConfigAttributeTypeID` int(11) NOT NULL,
  `parentID` int(11) NOT NULL,
  `vcdbID` int(11) DEFAULT NULL,
  `value` varchar(255) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_ConfigAttribute_IX` (`ConfigAttributeTypeID`,`parentID`),
  CONSTRAINT `FK__ConfigAtt__Confi__07D43958` FOREIGN KEY (`ConfigAttributeTypeID`) REFERENCES `ConfigAttributeType` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=293 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ConfigAttributeType`
--

DROP TABLE IF EXISTS `ConfigAttributeType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ConfigAttributeType` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `AcesTypeID` int(11) DEFAULT NULL,
  `sort` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `FK__ConfigAtt__AcesT__030F843B` (`AcesTypeID`),
  CONSTRAINT `FK__ConfigAtt__AcesT__030F843B` FOREIGN KEY (`AcesTypeID`) REFERENCES `AcesType` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=74 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Contact`
--

DROP TABLE IF EXISTS `Contact`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Contact` (
  `contactID` int(11) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `phone` varchar(30) DEFAULT NULL,
  `subject` varchar(500) DEFAULT NULL,
  `message` longtext,
  `createdDate` datetime NOT NULL,
  `type` varchar(255) DEFAULT NULL,
  `address1` varchar(500) DEFAULT NULL,
  `address2` varchar(500) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `state` varchar(10) DEFAULT NULL,
  `postalcode` varchar(20) DEFAULT NULL,
  `country` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`contactID`)
) ENGINE=InnoDB AUTO_INCREMENT=2147483647 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ContactReceiver`
--

DROP TABLE IF EXISTS `ContactReceiver`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ContactReceiver` (
  `contactReceiverID` int(11) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`contactReceiverID`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ContactReceiver_ContactType`
--

DROP TABLE IF EXISTS `ContactReceiver_ContactType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ContactReceiver_ContactType` (
  `receiverTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `contactReceiverID` int(11) NOT NULL,
  `contactTypeID` int(11) NOT NULL,
  PRIMARY KEY (`receiverTypeID`),
  KEY `FK__ContactRe__conta__6FB49575` (`contactReceiverID`),
  KEY `FK__ContactRe__conta__70A8B9AE` (`contactTypeID`),
  CONSTRAINT `FK__ContactRe__conta__6FB49575` FOREIGN KEY (`contactReceiverID`) REFERENCES `ContactReceiver` (`contactReceiverID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__ContactRe__conta__70A8B9AE` FOREIGN KEY (`contactTypeID`) REFERENCES `ContactType` (`contactTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=111 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ContactType`
--

DROP TABLE IF EXISTS `ContactType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ContactType` (
  `contactTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`contactTypeID`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Content`
--

DROP TABLE IF EXISTS `Content`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Content` (
  `contentID` int(11) NOT NULL AUTO_INCREMENT,
  `text` longtext,
  `cTypeID` int(11) NOT NULL,
  PRIMARY KEY (`contentID`),
  KEY `FK__Content__cTypeID__0B457116` (`cTypeID`),
  CONSTRAINT `FK__Content__cTypeID__0B457116` FOREIGN KEY (`cTypeID`) REFERENCES `ContentType` (`cTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=296001 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ContentBridge`
--

DROP TABLE IF EXISTS `ContentBridge`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ContentBridge` (
  `cBridgeID` int(11) NOT NULL AUTO_INCREMENT,
  `catID` int(11) DEFAULT NULL,
  `partID` int(11) DEFAULT NULL,
  `contentID` int(11) NOT NULL,
  PRIMARY KEY (`cBridgeID`),
  KEY `IX_ContentBridge_catIDContent` (`catID`,`contentID`),
  KEY `IX_ContentBridge_partIDContent` (`partID`,`contentID`),
  KEY `FK__ContentBr__conte__390C3BC6` (`contentID`),
  CONSTRAINT `FK__ContentBr__catID__3A005FFF` FOREIGN KEY (`catID`) REFERENCES `Categories` (`catID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__ContentBr__conte__390C3BC6` FOREIGN KEY (`contentID`) REFERENCES `Content` (`contentID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__ContentBr__partI__3AF48438` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=26230 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ContentType`
--

DROP TABLE IF EXISTS `ContentType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ContentType` (
  `cTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(255) DEFAULT NULL,
  `allowHTML` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`cTypeID`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Country`
--

DROP TABLE IF EXISTS `Country`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Country` (
  `countryID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `abbr` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`countryID`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `CustUserWebProperties`
--

DROP TABLE IF EXISTS `CustUserWebProperties`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `CustUserWebProperties` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userID` varchar(64) NOT NULL,
  `webPropID` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=165 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Customer`
--

DROP TABLE IF EXISTS `Customer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Customer` (
  `cust_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `address` varchar(500) DEFAULT NULL,
  `city` varchar(150) DEFAULT NULL,
  `stateID` int(11) NOT NULL,
  `phone` varchar(50) DEFAULT NULL,
  `fax` varchar(50) DEFAULT NULL,
  `contact_person` varchar(300) DEFAULT NULL,
  `dealer_type` int(11) NOT NULL,
  `latitude` varchar(200) DEFAULT NULL,
  `longitude` varchar(200) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `website` varchar(500) DEFAULT NULL,
  `customerID` int(11) DEFAULT NULL,
  `isDummy` tinyint(1) NOT NULL DEFAULT '0',
  `parentID` int(11) DEFAULT NULL,
  `searchURL` varchar(500) DEFAULT NULL,
  `eLocalURL` varchar(500) DEFAULT NULL,
  `logo` varchar(500) DEFAULT NULL,
  `address2` varchar(500) DEFAULT NULL,
  `postal_code` varchar(25) DEFAULT NULL,
  `mCodeID` int(11) NOT NULL DEFAULT '1',
  `salesRepID` int(11) DEFAULT NULL,
  `APIKey` varchar(64) DEFAULT NULL,
  `tier` int(11) NOT NULL DEFAULT '1',
  `showWebsite` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`cust_id`),
  KEY `CustomerCustomerID` (`customerID`),
  KEY `IX_CustomerID` (`customerID`)
) ENGINE=InnoDB AUTO_INCREMENT=10443192 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `CustomerCost`
--

DROP TABLE IF EXISTS `CustomerCost`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `CustomerCost` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `cust_id` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  `cost` decimal(18,2) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `CustomerLocations`
--

DROP TABLE IF EXISTS `CustomerLocations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `CustomerLocations` (
  `locationID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(500) DEFAULT NULL,
  `address` varchar(500) DEFAULT NULL,
  `city` varchar(500) DEFAULT NULL,
  `stateID` int(11) NOT NULL,
  `email` varchar(500) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `fax` varchar(20) DEFAULT NULL,
  `latitude` double NOT NULL,
  `longitude` double NOT NULL,
  `cust_id` int(11) NOT NULL DEFAULT '0',
  `contact_person` varchar(300) DEFAULT NULL,
  `isprimary` tinyint(1) NOT NULL DEFAULT '0',
  `postalCode` varchar(30) DEFAULT NULL,
  `ShippingDefault` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`locationID`),
  KEY `IX_CustomerLocations_Customer` (`cust_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7727 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `CustomerPricing`
--

DROP TABLE IF EXISTS `CustomerPricing`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `CustomerPricing` (
  `cust_price_id` int(11) NOT NULL AUTO_INCREMENT,
  `cust_id` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  `price` decimal(8,2) DEFAULT NULL,
  `isSale` int(11) NOT NULL DEFAULT '0',
  `sale_start` date DEFAULT NULL,
  `sale_end` date DEFAULT NULL,
  PRIMARY KEY (`cust_price_id`)
) ENGINE=InnoDB AUTO_INCREMENT=330807 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `CustomerReport`
--

DROP TABLE IF EXISTS `CustomerReport`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `CustomerReport` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `customerID` int(11) NOT NULL,
  `created` datetime NOT NULL,
  `ReportTypeID` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `FK__CustomerR__Repor__0F604C87` (`ReportTypeID`),
  CONSTRAINT `FK__CustomerR__Repor__0F604C87` FOREIGN KEY (`ReportTypeID`) REFERENCES `ReportType` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `CustomerReportPart`
--

DROP TABLE IF EXISTS `CustomerReportPart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `CustomerReportPart` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `customerID` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=279 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `CustomerUser`
--

DROP TABLE IF EXISTS `CustomerUser`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `CustomerUser` (
  `id` varchar(64) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `customerID` int(11) DEFAULT NULL,
  `date_added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `active` tinyint(1) NOT NULL DEFAULT '0',
  `locationID` int(11) NOT NULL DEFAULT '0',
  `isSudo` tinyint(1) NOT NULL DEFAULT '0',
  `cust_ID` int(11) NOT NULL,
  `NotCustomer` tinyint(1) DEFAULT NULL,
  `proper_password` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `DealerTiers`
--

DROP TABLE IF EXISTS `DealerTiers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `DealerTiers` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `tier` varchar(255) DEFAULT NULL,
  `sort` int(11) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `DealerTypes`
--

DROP TABLE IF EXISTS `DealerTypes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `DealerTypes` (
  `dealer_type` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(100) DEFAULT NULL,
  `online` tinyint(1) NOT NULL DEFAULT '0',
  `show` tinyint(1) NOT NULL DEFAULT '1',
  `label` varchar(300) DEFAULT NULL,
  PRIMARY KEY (`dealer_type`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `FAQ`
--

DROP TABLE IF EXISTS `FAQ`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `FAQ` (
  `faqID` int(11) NOT NULL AUTO_INCREMENT,
  `question` varchar(500) DEFAULT NULL,
  `answer` longtext,
  PRIMARY KEY (`faqID`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `File`
--

DROP TABLE IF EXISTS `File`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `File` (
  `fileID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(500) NOT NULL,
  `path` varchar(500) NOT NULL,
  `height` int(11) NOT NULL,
  `width` int(11) NOT NULL,
  `createdDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `fileGalleryID` int(11) NOT NULL,
  `fileExtID` int(11) NOT NULL DEFAULT '1',
  `size` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`fileID`),
  KEY `FK__File__fileExtID__6C390A4C` (`fileExtID`),
  CONSTRAINT `FK__File__fileExtID__6C390A4C` FOREIGN KEY (`fileExtID`) REFERENCES `FileExt` (`fileExtID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1286 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `FileExt`
--

DROP TABLE IF EXISTS `FileExt`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `FileExt` (
  `fileExtID` int(11) NOT NULL AUTO_INCREMENT,
  `fileExt` varchar(10) NOT NULL,
  `fileExtIcon` varchar(1000) DEFAULT NULL,
  `fileTypeID` int(11) NOT NULL,
  PRIMARY KEY (`fileExtID`),
  KEY `FK__FileExt__fileTyp__6A50C1DA` (`fileTypeID`),
  CONSTRAINT `FK__FileExt__fileTyp__6A50C1DA` FOREIGN KEY (`fileTypeID`) REFERENCES `FileType` (`fileTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `FileGallery`
--

DROP TABLE IF EXISTS `FileGallery`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `FileGallery` (
  `fileGalleryID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` varchar(4000) DEFAULT NULL,
  `parentID` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`fileGalleryID`)
) ENGINE=InnoDB AUTO_INCREMENT=98 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `FileType`
--

DROP TABLE IF EXISTS `FileType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `FileType` (
  `fileTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `fileType` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`fileTypeID`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ForumGroup`
--

DROP TABLE IF EXISTS `ForumGroup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ForumGroup` (
  `forumGroupID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` longtext,
  `createdDate` datetime NOT NULL,
  PRIMARY KEY (`forumGroupID`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ForumPost`
--

DROP TABLE IF EXISTS `ForumPost`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ForumPost` (
  `postID` int(11) NOT NULL AUTO_INCREMENT,
  `parentID` int(11) NOT NULL,
  `threadID` int(11) NOT NULL,
  `createdDate` datetime NOT NULL,
  `title` varchar(255) NOT NULL,
  `post` longtext NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `company` varchar(255) DEFAULT NULL,
  `notify` tinyint(1) NOT NULL,
  `approved` tinyint(1) NOT NULL,
  `active` tinyint(1) NOT NULL,
  `IPAddress` varchar(255) NOT NULL,
  `flag` tinyint(1) NOT NULL,
  `sticky` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`postID`),
  KEY `FK__ForumPost__threa__22B5E1E5` (`threadID`),
  CONSTRAINT `FK__ForumPost__threa__22B5E1E5` FOREIGN KEY (`threadID`) REFERENCES `ForumThread` (`threadID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ForumThread`
--

DROP TABLE IF EXISTS `ForumThread`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ForumThread` (
  `threadID` int(11) NOT NULL AUTO_INCREMENT,
  `topicID` int(11) NOT NULL,
  `createdDate` datetime NOT NULL,
  `active` tinyint(1) NOT NULL,
  `closed` tinyint(1) NOT NULL,
  PRIMARY KEY (`threadID`),
  KEY `FK__ForumThre__topic__1DF12CC8` (`topicID`),
  CONSTRAINT `FK__ForumThre__topic__1DF12CC8` FOREIGN KEY (`topicID`) REFERENCES `ForumTopic` (`topicID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ForumTopic`
--

DROP TABLE IF EXISTS `ForumTopic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ForumTopic` (
  `topicID` int(11) NOT NULL AUTO_INCREMENT,
  `TopicGroupID` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` longtext,
  `image` varchar(255) DEFAULT NULL,
  `createdDate` datetime NOT NULL,
  `active` tinyint(1) NOT NULL,
  `closed` tinyint(1) NOT NULL,
  PRIMARY KEY (`topicID`),
  KEY `FK__ForumTopi__Topic__192C77AB` (`TopicGroupID`),
  CONSTRAINT `FK__ForumTopi__Topic__192C77AB` FOREIGN KEY (`TopicGroupID`) REFERENCES `ForumGroup` (`forumGroupID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Gallery`
--

DROP TABLE IF EXISTS `Gallery`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Gallery` (
  `imgID` int(11) NOT NULL AUTO_INCREMENT,
  `img_path` varchar(500) NOT NULL,
  `title` varchar(200) DEFAULT NULL,
  `sort_order` int(11) NOT NULL,
  PRIMARY KEY (`imgID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `IPBlock`
--

DROP TABLE IF EXISTS `IPBlock`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IPBlock` (
  `blockID` int(11) NOT NULL AUTO_INCREMENT,
  `IPAddress` varchar(255) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `createdDate` datetime NOT NULL,
  `userID` int(11) NOT NULL,
  PRIMARY KEY (`blockID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `IncludedPart`
--

DROP TABLE IF EXISTS `IncludedPart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IncludedPart` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `includedID` int(11) NOT NULL,
  `quantity` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `FK__IncludedP__partI__33DF3ACD` (`partID`),
  KEY `FK__IncludedP__inclu__34D35F06` (`includedID`),
  CONSTRAINT `FK__IncludedP__inclu__34D35F06` FOREIGN KEY (`includedID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__IncludedP__partI__33DF3ACD` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `KioskOrderItems`
--

DROP TABLE IF EXISTS `KioskOrderItems`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `KioskOrderItems` (
  `itemID` int(11) NOT NULL AUTO_INCREMENT,
  `orderID` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  `quantity` int(11) NOT NULL,
  `price` decimal(19,4) NOT NULL,
  `isFulfilled` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`itemID`)
) ENGINE=InnoDB AUTO_INCREMENT=120 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `KioskOrders`
--

DROP TABLE IF EXISTS `KioskOrders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `KioskOrders` (
  `orderID` int(11) NOT NULL AUTO_INCREMENT,
  `order_date` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `vehicleID` int(11) DEFAULT NULL,
  `acctID` int(11) DEFAULT NULL,
  `fname` varchar(255) DEFAULT NULL,
  `lname` varchar(500) DEFAULT NULL,
  `email` varchar(300) DEFAULT NULL,
  `phone` varchar(100) DEFAULT NULL,
  `isProcessed` int(11) NOT NULL DEFAULT '0',
  `isHidden` int(11) NOT NULL DEFAULT '0',
  `address` varchar(600) DEFAULT NULL,
  `city` varchar(500) DEFAULT NULL,
  `stateID` int(11) NOT NULL DEFAULT '0',
  `zip` varchar(100) DEFAULT NULL,
  `locationID` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`orderID`)
) ENGINE=InnoDB AUTO_INCREMENT=10000077 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `LandingPage`
--

DROP TABLE IF EXISTS `LandingPage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `LandingPage` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `websiteID` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `startDate` datetime NOT NULL,
  `endDate` datetime NOT NULL,
  `url` varchar(255) NOT NULL,
  `pageContent` longtext,
  `linkClasses` varchar(255) DEFAULT NULL,
  `conversionID` varchar(150) DEFAULT NULL,
  `conversionLabel` varchar(150) DEFAULT NULL,
  `newWindow` tinyint(1) NOT NULL DEFAULT '0',
  `menuPosition` varchar(15) NOT NULL DEFAULT 'top',
  PRIMARY KEY (`id`),
  KEY `FK__LandingPa__websi__509AA9B5` (`websiteID`),
  CONSTRAINT `FK__LandingPa__websi__509AA9B5` FOREIGN KEY (`websiteID`) REFERENCES `Website` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `LandingPageData`
--

DROP TABLE IF EXISTS `LandingPageData`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `LandingPageData` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `landingPageID` int(11) NOT NULL,
  `dataKey` varchar(100) NOT NULL,
  `dataValue` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `FK__LandingPa__landi__5A2413EF` (`landingPageID`),
  CONSTRAINT `FK__LandingPa__landi__5A2413EF` FOREIGN KEY (`landingPageID`) REFERENCES `LandingPage` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `LandingPageImages`
--

DROP TABLE IF EXISTS `LandingPageImages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `LandingPageImages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `landingPageID` int(11) NOT NULL,
  `url` varchar(255) NOT NULL,
  `sort` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `FK__LandingPa__landi__555F5ED2` (`landingPageID`),
  CONSTRAINT `FK__LandingPa__landi__555F5ED2` FOREIGN KEY (`landingPageID`) REFERENCES `LandingPage` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Lifestyle_Trailer`
--

DROP TABLE IF EXISTS `Lifestyle_Trailer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Lifestyle_Trailer` (
  `lifestyleTrailerID` int(11) NOT NULL AUTO_INCREMENT,
  `catID` int(11) NOT NULL,
  `trailerID` int(11) NOT NULL,
  PRIMARY KEY (`lifestyleTrailerID`),
  KEY `FK__Lifestyle__catID__0869046B` (`catID`),
  KEY `FK__Lifestyle__trail__095D28A4` (`trailerID`),
  CONSTRAINT `FK__Lifestyle__catID__0869046B` FOREIGN KEY (`catID`) REFERENCES `Categories` (`catID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__Lifestyle__trail__095D28A4` FOREIGN KEY (`trailerID`) REFERENCES `Trailer` (`trailerID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=165 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Location_Services`
--

DROP TABLE IF EXISTS `Location_Services`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Location_Services` (
  `loc_service_id` int(11) NOT NULL AUTO_INCREMENT,
  `serviceID` int(11) NOT NULL,
  `locationID` int(11) NOT NULL,
  PRIMARY KEY (`loc_service_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Locations`
--

DROP TABLE IF EXISTS `Locations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Locations` (
  `locationID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `phone` varchar(15) DEFAULT NULL,
  `fax` varchar(15) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `address` varchar(500) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `stateID` int(11) NOT NULL,
  `zip` int(11) DEFAULT NULL,
  `isPrimary` int(11) NOT NULL DEFAULT '0',
  `latitude` decimal(18,8) NOT NULL DEFAULT '0.00000000',
  `longitude` decimal(18,8) NOT NULL DEFAULT '0.00000000',
  `places_status` varchar(10) DEFAULT NULL,
  `places_reference` varchar(300) DEFAULT NULL,
  `places_id` varchar(300) DEFAULT NULL,
  PRIMARY KEY (`locationID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Logger`
--

DROP TABLE IF EXISTS `Logger`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Logger` (
  `id` varchar(64) NOT NULL,
  `Message` varchar(500) DEFAULT NULL,
  `Source` longtext,
  `StackTrace` longtext,
  `TargetSite` varchar(500) DEFAULT NULL,
  `loggedType` varchar(64) NOT NULL,
  `Date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  KEY `FK_Logger_LoggerTypes` (`loggedType`),
  CONSTRAINT `FK_Logger_LoggerTypes` FOREIGN KEY (`loggedType`) REFERENCES `LoggerTypes` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `LoggerTypes`
--

DROP TABLE IF EXISTS `LoggerTypes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `LoggerTypes` (
  `id` varchar(64) NOT NULL,
  `type` varchar(200) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Make`
--

DROP TABLE IF EXISTS `Make`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Make` (
  `makeID` int(11) NOT NULL AUTO_INCREMENT,
  `make` varchar(255) NOT NULL,
  PRIMARY KEY (`makeID`)
) ENGINE=InnoDB AUTO_INCREMENT=54 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `MakeModel`
--

DROP TABLE IF EXISTS `MakeModel`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `MakeModel` (
  `mmID` int(11) NOT NULL AUTO_INCREMENT,
  `makeID` int(11) NOT NULL,
  `modelID` int(11) NOT NULL,
  PRIMARY KEY (`mmID`),
  KEY `IX_MakeModel` (`makeID`,`modelID`),
  KEY `FK__MakeModel__model__4977ADB9` (`modelID`),
  CONSTRAINT `FK__MakeModel__makeI__48838980` FOREIGN KEY (`makeID`) REFERENCES `Make` (`makeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__MakeModel__model__4977ADB9` FOREIGN KEY (`modelID`) REFERENCES `Model` (`modelID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=733 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `MapIcons`
--

DROP TABLE IF EXISTS `MapIcons`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `MapIcons` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `tier` int(11) NOT NULL,
  `dealer_type` int(11) NOT NULL,
  `mapicon` varchar(300) NOT NULL,
  `mapiconshadow` varchar(300) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `FK__MapIcons__tier__4F707E31` (`tier`),
  KEY `FK__MapIcons__dealer__5064A26A` (`dealer_type`),
  CONSTRAINT `FK__MapIcons__dealer__5064A26A` FOREIGN KEY (`dealer_type`) REFERENCES `DealerTypes` (`dealer_type`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__MapIcons__tier__4F707E31` FOREIGN KEY (`tier`) REFERENCES `DealerTiers` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `MapPolygon`
--

DROP TABLE IF EXISTS `MapPolygon`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `MapPolygon` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `stateID` int(11) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=148 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `MapPolygonCoordinates`
--

DROP TABLE IF EXISTS `MapPolygonCoordinates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `MapPolygonCoordinates` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `MapPolygonID` int(11) NOT NULL,
  `latitude` double NOT NULL,
  `longitude` double NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=16225 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `MapixCode`
--

DROP TABLE IF EXISTS `MapixCode`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `MapixCode` (
  `mCodeID` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`mCodeID`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Menu`
--

DROP TABLE IF EXISTS `Menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Menu` (
  `menuID` int(11) NOT NULL AUTO_INCREMENT,
  `menu_name` varchar(255) NOT NULL,
  `isPrimary` tinyint(1) NOT NULL DEFAULT '0',
  `active` tinyint(1) NOT NULL DEFAULT '1',
  `display_name` varchar(255) DEFAULT NULL,
  `requireAuthentication` tinyint(1) NOT NULL DEFAULT '0',
  `showOnSitemap` tinyint(1) NOT NULL DEFAULT '0',
  `sort` int(11) NOT NULL DEFAULT '1',
  `websiteID` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`menuID`),
  KEY `FK__Menu__websiteID__1FCC89D1` (`websiteID`),
  CONSTRAINT `FK__Menu__websiteID__1FCC89D1` FOREIGN KEY (`websiteID`) REFERENCES `Website` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Menu_SiteContent`
--

DROP TABLE IF EXISTS `Menu_SiteContent`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Menu_SiteContent` (
  `menuContentID` int(11) NOT NULL AUTO_INCREMENT,
  `menuID` int(11) NOT NULL,
  `contentID` int(11) DEFAULT NULL,
  `menuSort` int(11) NOT NULL,
  `menuTitle` varchar(255) DEFAULT NULL,
  `menuLink` varchar(500) DEFAULT NULL,
  `parentID` int(11) DEFAULT NULL,
  `linkTarget` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`menuContentID`),
  KEY `FK__Menu_Site__menuI__208CD6FA` (`menuID`),
  KEY `FK__Menu_Site__conte__2180FB33` (`contentID`),
  KEY `FK__Menu_Site__paren__22751F6C` (`parentID`),
  CONSTRAINT `FK__Menu_Site__conte__2180FB33` FOREIGN KEY (`contentID`) REFERENCES `SiteContent` (`contentID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__Menu_Site__menuI__208CD6FA` FOREIGN KEY (`menuID`) REFERENCES `Menu` (`menuID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__Menu_Site__paren__22751F6C` FOREIGN KEY (`parentID`) REFERENCES `Menu_SiteContent` (`menuContentID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=130 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Model`
--

DROP TABLE IF EXISTS `Model`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Model` (
  `modelID` int(11) NOT NULL AUTO_INCREMENT,
  `model` varchar(255) NOT NULL,
  PRIMARY KEY (`modelID`)
) ENGINE=InnoDB AUTO_INCREMENT=704 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ModelStyle`
--

DROP TABLE IF EXISTS `ModelStyle`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ModelStyle` (
  `msID` int(11) NOT NULL AUTO_INCREMENT,
  `modelID` int(11) NOT NULL,
  `styleID` int(11) NOT NULL,
  PRIMARY KEY (`msID`),
  KEY `IX_ModelStyle` (`modelID`,`styleID`),
  KEY `FK__ModelStyl__style__4B5FF62B` (`styleID`),
  CONSTRAINT `FK__ModelStyl__model__4A6BD1F2` FOREIGN KEY (`modelID`) REFERENCES `Model` (`modelID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__ModelStyl__style__4B5FF62B` FOREIGN KEY (`styleID`) REFERENCES `Style` (`styleID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1371 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Modules`
--

DROP TABLE IF EXISTS `Modules`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Modules` (
  `moduleID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) DEFAULT NULL,
  `path` varchar(100) DEFAULT NULL,
  `image` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`moduleID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `NewsItem`
--

DROP TABLE IF EXISTS `NewsItem`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `NewsItem` (
  `newsItemID` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(500) DEFAULT NULL,
  `lead` longtext,
  `content` longtext,
  `publishStart` datetime DEFAULT NULL,
  `publishEnd` datetime DEFAULT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  `slug` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`newsItemID`)
) ENGINE=InnoDB AUTO_INCREMENT=212 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Note`
--

DROP TABLE IF EXISTS `Note`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Note` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `vehiclePartID` int(11) NOT NULL,
  `note` varchar(255) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_Note_IX` (`vehiclePartID`),
  CONSTRAINT `FK__Note__vehiclePar__2BFEED3A` FOREIGN KEY (`vehiclePartID`) REFERENCES `vcdb_VehiclePart` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=289708 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PackageType`
--

DROP TABLE IF EXISTS `PackageType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PackageType` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `idx` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Part`
--

DROP TABLE IF EXISTS `Part`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Part` (
  `partID` int(11) NOT NULL,
  `status` int(11) NOT NULL,
  `dateModified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `dateAdded` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `shortDesc` varchar(255) DEFAULT NULL,
  `oldPartNumber` varchar(100) DEFAULT NULL,
  `priceCode` int(11) DEFAULT NULL,
  `classID` int(11) NOT NULL DEFAULT '0',
  `featured` tinyint(1) NOT NULL DEFAULT '0',
  `ACESPartTypeID` int(11) DEFAULT NULL,
  PRIMARY KEY (`partID`),
  KEY `IX_Part_status` (`status`),
  KEY `IX_Part_Class` (`classID`),
  KEY `idx` (`partID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PartAttribute`
--

DROP TABLE IF EXISTS `PartAttribute`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PartAttribute` (
  `pAttrID` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `value` varchar(255) DEFAULT NULL,
  `field` varchar(255) DEFAULT NULL,
  `sort` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`pAttrID`),
  KEY `IX_PartAttribute_Part` (`partID`),
  KEY `idx` (`partID`),
  CONSTRAINT `FK__PartAttri__partI__4C541A64` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=65075 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PartGroup`
--

DROP TABLE IF EXISTS `PartGroup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PartGroup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=552 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PartGroupPart`
--

DROP TABLE IF EXISTS `PartGroupPart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PartGroupPart` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `partGroupID` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  `sort` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `FK__PartGroup__partG__2D323D3E` (`partGroupID`),
  KEY `FK__PartGroup__partI__2E266177` (`partID`),
  CONSTRAINT `FK__PartGroup__partG__2D323D3E` FOREIGN KEY (`partGroupID`) REFERENCES `PartGroup` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__PartGroup__partI__2E266177` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=2202 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PartImageSizes`
--

DROP TABLE IF EXISTS `PartImageSizes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PartImageSizes` (
  `sizeID` int(11) NOT NULL AUTO_INCREMENT,
  `size` varchar(25) DEFAULT NULL,
  `dimensions` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`sizeID`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PartImages`
--

DROP TABLE IF EXISTS `PartImages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PartImages` (
  `imageID` int(11) NOT NULL AUTO_INCREMENT,
  `sizeID` int(11) NOT NULL,
  `sort` char(2) NOT NULL,
  `path` varchar(500) NOT NULL,
  `height` int(11) NOT NULL,
  `width` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  PRIMARY KEY (`imageID`),
  KEY `IX_PartImages_Part` (`partID`),
  KEY `IX_PartImages_Size` (`sizeID`),
  KEY `idx` (`partID`),
  CONSTRAINT `FK__PartImage__partI__0E21DDC1` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__PartImage__sizeI__0D2DB988` FOREIGN KEY (`sizeID`) REFERENCES `PartImageSizes` (`sizeID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=498099 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PartIndex`
--

DROP TABLE IF EXISTS `PartIndex`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PartIndex` (
  `partIndexID` bigint(20) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `partIndex` longtext,
  PRIMARY KEY (`partIndexID`)
) ENGINE=InnoDB AUTO_INCREMENT=47779 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PartPackage`
--

DROP TABLE IF EXISTS `PartPackage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PartPackage` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `height` double DEFAULT NULL,
  `width` double DEFAULT NULL,
  `length` double DEFAULT NULL,
  `weight` double DEFAULT NULL,
  `dimensionUOM` int(11) NOT NULL,
  `weightUOM` int(11) NOT NULL,
  `packageUOM` int(11) NOT NULL,
  `quantity` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `idx` (`partID`)
) ENGINE=InnoDB AUTO_INCREMENT=4384 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PartVideo`
--

DROP TABLE IF EXISTS `PartVideo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PartVideo` (
  `pVideoID` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `video` varchar(255) NOT NULL,
  `vTypeID` int(11) NOT NULL,
  `isPrimary` tinyint(1) NOT NULL,
  PRIMARY KEY (`pVideoID`),
  KEY `FK__PartVideo__vType__3723F354` (`vTypeID`),
  KEY `FK__PartVideo__partI__3818178D` (`partID`),
  CONSTRAINT `FK__PartVideo__partI__3818178D` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__PartVideo__vType__3723F354` FOREIGN KEY (`vTypeID`) REFERENCES `videoType` (`vTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=3014 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Post_Category`
--

DROP TABLE IF EXISTS `Post_Category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Post_Category` (
  `postCategoryID` int(11) NOT NULL AUTO_INCREMENT,
  `postID` int(11) NOT NULL,
  `CategoryID` int(11) NOT NULL,
  PRIMARY KEY (`postCategoryID`),
  KEY `FK__Post_Cate__postI__2BC97F7C` (`postID`),
  CONSTRAINT `FK__Post_Cate__postI__2BC97F7C` FOREIGN KEY (`postID`) REFERENCES `Posts` (`postID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Posts`
--

DROP TABLE IF EXISTS `Posts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Posts` (
  `postID` int(11) NOT NULL AUTO_INCREMENT,
  `siteContentID` int(11) DEFAULT NULL,
  `publishedDate` datetime DEFAULT NULL,
  `createdDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `lastModified` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `authorID` int(11) NOT NULL,
  `meta_title` varchar(510) DEFAULT NULL,
  `meta_description` varchar(510) DEFAULT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`postID`),
  KEY `FK__Posts__authorID__28ED12D1` (`authorID`),
  KEY `FK__Posts__siteConte__29E1370A` (`siteContentID`),
  CONSTRAINT `FK__Posts__authorID__28ED12D1` FOREIGN KEY (`authorID`) REFERENCES `Authors` (`authorID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__Posts__siteConte__29E1370A` FOREIGN KEY (`siteContentID`) REFERENCES `SiteContent` (`contentID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Price`
--

DROP TABLE IF EXISTS `Price`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Price` (
  `priceID` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `priceType` varchar(255) DEFAULT NULL,
  `price` decimal(19,4) NOT NULL,
  `enforced` bit(1) NOT NULL,
  PRIMARY KEY (`priceID`),
  KEY `IX_Price_Part` (`partID`),
  KEY `idx` (`partID`),
  CONSTRAINT `FK__Price__partID__0A514CDD` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=29534 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Region`
--

DROP TABLE IF EXISTS `Region`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Region` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `RelatedPart`
--

DROP TABLE IF EXISTS `RelatedPart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `RelatedPart` (
  `relPartID` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `relatedID` bigint(20) NOT NULL,
  `rTypeID` int(11) NOT NULL,
  PRIMARY KEY (`relPartID`),
  KEY `IX_RelatedPart_Part` (`partID`,`relatedID`)
) ENGINE=InnoDB AUTO_INCREMENT=23312 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `RelatedType`
--

DROP TABLE IF EXISTS `RelatedType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `RelatedType` (
  `rTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`rTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ReportType`
--

DROP TABLE IF EXISTS `ReportType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ReportType` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Review`
--

DROP TABLE IF EXISTS `Review`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Review` (
  `reviewID` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) DEFAULT NULL,
  `rating` int(11) NOT NULL,
  `subject` varchar(255) DEFAULT NULL,
  `review_text` longtext,
  `name` varchar(500) DEFAULT NULL,
  `email` varchar(500) DEFAULT NULL,
  `active` tinyint(1) NOT NULL,
  `approved` tinyint(1) NOT NULL,
  `createdDate` datetime NOT NULL,
  `cust_id` int(11) NOT NULL,
  PRIMARY KEY (`reviewID`),
  KEY `ReviewPartID` (`partID`),
  KEY `IX_Review_Part` (`partID`,`createdDate`),
  CONSTRAINT `FK__Review__partID__0C39954F` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=390 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `SalesRepresentative`
--

DROP TABLE IF EXISTS `SalesRepresentative`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `SalesRepresentative` (
  `salesRepID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `code` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`salesRepID`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Services`
--

DROP TABLE IF EXISTS `Services`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Services` (
  `serviceID` int(11) NOT NULL AUTO_INCREMENT,
  `service_title` varchar(255) DEFAULT NULL,
  `description` longtext,
  `service_price` decimal(19,4) NOT NULL DEFAULT '0.0000',
  `hourly` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`serviceID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `SiteContent`
--

DROP TABLE IF EXISTS `SiteContent`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `SiteContent` (
  `contentID` int(11) NOT NULL AUTO_INCREMENT,
  `content_type` varchar(255) DEFAULT NULL,
  `page_title` varchar(500) DEFAULT NULL,
  `createdDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `lastModified` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `meta_title` varchar(255) DEFAULT NULL,
  `meta_description` varchar(255) DEFAULT NULL,
  `keywords` longtext,
  `isPrimary` tinyint(1) NOT NULL DEFAULT '0',
  `published` tinyint(1) NOT NULL DEFAULT '0',
  `active` tinyint(1) NOT NULL DEFAULT '1',
  `slug` varchar(500) DEFAULT NULL,
  `requireAuthentication` tinyint(1) NOT NULL DEFAULT '0',
  `canonical` varchar(255) DEFAULT NULL,
  `websiteID` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`contentID`),
  KEY `FK__SiteConte__websi__21B4D243` (`websiteID`),
  CONSTRAINT `FK__SiteConte__websi__21B4D243` FOREIGN KEY (`websiteID`) REFERENCES `Website` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=66 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `State`
--

DROP TABLE IF EXISTS `State`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `State` (
  `state` varchar(128) NOT NULL,
  `abbr` varchar(128) NOT NULL,
  `stateID` int(11) NOT NULL,
  `countryID` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `States`
--

DROP TABLE IF EXISTS `States`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `States` (
  `stateID` int(11) NOT NULL AUTO_INCREMENT,
  `state` varchar(100) NOT NULL,
  `abbr` varchar(3) NOT NULL,
  `countryID` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`stateID`),
  KEY `FK__States__countryI__607251E5` (`countryID`),
  CONSTRAINT `FK__States__countryI__607251E5` FOREIGN KEY (`countryID`) REFERENCES `Country` (`countryID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=85 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Style`
--

DROP TABLE IF EXISTS `Style`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Style` (
  `styleID` int(11) NOT NULL AUTO_INCREMENT,
  `style` varchar(255) NOT NULL,
  `aaiaID` int(11) NOT NULL,
  PRIMARY KEY (`styleID`)
) ENGINE=InnoDB AUTO_INCREMENT=613 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Submodel`
--

DROP TABLE IF EXISTS `Submodel`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Submodel` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `AAIASubmodelID` int(11) DEFAULT NULL,
  `SubmodelName` varchar(50) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_Submodel_IX` (`AAIASubmodelID`),
  KEY `idx` (`SubmodelName`)
) ENGINE=InnoDB AUTO_INCREMENT=1708 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Testimonial`
--

DROP TABLE IF EXISTS `Testimonial`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Testimonial` (
  `testimonialID` int(11) NOT NULL AUTO_INCREMENT,
  `rating` double NOT NULL,
  `title` varchar(500) DEFAULT NULL,
  `testimonial` longtext,
  `dateAdded` datetime NOT NULL,
  `approved` tinyint(1) NOT NULL,
  `active` tinyint(1) NOT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `location` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`testimonialID`)
) ENGINE=InnoDB AUTO_INCREMENT=88 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Testimonials`
--

DROP TABLE IF EXISTS `Testimonials`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Testimonials` (
  `reviewID` int(11) NOT NULL AUTO_INCREMENT,
  `reviewer` varchar(400) DEFAULT NULL,
  `review` longtext,
  `date_added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_new` int(11) NOT NULL DEFAULT '0',
  `is_hidden` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`reviewID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Trailer`
--

DROP TABLE IF EXISTS `Trailer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Trailer` (
  `trailerID` int(11) NOT NULL AUTO_INCREMENT,
  `image` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `TW` int(11) DEFAULT NULL,
  `GTW` int(11) DEFAULT NULL,
  `hitchClass` varchar(255) DEFAULT NULL,
  `shortDesc` varchar(1000) DEFAULT NULL,
  `message` varchar(1000) DEFAULT NULL,
  PRIMARY KEY (`trailerID`)
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `UnitOfMeasure`
--

DROP TABLE IF EXISTS `UnitOfMeasure`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `UnitOfMeasure` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `code` varchar(5) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `UserProfiles`
--

DROP TABLE IF EXISTS `UserProfiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `UserProfiles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `customerID` int(11) DEFAULT NULL,
  `custID` int(11) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `IP` varchar(25) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=527 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Vehicle`
--

DROP TABLE IF EXISTS `Vehicle`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Vehicle` (
  `vehicleID` int(11) NOT NULL AUTO_INCREMENT,
  `yearID` int(11) NOT NULL,
  `makeID` int(11) NOT NULL,
  `modelID` int(11) NOT NULL,
  `styleID` int(11) NOT NULL,
  `dateAdded` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`vehicleID`),
  KEY `IX_Vehicle_Filter` (`yearID`,`makeID`,`modelID`,`styleID`),
  KEY `FK__Vehicle__makeID__4E3C62D6` (`makeID`),
  KEY `FK__Vehicle__modelID__4F30870F` (`modelID`),
  KEY `FK__Vehicle__styleID__5024AB48` (`styleID`),
  CONSTRAINT `FK__Vehicle__makeID__4E3C62D6` FOREIGN KEY (`makeID`) REFERENCES `Make` (`makeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__Vehicle__modelID__4F30870F` FOREIGN KEY (`modelID`) REFERENCES `Model` (`modelID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__Vehicle__styleID__5024AB48` FOREIGN KEY (`styleID`) REFERENCES `Style` (`styleID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__Vehicle__yearID__4D483E9D` FOREIGN KEY (`yearID`) REFERENCES `Year` (`yearID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=248482 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `VehicleConfig`
--

DROP TABLE IF EXISTS `VehicleConfig`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `VehicleConfig` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `AAIAVehicleConfigID` int(11) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_VehicleConfig_IX` (`AAIAVehicleConfigID`)
) ENGINE=InnoDB AUTO_INCREMENT=28002 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `VehicleConfigAttribute`
--

DROP TABLE IF EXISTS `VehicleConfigAttribute`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `VehicleConfigAttribute` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `AttributeID` int(11) NOT NULL,
  `VehicleConfigID` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `FK__VehicleCo__Attri__19E03CFF` (`AttributeID`),
  KEY `FK__VehicleCo__Vehic__1AD46138` (`VehicleConfigID`),
  CONSTRAINT `FK__VehicleCo__Attri__19E03CFF` FOREIGN KEY (`AttributeID`) REFERENCES `ConfigAttribute` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__VehicleCo__Vehic__1AD46138` FOREIGN KEY (`VehicleConfigID`) REFERENCES `VehicleConfig` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=38264 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `VehiclePart`
--

DROP TABLE IF EXISTS `VehiclePart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `VehiclePart` (
  `vPartID` int(11) NOT NULL AUTO_INCREMENT,
  `vehicleID` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  `drilling` varchar(100) DEFAULT NULL,
  `exposed` varchar(100) DEFAULT NULL,
  `installTime` int(11) DEFAULT NULL,
  PRIMARY KEY (`vPartID`),
  KEY `IX_VehiclePart_Part` (`vehicleID`,`partID`),
  KEY `FK__VehiclePa__partI__0F1601FA` (`partID`),
  CONSTRAINT `FK__VehiclePa__partI__0F1601FA` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__VehiclePa__vehic__5118CF81` FOREIGN KEY (`vehicleID`) REFERENCES `Vehicle` (`vehicleID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=39419 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `VehiclePartAttribute`
--

DROP TABLE IF EXISTS `VehiclePartAttribute`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `VehiclePartAttribute` (
  `vpAttrID` int(11) NOT NULL AUTO_INCREMENT,
  `vPartID` int(11) NOT NULL,
  `value` varchar(255) DEFAULT NULL,
  `field` varchar(255) DEFAULT NULL,
  `sort` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`vpAttrID`),
  KEY `IX_VehiclePartAttr_VPart` (`vPartID`),
  CONSTRAINT `FK__VehiclePa__vPart__520CF3BA` FOREIGN KEY (`vPartID`) REFERENCES `VehiclePart` (`vPartID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=98331 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `VehicleType`
--

DROP TABLE IF EXISTS `VehicleType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `VehicleType` (
  `VehicleTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `VehicleTypeName` varchar(50) NOT NULL,
  `VehicleTypeGroupID` int(11) DEFAULT NULL,
  PRIMARY KEY (`VehicleTypeID`),
  KEY `FK__VehicleTy__Vehic__648AFD1B` (`VehicleTypeGroupID`),
  CONSTRAINT `FK__VehicleTy__Vehic__648AFD1B` FOREIGN KEY (`VehicleTypeGroupID`) REFERENCES `VehicleTypeGroup` (`VehicleTypeGroupID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `VehicleTypeGroup`
--

DROP TABLE IF EXISTS `VehicleTypeGroup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `VehicleTypeGroup` (
  `VehicleTypeGroupID` int(11) NOT NULL AUTO_INCREMENT,
  `VehicleTypeGroupName` varchar(50) NOT NULL,
  PRIMARY KEY (`VehicleTypeGroupID`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Video`
--

DROP TABLE IF EXISTS `Video`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Video` (
  `videoID` int(11) NOT NULL AUTO_INCREMENT,
  `embed_link` varchar(200) DEFAULT NULL,
  `dateAdded` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `sort` int(11) NOT NULL DEFAULT '1',
  `title` varchar(255) DEFAULT NULL,
  `description` longtext,
  `youtubeID` varchar(255) DEFAULT NULL,
  `watchpage` varchar(255) DEFAULT NULL,
  `screenshot` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`videoID`)
) ENGINE=InnoDB AUTO_INCREMENT=59 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `WebProperties`
--

DROP TABLE IF EXISTS `WebProperties`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `WebProperties` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(75) NOT NULL,
  `cust_ID` int(11) NOT NULL,
  `badgeID` varchar(64) NOT NULL,
  `url` longtext,
  `isEnabled` tinyint(1) NOT NULL,
  `sellerID` varchar(50) DEFAULT NULL,
  `typeID` int(11) DEFAULT NULL,
  `isFinalApproved` tinyint(1) NOT NULL,
  `isEnabledDate` datetime DEFAULT NULL,
  `isDenied` tinyint(1) NOT NULL,
  `requestedDate` datetime DEFAULT NULL,
  `addedDate` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `badgeID` (`badgeID`)
) ENGINE=InnoDB AUTO_INCREMENT=165 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `WebPropertyTypes`
--

DROP TABLE IF EXISTS `WebPropertyTypes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `WebPropertyTypes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `typeID` int(11) NOT NULL,
  `type` varchar(50) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Website`
--

DROP TABLE IF EXISTS `Website`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Website` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `WidgetDeployments`
--

DROP TABLE IF EXISTS `WidgetDeployments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `WidgetDeployments` (
  `trackerID` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(400) NOT NULL,
  `date_added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`trackerID`)
) ENGINE=InnoDB AUTO_INCREMENT=76 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Year`
--

DROP TABLE IF EXISTS `Year`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Year` (
  `yearID` int(11) NOT NULL AUTO_INCREMENT,
  `year` double NOT NULL,
  PRIMARY KEY (`yearID`)
) ENGINE=InnoDB AUTO_INCREMENT=283 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `YearMake`
--

DROP TABLE IF EXISTS `YearMake`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `YearMake` (
  `ymID` int(11) NOT NULL AUTO_INCREMENT,
  `yearID` int(11) DEFAULT NULL,
  `makeID` int(11) NOT NULL,
  PRIMARY KEY (`ymID`),
  KEY `IX_YearMake` (`yearID`,`makeID`),
  KEY `FK__YearMake__makeID__478F6547` (`makeID`),
  CONSTRAINT `FK__YearMake__makeID__478F6547` FOREIGN KEY (`makeID`) REFERENCES `Make` (`makeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__YearMake__yearID__469B410E` FOREIGN KEY (`yearID`) REFERENCES `Year` (`yearID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1301 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `vcdb_Make`
--

DROP TABLE IF EXISTS `vcdb_Make`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vcdb_Make` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `AAIAMakeID` int(11) DEFAULT NULL,
  `MakeName` varchar(50) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_Make_IX` (`AAIAMakeID`),
  KEY `idx` (`MakeName`)
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `vcdb_Model`
--

DROP TABLE IF EXISTS `vcdb_Model`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vcdb_Model` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `AAIAModelID` int(11) DEFAULT NULL,
  `ModelName` varchar(100) DEFAULT NULL,
  `VehicleTypeID` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_Model_IX` (`AAIAModelID`),
  KEY `idx` (`ModelName`)
) ENGINE=InnoDB AUTO_INCREMENT=3868 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `vcdb_Vehicle`
--

DROP TABLE IF EXISTS `vcdb_Vehicle`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vcdb_Vehicle` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `BaseVehicleID` int(11) NOT NULL,
  `SubModelID` int(11) DEFAULT NULL,
  `ConfigID` int(11) DEFAULT NULL,
  `AppID` int(11) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_Vehicle_IX` (`BaseVehicleID`,`SubModelID`,`ConfigID`),
  KEY `FK__vcdb_Vehi__SubMo__208D3A8E` (`SubModelID`),
  KEY `FK__vcdb_Vehi__Confi__21815EC7` (`ConfigID`),
  CONSTRAINT `FK__vcdb_Vehi__BaseV__1F991655` FOREIGN KEY (`BaseVehicleID`) REFERENCES `BaseVehicle` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__vcdb_Vehi__Confi__21815EC7` FOREIGN KEY (`ConfigID`) REFERENCES `VehicleConfig` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__vcdb_Vehi__SubMo__208D3A8E` FOREIGN KEY (`SubModelID`) REFERENCES `Submodel` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=33906 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `vcdb_VehiclePart`
--

DROP TABLE IF EXISTS `vcdb_VehiclePart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vcdb_VehiclePart` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `VehicleID` int(11) NOT NULL,
  `PartNumber` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_VehiclePart_Part_IX` (`VehicleID`,`PartNumber`),
  KEY `FK__vcdb_Vehi__PartN__273A381D` (`PartNumber`),
  CONSTRAINT `FK__vcdb_Vehi__PartN__273A381D` FOREIGN KEY (`PartNumber`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__vcdb_Vehi__Vehic__264613E4` FOREIGN KEY (`VehicleID`) REFERENCES `vcdb_Vehicle` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=191945 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `vcdb_Year`
--

DROP TABLE IF EXISTS `vcdb_Year`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vcdb_Year` (
  `YearID` int(11) NOT NULL,
  PRIMARY KEY (`YearID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `videoType`
--

DROP TABLE IF EXISTS `videoType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `videoType` (
  `vTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `icon` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`vTypeID`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2013-04-10 16:02:49
