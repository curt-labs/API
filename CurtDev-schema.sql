/*
 Navicat Premium Data Transfer

 Source Server         : Localhost
 Source Server Type    : MySQL
 Source Server Version : 50614
 Source Host           : localhost
 Source Database       : CurtDev

 Target Server Type    : MySQL
 Target Server Version : 50614
 File Encoding         : utf-8

 Date: 12/03/2013 09:06:27 AM
*/

CREATE DATABASE IF NOT EXISTS CurtDev;

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `AcesType`
-- ----------------------------
DROP TABLE IF EXISTS `AcesType`;
CREATE TABLE `AcesType` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ApiAccess`
-- ----------------------------
DROP TABLE IF EXISTS `ApiAccess`;
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

-- ----------------------------
--  Table structure for `ApiKey`
-- ----------------------------
DROP TABLE IF EXISTS `ApiKey`;
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

-- ----------------------------
--  Table structure for `ApiKeyType`
-- ----------------------------
DROP TABLE IF EXISTS `ApiKeyType`;
CREATE TABLE `ApiKeyType` (
  `id` varchar(64) NOT NULL,
  `type` varchar(500) DEFAULT NULL,
  `date_added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ApiModules`
-- ----------------------------
DROP TABLE IF EXISTS `ApiModules`;
CREATE TABLE `ApiModules` (
  `id` varchar(64) NOT NULL,
  `name` varchar(500) DEFAULT NULL,
  `access_level` varchar(64) DEFAULT NULL,
  `date_added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `AuthAccess`
