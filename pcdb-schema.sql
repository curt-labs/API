/*
 Navicat Premium Data Transfer

 Source Server         : Localhost
 Source Server Type    : MySQL
 Source Server Version : 50614
 Source Host           : localhost
 Source Database       : pcdb

 Target Server Type    : MySQL
 Target Server Version : 50614
 File Encoding         : utf-8

 Date: 12/03/2013 09:07:13 AM
*/

CREATE DATABASE IF NOT EXISTS pcdb;

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `Categories`
-- ----------------------------
DROP TABLE IF EXISTS `Categories`;
CREATE TABLE `Categories` (
  `CategoryID` int(10) NOT NULL,
  `CategoryName` varchar(100) NOT NULL,
  PRIMARY KEY (`CategoryID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `CodeMaster`
-- ----------------------------
DROP TABLE IF EXISTS `CodeMaster`;
CREATE TABLE `CodeMaster` (
  `CodeMasterID` int(10) NOT NULL,
  `CategoryID` int(10) NOT NULL,
  `SubCategoryID` int(10) NOT NULL,
  `PartTerminologyID` int(10) NOT NULL,
  `PositionID` int(10) NOT NULL,
  `RevDate` datetime DEFAULT NULL,
  PRIMARY KEY (`CodeMasterID`),
  UNIQUE KEY `CodeMaster_Unique` (`CategoryID`,`SubCategoryID`,`PartTerminologyID`,`PositionID`),
  KEY `IX_CodeMaster_CategoryID` (`CategoryID`),
  KEY `IX_CodeMaster_PartTerminologyID` (`PartTerminologyID`),
  KEY `IX_CodeMaster_PositionID` (`PositionID`),
  KEY `IX_CodeMaster_SubCategoryID` (`SubCategoryID`),
  CONSTRAINT `FK_Categories_CodeMaster` FOREIGN KEY (`CategoryID`) REFERENCES `categories` (`CategoryID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_Parts_CodeMaster` FOREIGN KEY (`PartTerminologyID`) REFERENCES `parts` (`PartTerminologyID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_Positions_CodeMaster` FOREIGN KEY (`PositionID`) REFERENCES `positions` (`PositionID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_Subcategories_CodeMaster` FOREIGN KEY (`SubCategoryID`) REFERENCES `subcategories` (`SubCategoryID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `PCdbChanges`
-- ----------------------------
DROP TABLE IF EXISTS `PCdbChanges`;
CREATE TABLE `PCdbChanges` (
  `CodeMasterID` int(10) DEFAULT NULL,
  `CategoryID` int(10) DEFAULT NULL,
  `CategoryName` varchar(100) DEFAULT NULL,
  `SubCategoryID` int(10) DEFAULT NULL,
  `SubCategoryName` varchar(100) DEFAULT NULL,
  `PartTerminologyID` int(10) DEFAULT NULL,
  `PartTerminologyName` varchar(100) DEFAULT NULL,
  `UseID` int(10) DEFAULT NULL,
  `UseDescription` varchar(100) DEFAULT NULL,
  `PositionID` int(10) DEFAULT NULL,
  `Position` varchar(100) DEFAULT NULL,
  `RevDATE` datetime DEFAULT NULL,
  `Action` varchar(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Parts`
-- ----------------------------
DROP TABLE IF EXISTS `Parts`;
CREATE TABLE `Parts` (
  `PartTerminologyID` int(10) NOT NULL,
  `PartTerminologyName` varchar(100) NOT NULL,
  `RevDate` datetime DEFAULT NULL,
  PRIMARY KEY (`PartTerminologyID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `PartsToUse`
-- ----------------------------
DROP TABLE IF EXISTS `PartsToUse`;
CREATE TABLE `PartsToUse` (
  `PartTerminologyId` int(10) NOT NULL,
  `UseId` int(10) NOT NULL,
  PRIMARY KEY (`PartTerminologyId`,`UseId`),
  KEY `FK_PartsToUse_Use` (`UseId`),
  CONSTRAINT `FK_PartsToUse_Parts` FOREIGN KEY (`PartTerminologyId`) REFERENCES `parts` (`PartTerminologyID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_PartsToUse_Use` FOREIGN KEY (`UseId`) REFERENCES `use` (`UseId`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Positions`
-- ----------------------------
DROP TABLE IF EXISTS `Positions`;
CREATE TABLE `Positions` (
  `PositionID` int(10) NOT NULL,
  `Position` varchar(100) NOT NULL,
  PRIMARY KEY (`PositionID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Retired Terms`
-- ----------------------------
DROP TABLE IF EXISTS `Retired Terms`;
CREATE TABLE `Retired Terms` (
  `PartName` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `PartIDCode` int(10) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Subcategories`
-- ----------------------------
DROP TABLE IF EXISTS `Subcategories`;
CREATE TABLE `Subcategories` (
  `SubCategoryID` int(10) NOT NULL,
  `SubCategoryName` varchar(100) NOT NULL,
  PRIMARY KEY (`SubCategoryID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Use`
-- ----------------------------
DROP TABLE IF EXISTS `Use`;
CREATE TABLE `Use` (
  `UseId` int(10) NOT NULL,
  `UseDescription` varchar(100) NOT NULL,
  PRIMARY KEY (`UseId`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Version`
-- ----------------------------
DROP TABLE IF EXISTS `Version`;
CREATE TABLE `Version` (
  `VersionDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

SET FOREIGN_KEY_CHECKS = 1;