-- ----------------------------
DROP TABLE IF EXISTS `AuthAccess`;
CREATE TABLE `AuthAccess` (
  `id` varchar(64) NOT NULL,
  `userID` varchar(64) NOT NULL,
  `AreaID` varchar(64) NOT NULL,
  `dateAdded` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  KEY `authArea_ref_idx` (`AreaID`),
  KEY `custUser_ref_idx` (`userID`),
  CONSTRAINT `custUserAuthAccess_ref` FOREIGN KEY (`userID`) REFERENCES `CustomerUser` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `AuthAreas`
-- ----------------------------
DROP TABLE IF EXISTS `AuthAreas`;
CREATE TABLE `AuthAreas` (
  `id` varchar(64) NOT NULL,
  `path` varchar(50) NOT NULL,
  `DomainID` varchar(64) NOT NULL,
  `name` varchar(50) NOT NULL,
  `parentAreaID` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  KEY `authDomain_ref_idx` (`DomainID`),
  CONSTRAINT `authDomains` FOREIGN KEY (`DomainID`) REFERENCES `AuthDomains` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `AuthDomains`
-- ----------------------------
DROP TABLE IF EXISTS `AuthDomains`;
CREATE TABLE `AuthDomains` (
  `id` varchar(64) NOT NULL,
  `url` varchar(50) NOT NULL,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `AuthorizedTracking`
-- ----------------------------
DROP TABLE IF EXISTS `AuthorizedTracking`;
CREATE TABLE `AuthorizedTracking` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `property` varchar(500) NOT NULL,
  `view_count` int(11) NOT NULL,
  `authorized_id` varchar(500) DEFAULT NULL,
  `date_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `date_modified` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Authors`
-- ----------------------------
DROP TABLE IF EXISTS `Authors`;
CREATE TABLE `Authors` (
  `authorID` int(11) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(510) DEFAULT NULL,
  `last_name` varchar(510) DEFAULT NULL,
  `email` varchar(510) DEFAULT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`authorID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Banners`
-- ----------------------------
DROP TABLE IF EXISTS `Banners`;
CREATE TABLE `Banners` (
  `bannerID` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) DEFAULT NULL,
  `link` varchar(400) DEFAULT NULL,
  `starts` datetime DEFAULT NULL,
  `ends` datetime DEFAULT NULL,
  `path` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`bannerID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `BaseVehicle`
-- ----------------------------
DROP TABLE IF EXISTS `BaseVehicle`;
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
) ENGINE=InnoDB AUTO_INCREMENT=22513 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `BlogCategories`
-- ----------------------------
DROP TABLE IF EXISTS `BlogCategories`;
CREATE TABLE `BlogCategories` (
  `blogCategoryID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `slug` varchar(255) DEFAULT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`blogCategoryID`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `BlogPost_BlogCategory`
-- ----------------------------
DROP TABLE IF EXISTS `BlogPost_BlogCategory`;
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

-- ----------------------------
--  Table structure for `BlogPosts`
-- ----------------------------
DROP TABLE IF EXISTS `BlogPosts`;
CREATE TABLE `BlogPosts` (
  `blogPostID` int(11) NOT NULL AUTO_INCREMENT,
  `post_title` varchar(500) NOT NULL,
  `slug` varchar(500) NOT NULL,
  `post_text` longtext,
  `publishedDate` datetime DEFAULT NULL,
  `createdDate` datetime NOT NULL,
  `lastModified` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `userID` int(11) NOT NULL,
  `meta_title` varchar(510) DEFAULT NULL,
  `meta_description` varchar(510) DEFAULT NULL,
  `keywords` longtext,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`blogPostID`),
  KEY `BlogPostAuthorID` (`userID`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Brand`
-- ----------------------------
DROP TABLE IF EXISTS `Brand`;
CREATE TABLE `Brand` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `code` varchar(255) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `BrandPart`
-- ----------------------------
DROP TABLE IF EXISTS `BrandPart`;
CREATE TABLE `BrandPart` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `brandID` int(11) NOT NULL,
  `brandPartID` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  `interchangeType` char(1) NOT NULL,
  `dateAdded` datetime NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `BusinessClass`
-- ----------------------------
DROP TABLE IF EXISTS `BusinessClass`;
CREATE TABLE `BusinessClass` (
  `BusinessClassID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `sort` int(11) NOT NULL,
  `showOnWebsite` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`BusinessClassID`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Cabelas`
-- ----------------------------
DROP TABLE IF EXISTS `Cabelas`;
CREATE TABLE `Cabelas` (
  `cabelasID` int(11) NOT NULL AUTO_INCREMENT,
  `priceCode` int(11) DEFAULT NULL,
  `cabelasPart` varchar(50) NOT NULL,
  PRIMARY KEY (`cabelasID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `CartIntegration`
-- ----------------------------
DROP TABLE IF EXISTS `CartIntegration`;
CREATE TABLE `CartIntegration` (
  `referenceID` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `custPartID` int(11) NOT NULL,
  `custID` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`referenceID`),
  KEY `partID` (`partID`)
) ENGINE=InnoDB AUTO_INCREMENT=69508 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `CatPart`
-- ----------------------------
DROP TABLE IF EXISTS `CatPart`;
CREATE TABLE `CatPart` (
  `catPartID` int(11) NOT NULL AUTO_INCREMENT,
  `catID` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  PRIMARY KEY (`catPartID`),
  KEY `IX_CatPart_Cat_Part` (`catID`,`partID`),
  KEY `FK__CatPart__partID__54945AAA` (`partID`),
  CONSTRAINT `FK__CatPart__catID__55887EE3` FOREIGN KEY (`catID`) REFERENCES `Categories` (`catID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__CatPart__partID__54945AAA` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=5228 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Categories`
-- ----------------------------
DROP TABLE IF EXISTS `Categories`;
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
  KEY `IX_Categories_Sort` (`sort`)
) ENGINE=InnoDB AUTO_INCREMENT=282 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Class`
-- ----------------------------
DROP TABLE IF EXISTS `Class`;
CREATE TABLE `Class` (
  `classID` int(11) NOT NULL AUTO_INCREMENT,
  `class` varchar(255) DEFAULT NULL,
  `image` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`classID`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ColorCode`
-- ----------------------------
DROP TABLE IF EXISTS `ColorCode`;
CREATE TABLE `ColorCode` (
  `codeID` int(11) NOT NULL,
  `code` varchar(100) DEFAULT NULL,
  `font` varchar(100) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Comments`
-- ----------------------------
DROP TABLE IF EXISTS `Comments`;
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

-- ----------------------------
--  Table structure for `Company`
-- ----------------------------
DROP TABLE IF EXISTS `Company`;
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

-- ----------------------------
--  Table structure for `ConfigAttribute`
-- ----------------------------
DROP TABLE IF EXISTS `ConfigAttribute`;
CREATE TABLE `ConfigAttribute` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `ConfigAttributeTypeID` int(11) NOT NULL,
  `parentID` int(11) NOT NULL,
  `vcdbID` int(11) DEFAULT NULL,
  `value` varchar(255) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_ConfigAttribute_IX` (`ConfigAttributeTypeID`,`parentID`),
  CONSTRAINT `FK__ConfigAtt__Confi__07D43958` FOREIGN KEY (`ConfigAttributeTypeID`) REFERENCES `ConfigAttributeType` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=308 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ConfigAttributeType`
-- ----------------------------
DROP TABLE IF EXISTS `ConfigAttributeType`;
CREATE TABLE `ConfigAttributeType` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `AcesTypeID` int(11) DEFAULT NULL,
  `sort` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `FK__ConfigAtt__AcesT__030F843B` (`AcesTypeID`),
  CONSTRAINT `FK__ConfigAtt__AcesT__030F843B` FOREIGN KEY (`AcesTypeID`) REFERENCES `AcesType` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=77 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Contact`
-- ----------------------------
DROP TABLE IF EXISTS `Contact`;
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
) ENGINE=InnoDB AUTO_INCREMENT=9610 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ContactReceiver`
-- ----------------------------
DROP TABLE IF EXISTS `ContactReceiver`;
CREATE TABLE `ContactReceiver` (
  `contactReceiverID` int(11) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`contactReceiverID`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ContactReceiver_ContactType`
-- ----------------------------
DROP TABLE IF EXISTS `ContactReceiver_ContactType`;
CREATE TABLE `ContactReceiver_ContactType` (
  `receiverTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `contactReceiverID` int(11) NOT NULL,
  `contactTypeID` int(11) NOT NULL,
  PRIMARY KEY (`receiverTypeID`),
  KEY `FK__ContactRe__conta__6FB49575` (`contactReceiverID`),
  KEY `FK__ContactRe__conta__70A8B9AE` (`contactTypeID`),
  CONSTRAINT `FK__ContactRe__conta__6FB49575` FOREIGN KEY (`contactReceiverID`) REFERENCES `ContactReceiver` (`contactReceiverID`) ON DELETE CASCADE,
  CONSTRAINT `FK__ContactRe__conta__70A8B9AE` FOREIGN KEY (`contactTypeID`) REFERENCES `ContactType` (`contactTypeID`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=113 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ContactType`
-- ----------------------------
DROP TABLE IF EXISTS `ContactType`;
CREATE TABLE `ContactType` (
  `contactTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`contactTypeID`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Content`
-- ----------------------------
DROP TABLE IF EXISTS `Content`;
CREATE TABLE `Content` (
  `contentID` int(11) NOT NULL AUTO_INCREMENT,
  `text` longtext,
  `cTypeID` int(11) NOT NULL,
  PRIMARY KEY (`contentID`),
  KEY `FK__Content__cTypeID__0B457116` (`cTypeID`),
  CONSTRAINT `FK__Content__cTypeID__0B457116` FOREIGN KEY (`cTypeID`) REFERENCES `ContentType` (`cTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=298062 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ContentBridge`
-- ----------------------------
DROP TABLE IF EXISTS `ContentBridge`;
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
) ENGINE=InnoDB AUTO_INCREMENT=28708 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ContentType`
-- ----------------------------
DROP TABLE IF EXISTS `ContentType`;
CREATE TABLE `ContentType` (
  `cTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(255) DEFAULT NULL,
  `allowHTML` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`cTypeID`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Country`
-- ----------------------------
DROP TABLE IF EXISTS `Country`;
CREATE TABLE `Country` (
  `countryID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `abbr` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`countryID`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `CustUserWebProperties`
-- ----------------------------
DROP TABLE IF EXISTS `CustUserWebProperties`;
CREATE TABLE `CustUserWebProperties` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userID` varchar(64) NOT NULL,
  `webPropID` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `webProp_ref_idx` (`webPropID`),
  KEY `custUser_ref_idx` (`userID`),
  CONSTRAINT `custUser_ref` FOREIGN KEY (`userID`) REFERENCES `CustomerUser` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `webProps_ref` FOREIGN KEY (`webPropID`) REFERENCES `WebProperties` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=193 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Customer`
-- ----------------------------
DROP TABLE IF EXISTS `Customer`;
CREATE TABLE `Customer` (
  `cust_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `address` varchar(500) DEFAULT NULL,
  `city` varchar(150) DEFAULT NULL,
  `stateID` int(11) DEFAULT NULL,
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
) ENGINE=InnoDB AUTO_INCREMENT=10443411 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `CustomerContent`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerContent`;
CREATE TABLE `CustomerContent` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `text` longtext NOT NULL,
  `custID` int(11) NOT NULL,
  `added` datetime NOT NULL,
  `modified` datetime NOT NULL,
  `userID` varchar(64) NOT NULL,
  `typeID` int(11) NOT NULL,
  `deleted` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `cTypeID_idx` (`typeID`),
  KEY `cust_id_idx` (`custID`),
  KEY `id_idx` (`userID`),
  KEY `deleted_idx` (`deleted`),
  CONSTRAINT `cTypeID` FOREIGN KEY (`typeID`) REFERENCES `ContentType` (`cTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `cust_id` FOREIGN KEY (`custID`) REFERENCES `Customer` (`cust_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `id` FOREIGN KEY (`userID`) REFERENCES `CustomerUser` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COMMENT='Customer Content';
delimiter ;;
CREATE TRIGGER `CustomerContent_INSERT` AFTER INSERT ON `CustomerContent` FOR EACH ROW BEGIN IF NEW.deleted THEN SET @changeType = 'DELETE'; ELSE SET @changeType ='NEW'; END IF; INSERT INTO CustomerContent_Revisions (userID,custID,new_text,date,changeType, contentID,new_type) VALUES (NEW.userID,NEW.custID,NEW.text,CURRENT_TIMESTAMP, @changeType, NEW.id, NEW.typeID); END;
 ;;
delimiter ;
delimiter ;;
CREATE TRIGGER `CustomerContent_Update` AFTER UPDATE ON `CustomerContent` FOR EACH ROW BEGIN IF NEW.deleted THEN SET @changeType = 'DELETE'; ELSE SET @changeType = 'EDIT'; END IF; INSERT INTO CustomerContent_Revisions (userID,custID,new_text,old_text, date, changeType,contentID,new_type,old_type) VALUES(NEW.userID,NEW.custID,NEW.text,OLD.text,NOW(),@changeType, NEW.id,NEW.typeID,OLD.typeID); END;
 ;;
delimiter ;

-- ----------------------------
--  Table structure for `CustomerContentBridge`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerContentBridge`;
CREATE TABLE `CustomerContentBridge` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `catID` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  `contentID` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `id_idx` (`contentID`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COMMENT='Customer Content Bridge';

-- ----------------------------
--  Table structure for `CustomerContent_Revisions`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerContent_Revisions`;
CREATE TABLE `CustomerContent_Revisions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userID` varchar(64) NOT NULL,
  `custID` int(11) NOT NULL,
  `old_text` longtext,
  `new_text` longtext,
  `date` datetime NOT NULL,
  `changeType` enum('NEW','EDIT','DELETE') NOT NULL,
  `contentID` int(11) NOT NULL,
  `old_type` int(11) NOT NULL,
  `new_type` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `id_idx` (`contentID`),
  KEY `cTypeID_idx` (`old_type`),
  KEY `cTypeID_idx1` (`new_type`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8 COMMENT='Revision History for CustomerContent';

-- ----------------------------
--  Table structure for `CustomerCost`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerCost`;
CREATE TABLE `CustomerCost` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `cust_id` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  `cost` decimal(18,2) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=4599 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `CustomerLocations`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerLocations`;
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
) ENGINE=InnoDB AUTO_INCREMENT=7820 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `CustomerPartAttributeFields`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerPartAttributeFields`;
CREATE TABLE `CustomerPartAttributeFields` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `field` varchar(255) DEFAULT NULL,
  `dataType` int(11) NOT NULL,
  `added` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `id_idx` (`dataType`),
  CONSTRAINT `datatype_id` FOREIGN KEY (`dataType`) REFERENCES `DataTypes` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `CustomerPartAttributeValues`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerPartAttributeValues`;
CREATE TABLE `CustomerPartAttributeValues` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `value` varchar(255) NOT NULL,
  `added` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `CustomerPartAttributes`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerPartAttributes`;
CREATE TABLE `CustomerPartAttributes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fieldID` int(11) NOT NULL,
  `valueID` int(11) NOT NULL,
  `sort` int(11) NOT NULL,
  `custID` int(11) NOT NULL,
  `userID` varchar(64) NOT NULL,
  `added` datetime NOT NULL,
  `modified` datetime NOT NULL,
  `deleted` tinyint(1) unsigned zerofill NOT NULL,
  `partID` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `id_idx` (`fieldID`),
  KEY `id_idx1` (`valueID`),
  KEY `cust_id_idx` (`custID`),
  KEY `id_idx2` (`userID`),
  KEY `partID_idx` (`partID`),
  CONSTRAINT `CustomerPartAttributeFields_id` FOREIGN KEY (`fieldID`) REFERENCES `CustomerPartAttributeFields` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `CustomerPartAttributeValues_id` FOREIGN KEY (`valueID`) REFERENCES `CustomerPartAttributeValues` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `CustomerUser_id` FOREIGN KEY (`userID`) REFERENCES `CustomerUser` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `Customer_cust_id` FOREIGN KEY (`custID`) REFERENCES `Customer` (`cust_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `Part_partID` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
delimiter ;;
CREATE TRIGGER `CustomerPartAttributes_Insert` AFTER INSERT ON `CustomerPartAttributes` FOR EACH ROW BEGIN 
	IF NEW.deleted THEN 
		SET @changeType = 'DELETE'; 
	ELSE 
		SET @changeType ='NEW'; 
	END IF; 
	INSERT INTO CustomerPartAttributes_Revisions
		(userID,custID,new_field,new_value,changeType,attributeID)
	VALUES
		(NEW.userID,NEW.custID, NEW.fieldID, NEW.valueID, @changeType, NEW.id);
END;
 ;;
delimiter ;
delimiter ;;
CREATE TRIGGER `CustomerPartAttributes_Update` AFTER UPDATE ON `CustomerPartAttributes` FOR EACH ROW BEGIN
	IF NEW.deleted THEN
		SET @changeType = 'DELETE';
	ELSE
		SET @changeType = 'EDIT';
	END IF;
	INSERT INTO CustomerPartAttributes_Revisions 
		(userID,custID,old_field,new_field,old_value,new_value,changeType,attributeID)
	VALUES
		(NEW.userID,NEW.custID,OLD.fieldID, NEW.fieldID, OLD.valueID, NEW.valueID, @changeType, NEW.id);
END;
 ;;
delimiter ;

-- ----------------------------
--  Table structure for `CustomerPartAttributes_Revisions`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerPartAttributes_Revisions`;
CREATE TABLE `CustomerPartAttributes_Revisions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userID` varchar(64) NOT NULL,
  `custID` int(11) DEFAULT NULL,
  `old_field` int(11) NOT NULL,
  `new_field` int(11) NOT NULL,
  `old_value` int(11) NOT NULL,
  `new_value` int(11) NOT NULL,
  `date` datetime NOT NULL,
  `changeType` enum('NEW','EDIT','DELETE') NOT NULL,
  `attributeID` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `id_idx` (`userID`),
  KEY `cust_id_idx` (`custID`),
  KEY `id_idx1` (`old_field`),
  KEY `id_idx2` (`new_field`),
  KEY `id_idx3` (`old_value`),
  KEY `id_idx4` (`new_value`),
  KEY `id_idx5` (`attributeID`),
  CONSTRAINT `CustomerPartAttributeFields_Revisions_new_id` FOREIGN KEY (`new_field`) REFERENCES `CustomerPartAttributeFields` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `CustomerPartAttributeFields_Revisions_old_id` FOREIGN KEY (`old_field`) REFERENCES `CustomerPartAttributeFields` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `CustomerPartAttributes_Revisions_id` FOREIGN KEY (`attributeID`) REFERENCES `CustomerPartAttributes` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `CustomerPartAttributeValues_Revisions_new_id` FOREIGN KEY (`new_value`) REFERENCES `CustomerPartAttributeValues` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `CustomerPartAttributeValues_Revisions_old_id` FOREIGN KEY (`old_value`) REFERENCES `CustomerPartAttributeValues` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `CustomerUser_Revisions_id` FOREIGN KEY (`userID`) REFERENCES `CustomerUser` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `Customer_Revisions_cust_id` FOREIGN KEY (`custID`) REFERENCES `Customer` (`cust_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `CustomerPricing`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerPricing`;
CREATE TABLE `CustomerPricing` (
  `cust_price_id` int(11) NOT NULL AUTO_INCREMENT,
  `cust_id` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  `price` decimal(8,2) DEFAULT NULL,
  `isSale` int(11) NOT NULL DEFAULT '0',
  `sale_start` date DEFAULT NULL,
  `sale_end` date DEFAULT NULL,
  PRIMARY KEY (`cust_price_id`),
  KEY `partID` (`partID`)
) ENGINE=InnoDB AUTO_INCREMENT=388362 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `CustomerReport`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerReport`;
CREATE TABLE `CustomerReport` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `customerID` int(11) NOT NULL,
  `created` datetime NOT NULL,
  `ReportTypeID` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `FK__CustomerR__Repor__0F604C87` (`ReportTypeID`),
  CONSTRAINT `FK__CustomerR__Repor__0F604C87` FOREIGN KEY (`ReportTypeID`) REFERENCES `ReportType` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `CustomerReportPart`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerReportPart`;
CREATE TABLE `CustomerReportPart` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `customerID` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=279 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `CustomerUser`
-- ----------------------------
DROP TABLE IF EXISTS `CustomerUser`;
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

-- ----------------------------
--  Table structure for `DataTypes`
-- ----------------------------
DROP TABLE IF EXISTS `DataTypes`;
CREATE TABLE `DataTypes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(45) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `DealerTiers`
-- ----------------------------
DROP TABLE IF EXISTS `DealerTiers`;
CREATE TABLE `DealerTiers` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `tier` varchar(255) DEFAULT NULL,
  `sort` int(11) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `DealerTypes`
-- ----------------------------
DROP TABLE IF EXISTS `DealerTypes`;
CREATE TABLE `DealerTypes` (
  `dealer_type` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(100) DEFAULT NULL,
  `online` tinyint(1) NOT NULL DEFAULT '0',
  `show` tinyint(1) NOT NULL DEFAULT '1',
  `label` varchar(300) DEFAULT NULL,
  PRIMARY KEY (`dealer_type`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `FAQ`
-- ----------------------------
DROP TABLE IF EXISTS `FAQ`;
CREATE TABLE `FAQ` (
  `faqID` int(11) NOT NULL AUTO_INCREMENT,
  `question` varchar(500) DEFAULT NULL,
  `answer` longtext,
  PRIMARY KEY (`faqID`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `File`
-- ----------------------------
DROP TABLE IF EXISTS `File`;
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
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `FileExt`
-- ----------------------------
DROP TABLE IF EXISTS `FileExt`;
CREATE TABLE `FileExt` (
  `fileExtID` int(11) NOT NULL AUTO_INCREMENT,
  `fileExt` varchar(10) NOT NULL,
  `fileExtIcon` varchar(1000) DEFAULT NULL,
  `fileTypeID` int(11) NOT NULL,
  PRIMARY KEY (`fileExtID`),
  KEY `FK__FileExt__fileTyp__6A50C1DA` (`fileTypeID`),
  CONSTRAINT `FK__FileExt__fileTyp__6A50C1DA` FOREIGN KEY (`fileTypeID`) REFERENCES `FileType` (`fileTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `FileGallery`
-- ----------------------------
DROP TABLE IF EXISTS `FileGallery`;
CREATE TABLE `FileGallery` (
  `fileGalleryID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` varchar(4000) DEFAULT NULL,
  `parentID` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`fileGalleryID`)
) ENGINE=InnoDB AUTO_INCREMENT=129 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `FileType`
-- ----------------------------
DROP TABLE IF EXISTS `FileType`;
CREATE TABLE `FileType` (
  `fileTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `fileType` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`fileTypeID`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ForumGroup`
-- ----------------------------
DROP TABLE IF EXISTS `ForumGroup`;
CREATE TABLE `ForumGroup` (
  `forumGroupID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` longtext,
  `createdDate` datetime NOT NULL,
  PRIMARY KEY (`forumGroupID`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ForumPost`
-- ----------------------------
DROP TABLE IF EXISTS `ForumPost`;
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

-- ----------------------------
--  Table structure for `ForumThread`
-- ----------------------------
DROP TABLE IF EXISTS `ForumThread`;
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

-- ----------------------------
--  Table structure for `ForumTopic`
-- ----------------------------
DROP TABLE IF EXISTS `ForumTopic`;
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

-- ----------------------------
--  Table structure for `Gallery`
-- ----------------------------
DROP TABLE IF EXISTS `Gallery`;
CREATE TABLE `Gallery` (
  `imgID` int(11) NOT NULL AUTO_INCREMENT,
  `img_path` varchar(500) NOT NULL,
  `title` varchar(200) DEFAULT NULL,
  `sort_order` int(11) NOT NULL,
  PRIMARY KEY (`imgID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `IPBlock`
-- ----------------------------
DROP TABLE IF EXISTS `IPBlock`;
CREATE TABLE `IPBlock` (
  `blockID` int(11) NOT NULL AUTO_INCREMENT,
  `IPAddress` varchar(255) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `createdDate` datetime NOT NULL,
  `userID` int(11) NOT NULL,
  PRIMARY KEY (`blockID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `IncludedPart`
-- ----------------------------
DROP TABLE IF EXISTS `IncludedPart`;
CREATE TABLE `IncludedPart` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `includedID` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `FK__IncludedP__partI__33DF3ACD` (`partID`),
  KEY `FK__IncludedP__inclu__34D35F06` (`includedID`),
  CONSTRAINT `FK__IncludedP__inclu__34D35F06` FOREIGN KEY (`includedID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__IncludedP__partI__33DF3ACD` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `KioskOrderItems`
-- ----------------------------
DROP TABLE IF EXISTS `KioskOrderItems`;
CREATE TABLE `KioskOrderItems` (
  `itemID` int(11) NOT NULL AUTO_INCREMENT,
  `orderID` int(11) NOT NULL,
  `partID` int(11) NOT NULL,
  `quantity` int(11) NOT NULL,
  `price` decimal(19,4) NOT NULL,
  `isFulfilled` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`itemID`)
) ENGINE=InnoDB AUTO_INCREMENT=124 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `KioskOrders`
-- ----------------------------
DROP TABLE IF EXISTS `KioskOrders`;
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
) ENGINE=InnoDB AUTO_INCREMENT=10000081 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `LandingPage`
-- ----------------------------
DROP TABLE IF EXISTS `LandingPage`;
CREATE TABLE `LandingPage` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `LandingPageData`
-- ----------------------------
DROP TABLE IF EXISTS `LandingPageData`;
CREATE TABLE `LandingPageData` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `landingPageID` int(11) NOT NULL,
  `dataKey` varchar(100) NOT NULL,
  `dataValue` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `FK__LandingPa__landi__5A2413EF` (`landingPageID`),
  CONSTRAINT `FK__LandingPa__landi__5A2413EF` FOREIGN KEY (`landingPageID`) REFERENCES `LandingPage` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `LandingPageImages`
-- ----------------------------
DROP TABLE IF EXISTS `LandingPageImages`;
CREATE TABLE `LandingPageImages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `landingPageID` int(11) NOT NULL,
  `url` varchar(255) NOT NULL,
  `sort` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `FK__LandingPa__landi__555F5ED2` (`landingPageID`),
  CONSTRAINT `FK__LandingPa__landi__555F5ED2` FOREIGN KEY (`landingPageID`) REFERENCES `LandingPage` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Lifestyle_Trailer`
-- ----------------------------
DROP TABLE IF EXISTS `Lifestyle_Trailer`;
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

-- ----------------------------
--  Table structure for `Location_Services`
-- ----------------------------
DROP TABLE IF EXISTS `Location_Services`;
CREATE TABLE `Location_Services` (
  `loc_service_id` int(11) NOT NULL AUTO_INCREMENT,
  `serviceID` int(11) NOT NULL,
  `locationID` int(11) NOT NULL,
  PRIMARY KEY (`loc_service_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Locations`
-- ----------------------------
DROP TABLE IF EXISTS `Locations`;
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

-- ----------------------------
--  Table structure for `Logger`
-- ----------------------------
DROP TABLE IF EXISTS `Logger`;
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

-- ----------------------------
--  Table structure for `LoggerTypes`
-- ----------------------------
DROP TABLE IF EXISTS `LoggerTypes`;
CREATE TABLE `LoggerTypes` (
  `id` varchar(64) NOT NULL,
  `type` varchar(200) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Make`
-- ----------------------------
DROP TABLE IF EXISTS `Make`;
CREATE TABLE `Make` (
  `makeID` int(11) NOT NULL AUTO_INCREMENT,
  `make` varchar(255) NOT NULL,
  PRIMARY KEY (`makeID`)
) ENGINE=InnoDB AUTO_INCREMENT=54 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `MakeModel`
-- ----------------------------
DROP TABLE IF EXISTS `MakeModel`;
CREATE TABLE `MakeModel` (
  `mmID` int(11) NOT NULL AUTO_INCREMENT,
  `makeID` int(11) NOT NULL,
  `modelID` int(11) NOT NULL,
  PRIMARY KEY (`mmID`),
  KEY `IX_MakeModel` (`makeID`,`modelID`),
  KEY `FK__MakeModel__model__4977ADB9` (`modelID`),
  CONSTRAINT `FK__MakeModel__makeI__48838980` FOREIGN KEY (`makeID`) REFERENCES `Make` (`makeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__MakeModel__model__4977ADB9` FOREIGN KEY (`modelID`) REFERENCES `Model` (`modelID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=743 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `MapIcons`
-- ----------------------------
DROP TABLE IF EXISTS `MapIcons`;
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

-- ----------------------------
--  Table structure for `MapPolygon`
-- ----------------------------
DROP TABLE IF EXISTS `MapPolygon`;
CREATE TABLE `MapPolygon` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `stateID` int(11) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=148 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `MapPolygonCoordinates`
-- ----------------------------
DROP TABLE IF EXISTS `MapPolygonCoordinates`;
CREATE TABLE `MapPolygonCoordinates` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `MapPolygonID` int(11) NOT NULL,
  `latitude` double NOT NULL,
  `longitude` double NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=16225 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `MapixCode`
-- ----------------------------
DROP TABLE IF EXISTS `MapixCode`;
CREATE TABLE `MapixCode` (
  `mCodeID` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`mCodeID`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Menu`
-- ----------------------------
DROP TABLE IF EXISTS `Menu`;
CREATE TABLE `Menu` (
  `menuID` int(11) NOT NULL AUTO_INCREMENT,
  `menu_name` varchar(255) NOT NULL,
  `isPrimary` tinyint(1) NOT NULL DEFAULT '0',
  `active` tinyint(1) NOT NULL DEFAULT '1',
  `display_name` varchar(255) DEFAULT NULL,
  `requireAuthentication` tinyint(1) NOT NULL DEFAULT '0',
  `showOnSitemap` tinyint(1) NOT NULL DEFAULT '0',
  `sort` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`menuID`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Menu_SiteContent`
-- ----------------------------
DROP TABLE IF EXISTS `Menu_SiteContent`;
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
  CONSTRAINT `FK__Menu_Site__menuI__208CD6FA` FOREIGN KEY (`menuID`) REFERENCES `Menu` (`menuID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=130 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Model`
-- ----------------------------
DROP TABLE IF EXISTS `Model`;
CREATE TABLE `Model` (
  `modelID` int(11) NOT NULL AUTO_INCREMENT,
  `model` varchar(255) NOT NULL,
  PRIMARY KEY (`modelID`)
) ENGINE=InnoDB AUTO_INCREMENT=714 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ModelStyle`
-- ----------------------------
DROP TABLE IF EXISTS `ModelStyle`;
CREATE TABLE `ModelStyle` (
  `msID` int(11) NOT NULL AUTO_INCREMENT,
  `modelID` int(11) NOT NULL,
  `styleID` int(11) NOT NULL,
  PRIMARY KEY (`msID`),
  KEY `IX_ModelStyle` (`modelID`,`styleID`),
  KEY `FK__ModelStyl__style__4B5FF62B` (`styleID`),
  CONSTRAINT `FK__ModelStyl__model__4A6BD1F2` FOREIGN KEY (`modelID`) REFERENCES `Model` (`modelID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__ModelStyl__style__4B5FF62B` FOREIGN KEY (`styleID`) REFERENCES `Style` (`styleID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1404 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Modules`
-- ----------------------------
DROP TABLE IF EXISTS `Modules`;
CREATE TABLE `Modules` (
  `moduleID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) DEFAULT NULL,
  `path` varchar(100) DEFAULT NULL,
  `image` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`moduleID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `NewsItem`
-- ----------------------------
DROP TABLE IF EXISTS `NewsItem`;
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
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Note`
-- ----------------------------
DROP TABLE IF EXISTS `Note`;
CREATE TABLE `Note` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `vehiclePartID` int(11) NOT NULL,
  `note` varchar(255) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_Note_IX` (`vehiclePartID`),
  CONSTRAINT `FK__Note__vehiclePar__2BFEED3A` FOREIGN KEY (`vehiclePartID`) REFERENCES `vcdb_VehiclePart` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=317674 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `PackageType`
-- ----------------------------
DROP TABLE IF EXISTS `PackageType`;
CREATE TABLE `PackageType` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Part`
-- ----------------------------
DROP TABLE IF EXISTS `Part`;
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
  KEY `IX_Part_Class` (`classID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `PartAttribute`
-- ----------------------------
DROP TABLE IF EXISTS `PartAttribute`;
CREATE TABLE `PartAttribute` (
  `pAttrID` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `value` varchar(255) DEFAULT NULL,
  `field` varchar(255) DEFAULT NULL,
  `sort` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`pAttrID`),
  KEY `IX_PartAttribute_Part` (`partID`),
  CONSTRAINT `FK__PartAttri__partI__4C541A64` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=66685 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `PartGroup`
-- ----------------------------
DROP TABLE IF EXISTS `PartGroup`;
CREATE TABLE `PartGroup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=552 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `PartGroupPart`
-- ----------------------------
DROP TABLE IF EXISTS `PartGroupPart`;
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

-- ----------------------------
--  Table structure for `PartImageSizes`
-- ----------------------------
DROP TABLE IF EXISTS `PartImageSizes`;
CREATE TABLE `PartImageSizes` (
  `sizeID` int(11) NOT NULL AUTO_INCREMENT,
  `size` varchar(25) DEFAULT NULL,
  `dimensions` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`sizeID`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `PartImages`
-- ----------------------------
DROP TABLE IF EXISTS `PartImages`;
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
  CONSTRAINT `FK__PartImage__partI__0E21DDC1` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__PartImage__sizeI__0D2DB988` FOREIGN KEY (`sizeID`) REFERENCES `PartImageSizes` (`sizeID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1328071 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `PartIndex`
-- ----------------------------
DROP TABLE IF EXISTS `PartIndex`;
CREATE TABLE `PartIndex` (
  `partIndexID` bigint(20) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `partIndex` longtext,
  PRIMARY KEY (`partIndexID`)
) ENGINE=InnoDB AUTO_INCREMENT=50339 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `PartPackage`
-- ----------------------------
DROP TABLE IF EXISTS `PartPackage`;
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
  `typeID` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `dimUnit_idx` (`dimensionUOM`),
  KEY `weightUnit_idx` (`weightUOM`),
  KEY `packageUnit_idx` (`packageUOM`),
  KEY `typeUnit_FK_idx` (`typeID`),
  CONSTRAINT `dimUinit_FK` FOREIGN KEY (`dimensionUOM`) REFERENCES `UnitOfMeasure` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `packageUnit_FK` FOREIGN KEY (`packageUOM`) REFERENCES `UnitOfMeasure` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `typeUnit_FK` FOREIGN KEY (`typeID`) REFERENCES `PackageType` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `weightUnit_FK` FOREIGN KEY (`weightUOM`) REFERENCES `UnitOfMeasure` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=4538 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `PartVideo`
-- ----------------------------
DROP TABLE IF EXISTS `PartVideo`;
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
) ENGINE=InnoDB AUTO_INCREMENT=3204 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Post_Category`
-- ----------------------------
DROP TABLE IF EXISTS `Post_Category`;
CREATE TABLE `Post_Category` (
  `postCategoryID` int(11) NOT NULL AUTO_INCREMENT,
  `postID` int(11) NOT NULL,
  `CategoryID` int(11) NOT NULL,
  PRIMARY KEY (`postCategoryID`),
  KEY `FK__Post_Cate__postI__2BC97F7C` (`postID`),
  CONSTRAINT `FK__Post_Cate__postI__2BC97F7C` FOREIGN KEY (`postID`) REFERENCES `Posts` (`postID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Posts`
-- ----------------------------
DROP TABLE IF EXISTS `Posts`;
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

-- ----------------------------
--  Table structure for `Price`
-- ----------------------------
DROP TABLE IF EXISTS `Price`;
CREATE TABLE `Price` (
  `priceID` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `priceType` varchar(255) DEFAULT NULL,
  `price` decimal(8,2) NOT NULL,
  `enforced` bit(1) NOT NULL,
  `dateModified` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`priceID`),
  KEY `IX_Price_Part` (`partID`),
  CONSTRAINT `FK__Price__partID__0A514CDD` FOREIGN KEY (`partID`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=30173 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Region`
-- ----------------------------
DROP TABLE IF EXISTS `Region`;
CREATE TABLE `Region` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `RelatedPart`
-- ----------------------------
DROP TABLE IF EXISTS `RelatedPart`;
CREATE TABLE `RelatedPart` (
  `relPartID` int(11) NOT NULL AUTO_INCREMENT,
  `partID` int(11) NOT NULL,
  `relatedID` bigint(20) NOT NULL,
  `rTypeID` int(11) NOT NULL,
  PRIMARY KEY (`relPartID`),
  KEY `IX_RelatedPart_Part` (`partID`,`relatedID`)
) ENGINE=InnoDB AUTO_INCREMENT=24952 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `RelatedType`
-- ----------------------------
DROP TABLE IF EXISTS `RelatedType`;
CREATE TABLE `RelatedType` (
  `rTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`rTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ReportType`
-- ----------------------------
DROP TABLE IF EXISTS `ReportType`;
CREATE TABLE `ReportType` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Review`
-- ----------------------------
DROP TABLE IF EXISTS `Review`;
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
) ENGINE=InnoDB AUTO_INCREMENT=524 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `SalesRepresentative`
-- ----------------------------
DROP TABLE IF EXISTS `SalesRepresentative`;
CREATE TABLE `SalesRepresentative` (
  `salesRepID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `code` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`salesRepID`)
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Services`
-- ----------------------------
DROP TABLE IF EXISTS `Services`;
CREATE TABLE `Services` (
  `serviceID` int(11) NOT NULL AUTO_INCREMENT,
  `service_title` varchar(255) DEFAULT NULL,
  `description` longtext,
  `service_price` decimal(19,4) NOT NULL DEFAULT '0.0000',
  `hourly` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`serviceID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `SiteContent`
-- ----------------------------
DROP TABLE IF EXISTS `SiteContent`;
CREATE TABLE `SiteContent` (
  `contentID` int(11) NOT NULL AUTO_INCREMENT,
  `content_type` varchar(255) DEFAULT NULL,
  `page_title` varchar(500) DEFAULT NULL,
  `createdDate` datetime NOT NULL,
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
  PRIMARY KEY (`contentID`)
) ENGINE=InnoDB AUTO_INCREMENT=66 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `SiteContentRevision`
-- ----------------------------
DROP TABLE IF EXISTS `SiteContentRevision`;
CREATE TABLE `SiteContentRevision` (
  `revisionID` int(11) NOT NULL AUTO_INCREMENT,
  `contentID` int(11) NOT NULL DEFAULT '1',
  `content_text` longtext,
  `createdOn` datetime NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`revisionID`),
  KEY `FK__SiteConte__conte__151B244E` (`contentID`),
  CONSTRAINT `FK__SiteConte__conte__151B244E` FOREIGN KEY (`contentID`) REFERENCES `SiteContent` (`contentID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=76 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `State`
-- ----------------------------
DROP TABLE IF EXISTS `State`;
CREATE TABLE `State` (
  `state` varchar(128) NOT NULL,
  `abbr` varchar(128) NOT NULL,
  `stateID` int(11) NOT NULL,
  `countryID` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `States`
-- ----------------------------
DROP TABLE IF EXISTS `States`;
CREATE TABLE `States` (
  `stateID` int(11) NOT NULL AUTO_INCREMENT,
  `state` varchar(100) NOT NULL,
  `abbr` varchar(3) NOT NULL,
  `countryID` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`stateID`),
  KEY `FK__States__countryI__607251E5` (`countryID`),
  CONSTRAINT `FK__States__countryI__607251E5` FOREIGN KEY (`countryID`) REFERENCES `Country` (`countryID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=85 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Style`
-- ----------------------------
DROP TABLE IF EXISTS `Style`;
CREATE TABLE `Style` (
  `styleID` int(11) NOT NULL AUTO_INCREMENT,
  `style` varchar(255) NOT NULL,
  `aaiaID` int(11) NOT NULL,
  PRIMARY KEY (`styleID`)
) ENGINE=InnoDB AUTO_INCREMENT=631 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Submodel`
-- ----------------------------
DROP TABLE IF EXISTS `Submodel`;
CREATE TABLE `Submodel` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `AAIASubmodelID` int(11) DEFAULT NULL,
  `SubmodelName` varchar(50) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_Submodel_IX` (`AAIASubmodelID`)
) ENGINE=InnoDB AUTO_INCREMENT=1729 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `TechNews`
-- ----------------------------
DROP TABLE IF EXISTS `TechNews`;
CREATE TABLE `TechNews` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pageContent` text,
  `showDealers` tinyint(1) NOT NULL,
  `showPublic` tinyint(1) NOT NULL,
  `dateModified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `displayOrder` int(11) DEFAULT NULL,
  `active` tinyint(1) NOT NULL,
  `title` varchar(500) NOT NULL,
  `subTitle` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Testimonial`
-- ----------------------------
DROP TABLE IF EXISTS `Testimonial`;
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
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Testimonials`
-- ----------------------------
DROP TABLE IF EXISTS `Testimonials`;
CREATE TABLE `Testimonials` (
  `reviewID` int(11) NOT NULL AUTO_INCREMENT,
  `reviewer` varchar(400) DEFAULT NULL,
  `review` longtext,
  `date_added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_new` int(11) NOT NULL DEFAULT '0',
  `is_hidden` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`reviewID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Trailer`
-- ----------------------------
DROP TABLE IF EXISTS `Trailer`;
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

-- ----------------------------
--  Table structure for `Tweets`
-- ----------------------------
DROP TABLE IF EXISTS `Tweets`;
CREATE TABLE `Tweets` (
  `tweetID` int(11) NOT NULL AUTO_INCREMENT,
  `twitterTweetID` varchar(500) NOT NULL,
  `tweet` varchar(150) NOT NULL,
  `postDate` datetime NOT NULL,
  `twitterUserID` varchar(500) NOT NULL,
  `screenName` varchar(100) NOT NULL,
  `profilePhoto` varchar(500) NOT NULL,
  PRIMARY KEY (`tweetID`),
  UNIQUE KEY `tweetID_UNIQUE` (`tweetID`)
) ENGINE=InnoDB AUTO_INCREMENT=755 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `UnitOfMeasure`
-- ----------------------------
DROP TABLE IF EXISTS `UnitOfMeasure`;
CREATE TABLE `UnitOfMeasure` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `code` varchar(5) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `UserProfiles`
-- ----------------------------
DROP TABLE IF EXISTS `UserProfiles`;
CREATE TABLE `UserProfiles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `customerID` int(11) DEFAULT NULL,
  `custID` int(11) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `IP` varchar(25) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=715 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Vehicle`
-- ----------------------------
DROP TABLE IF EXISTS `Vehicle`;
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
) ENGINE=InnoDB AUTO_INCREMENT=248743 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `VehicleConfig`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleConfig`;
CREATE TABLE `VehicleConfig` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `AAIAVehicleConfigID` int(11) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_VehicleConfig_IX` (`AAIAVehicleConfigID`)
) ENGINE=InnoDB AUTO_INCREMENT=28348 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `VehicleConfigAttribute`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleConfigAttribute`;
CREATE TABLE `VehicleConfigAttribute` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `AttributeID` int(11) NOT NULL,
  `VehicleConfigID` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `FK__VehicleCo__Attri__19E03CFF` (`AttributeID`),
  KEY `FK__VehicleCo__Vehic__1AD46138` (`VehicleConfigID`),
  CONSTRAINT `FK__VehicleCo__Attri__19E03CFF` FOREIGN KEY (`AttributeID`) REFERENCES `ConfigAttribute` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__VehicleCo__Vehic__1AD46138` FOREIGN KEY (`VehicleConfigID`) REFERENCES `VehicleConfig` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=39846 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `VehiclePart`
-- ----------------------------
DROP TABLE IF EXISTS `VehiclePart`;
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
) ENGINE=InnoDB AUTO_INCREMENT=42585 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `VehiclePartAttribute`
-- ----------------------------
DROP TABLE IF EXISTS `VehiclePartAttribute`;
CREATE TABLE `VehiclePartAttribute` (
  `vpAttrID` int(11) NOT NULL AUTO_INCREMENT,
  `vPartID` int(11) NOT NULL,
  `value` varchar(255) DEFAULT NULL,
  `field` varchar(255) DEFAULT NULL,
  `sort` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`vpAttrID`),
  KEY `IX_VehiclePartAttr_VPart` (`vPartID`),
  CONSTRAINT `FK__VehiclePa__vPart__520CF3BA` FOREIGN KEY (`vPartID`) REFERENCES `VehiclePart` (`vPartID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=102410 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `VehicleType`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleType`;
CREATE TABLE `VehicleType` (
  `VehicleTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `VehicleTypeName` varchar(50) NOT NULL,
  `VehicleTypeGroupID` int(11) DEFAULT NULL,
  PRIMARY KEY (`VehicleTypeID`),
  KEY `FK__VehicleTy__Vehic__648AFD1B` (`VehicleTypeGroupID`),
  CONSTRAINT `FK__VehicleTy__Vehic__648AFD1B` FOREIGN KEY (`VehicleTypeGroupID`) REFERENCES `VehicleTypeGroup` (`VehicleTypeGroupID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `VehicleTypeGroup`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleTypeGroup`;
CREATE TABLE `VehicleTypeGroup` (
  `VehicleTypeGroupID` int(11) NOT NULL AUTO_INCREMENT,
  `VehicleTypeGroupName` varchar(50) NOT NULL,
  PRIMARY KEY (`VehicleTypeGroupID`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Video`
-- ----------------------------
DROP TABLE IF EXISTS `Video`;
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
) ENGINE=InnoDB AUTO_INCREMENT=71 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `WebPropNotes`
-- ----------------------------
DROP TABLE IF EXISTS `WebPropNotes`;
CREATE TABLE `WebPropNotes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `webPropID` int(11) NOT NULL,
  `text` varchar(255) NOT NULL,
  `dateAdded` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `WebPropRequirementCheck`
-- ----------------------------
DROP TABLE IF EXISTS `WebPropRequirementCheck`;
CREATE TABLE `WebPropRequirementCheck` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `WebPropertiesID` int(11) DEFAULT NULL,
  `Compliance` tinyint(1) DEFAULT NULL,
  `WebPropRequirementsID` int(11) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `webPropID_ref_idx` (`WebPropertiesID`),
  KEY `webPropReqID_ref_idx` (`WebPropRequirementsID`),
  CONSTRAINT `webPropID_ref` FOREIGN KEY (`WebPropertiesID`) REFERENCES `WebProperties` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `webPropReqID_ref` FOREIGN KEY (`WebPropRequirementsID`) REFERENCES `WebPropRequirements` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=825 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `WebPropRequirements`
-- ----------------------------
DROP TABLE IF EXISTS `WebPropRequirements`;
CREATE TABLE `WebPropRequirements` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `ReqType` varchar(255) DEFAULT NULL,
  `Requirement` varchar(1000) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `WebProperties`
-- ----------------------------
DROP TABLE IF EXISTS `WebProperties`;
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
  UNIQUE KEY `badgeID` (`badgeID`),
  KEY `type_ref_idx` (`typeID`),
  CONSTRAINT `type_ref` FOREIGN KEY (`typeID`) REFERENCES `WebPropertyTypes` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=193 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `WebPropertyTypes`
-- ----------------------------
DROP TABLE IF EXISTS `WebPropertyTypes`;
CREATE TABLE `WebPropertyTypes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `typeID` int(11) NOT NULL,
  `type` varchar(50) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Website`
-- ----------------------------
DROP TABLE IF EXISTS `Website`;
CREATE TABLE `Website` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `WidgetDeployments`
-- ----------------------------
DROP TABLE IF EXISTS `WidgetDeployments`;
CREATE TABLE `WidgetDeployments` (
  `trackerID` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(400) NOT NULL,
  `date_added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`trackerID`)
) ENGINE=InnoDB AUTO_INCREMENT=80 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `Year`
-- ----------------------------
DROP TABLE IF EXISTS `Year`;
CREATE TABLE `Year` (
  `yearID` int(11) NOT NULL AUTO_INCREMENT,
  `year` double NOT NULL,
  PRIMARY KEY (`yearID`)
) ENGINE=InnoDB AUTO_INCREMENT=283 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `YearMake`
-- ----------------------------
DROP TABLE IF EXISTS `YearMake`;
CREATE TABLE `YearMake` (
  `ymID` int(11) NOT NULL AUTO_INCREMENT,
  `yearID` int(11) DEFAULT NULL,
  `makeID` int(11) NOT NULL,
  PRIMARY KEY (`ymID`),
  KEY `IX_YearMake` (`yearID`,`makeID`),
  KEY `FK__YearMake__makeID__478F6547` (`makeID`),
  CONSTRAINT `FK__YearMake__makeID__478F6547` FOREIGN KEY (`makeID`) REFERENCES `Make` (`makeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__YearMake__yearID__469B410E` FOREIGN KEY (`yearID`) REFERENCES `Year` (`yearID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1327 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `vcdb_Make`
-- ----------------------------
DROP TABLE IF EXISTS `vcdb_Make`;
CREATE TABLE `vcdb_Make` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `AAIAMakeID` int(11) DEFAULT NULL,
  `MakeName` varchar(50) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_Make_IX` (`AAIAMakeID`)
) ENGINE=InnoDB AUTO_INCREMENT=53 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `vcdb_Model`
-- ----------------------------
DROP TABLE IF EXISTS `vcdb_Model`;
CREATE TABLE `vcdb_Model` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `AAIAModelID` int(11) DEFAULT NULL,
  `ModelName` varchar(100) DEFAULT NULL,
  `VehicleTypeID` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_Model_IX` (`AAIAModelID`)
) ENGINE=InnoDB AUTO_INCREMENT=3877 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `vcdb_Vehicle`
-- ----------------------------
DROP TABLE IF EXISTS `vcdb_Vehicle`;
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
) ENGINE=InnoDB AUTO_INCREMENT=35130 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `vcdb_VehiclePart`
-- ----------------------------
DROP TABLE IF EXISTS `vcdb_VehiclePart`;
CREATE TABLE `vcdb_VehiclePart` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `VehicleID` int(11) NOT NULL,
  `PartNumber` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `AAIA_VehiclePart_Part_IX` (`VehicleID`,`PartNumber`),
  KEY `FK__vcdb_Vehi__PartN__273A381D` (`PartNumber`),
  CONSTRAINT `FK__vcdb_Vehi__PartN__273A381D` FOREIGN KEY (`PartNumber`) REFERENCES `Part` (`partID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK__vcdb_Vehi__Vehic__264613E4` FOREIGN KEY (`VehicleID`) REFERENCES `vcdb_Vehicle` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=202619 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `vcdb_Year`
-- ----------------------------
DROP TABLE IF EXISTS `vcdb_Year`;
CREATE TABLE `vcdb_Year` (
  `YearID` int(11) NOT NULL,
  PRIMARY KEY (`YearID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `videoType`
-- ----------------------------
DROP TABLE IF EXISTS `videoType`;
CREATE TABLE `videoType` (
  `vTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `icon` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`vTypeID`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
