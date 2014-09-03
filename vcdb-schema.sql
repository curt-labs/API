/*
 Navicat Premium Data Transfer

 Source Server         : 173.255.114.206
 Source Server Type    : MySQL
 Source Server Version : 50617
 Source Host           : 173.255.114.206
 Source Database       : vcdb

 Target Server Type    : MySQL
 Target Server Version : 50617
 File Encoding         : utf-8

 Date: 09/03/2014 11:47:00 AM
*/

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `Abbreviation`
-- ----------------------------
DROP TABLE IF EXISTS `Abbreviation`;
CREATE TABLE `Abbreviation` (
  `Abbreviation` char(3) NOT NULL,
  `Description` varchar(20) NOT NULL,
  `LongDescription` varchar(200) NOT NULL DEFAULT '',
  PRIMARY KEY (`Abbreviation`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Aspiration`
-- ----------------------------
DROP TABLE IF EXISTS `Aspiration`;
CREATE TABLE `Aspiration` (
  `AspirationID` int(10) NOT NULL,
  `AspirationName` varchar(30) NOT NULL,
  PRIMARY KEY (`AspirationID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Attachment`
-- ----------------------------
DROP TABLE IF EXISTS `Attachment`;
CREATE TABLE `Attachment` (
  `AttachmentID` int(10) NOT NULL AUTO_INCREMENT,
  `AttachmentTypeID` int(10) NOT NULL,
  `AttachmentFileName` varchar(50) NOT NULL,
  `AttachmentURL` varchar(100) NOT NULL,
  `AttachmentDescription` varchar(50) NOT NULL,
  PRIMARY KEY (`AttachmentID`),
  KEY `IX_Attachment_AttachmentTypeID` (`AttachmentTypeID`),
  CONSTRAINT `attachmenttypeattachment_fk` FOREIGN KEY (`AttachmentTypeID`) REFERENCES `attachmenttype` (`AttachmentTypeID`),
  CONSTRAINT `FK_AttachmentType_Attachment` FOREIGN KEY (`AttachmentTypeID`) REFERENCES `attachmenttype` (`AttachmentTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `AttachmentType`
-- ----------------------------
DROP TABLE IF EXISTS `AttachmentType`;
CREATE TABLE `AttachmentType` (
  `AttachmentTypeID` int(10) NOT NULL AUTO_INCREMENT,
  `AttachmentTypeName` varchar(20) NOT NULL,
  PRIMARY KEY (`AttachmentTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `BaseVehicle`
-- ----------------------------
DROP TABLE IF EXISTS `BaseVehicle`;
CREATE TABLE `BaseVehicle` (
  `BaseVehicleID` int(10) NOT NULL,
  `YearID` int(10) NOT NULL,
  `MakeID` int(10) NOT NULL,
  `ModelID` int(10) NOT NULL,
  PRIMARY KEY (`BaseVehicleID`),
  KEY `IDX_BaseVehicle_MakeID` (`MakeID`),
  KEY `IDX_BaseVehicle_ModelID` (`ModelID`),
  KEY `IDX_BaseVehicle_YearID` (`YearID`),
  CONSTRAINT `FK_Make_BaseVehicle` FOREIGN KEY (`MakeID`) REFERENCES `make` (`MakeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_Model_BaseVehicle` FOREIGN KEY (`ModelID`) REFERENCES `Model` (`ModelID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_Year_BaseVehicle` FOREIGN KEY (`YearID`) REFERENCES `year` (`YearID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `makebasevehicle_fk` FOREIGN KEY (`MakeID`) REFERENCES `make` (`MakeID`),
  CONSTRAINT `modelbasevehicle_fk` FOREIGN KEY (`ModelID`) REFERENCES `Model` (`ModelID`),
  CONSTRAINT `yearbasevehicle_fk` FOREIGN KEY (`YearID`) REFERENCES `year` (`YearID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `BedConfig`
-- ----------------------------
DROP TABLE IF EXISTS `BedConfig`;
CREATE TABLE `BedConfig` (
  `BedConfigID` int(10) NOT NULL,
  `BedLengthID` int(10) NOT NULL,
  `BedTypeID` int(10) NOT NULL,
  PRIMARY KEY (`BedConfigID`),
  KEY `IDX_BedConfig_BedLengthID` (`BedLengthID`),
  KEY `IDX_BedConfig_BedTypeID` (`BedTypeID`),
  CONSTRAINT `bedlengthbedconfig_fk` FOREIGN KEY (`BedLengthID`) REFERENCES `bedlength` (`BedLengthID`),
  CONSTRAINT `bedtypebedconfig_fk` FOREIGN KEY (`BedTypeID`) REFERENCES `bedtype` (`BedTypeID`),
  CONSTRAINT `FK_BedLength_BedConfig` FOREIGN KEY (`BedLengthID`) REFERENCES `bedlength` (`BedLengthID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_BedType_BedConfig` FOREIGN KEY (`BedTypeID`) REFERENCES `bedtype` (`BedTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `BedLength`
-- ----------------------------
DROP TABLE IF EXISTS `BedLength`;
CREATE TABLE `BedLength` (
  `BedLengthID` int(10) NOT NULL,
  `BedLength` char(10) NOT NULL,
  `BedLengthMetric` char(10) NOT NULL,
  PRIMARY KEY (`BedLengthID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `BedType`
-- ----------------------------
DROP TABLE IF EXISTS `BedType`;
CREATE TABLE `BedType` (
  `BedTypeID` int(10) NOT NULL,
  `BedTypeName` varchar(50) NOT NULL,
  PRIMARY KEY (`BedTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `BodyNumDoors`
-- ----------------------------
DROP TABLE IF EXISTS `BodyNumDoors`;
CREATE TABLE `BodyNumDoors` (
  `BodyNumDoorsID` int(10) NOT NULL,
  `BodyNumDoors` char(3) NOT NULL,
  PRIMARY KEY (`BodyNumDoorsID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `BodyStyleConfig`
-- ----------------------------
DROP TABLE IF EXISTS `BodyStyleConfig`;
CREATE TABLE `BodyStyleConfig` (
  `BodyStyleConfigID` int(10) NOT NULL,
  `BodyNumDoorsID` int(10) NOT NULL,
  `BodyTypeID` int(10) NOT NULL,
  PRIMARY KEY (`BodyStyleConfigID`),
  KEY `IDX_BodyStyleConfig_BodyNumDoor` (`BodyNumDoorsID`),
  KEY `IDX_BodyStyleConfig_BodyTypeID` (`BodyTypeID`),
  CONSTRAINT `bodynumdoorsbodystyleconfig_fk` FOREIGN KEY (`BodyNumDoorsID`) REFERENCES `bodynumdoors` (`BodyNumDoorsID`),
  CONSTRAINT `bodytypebodystyleconfig_fk` FOREIGN KEY (`BodyTypeID`) REFERENCES `bodytype` (`BodyTypeID`),
  CONSTRAINT `FK_BodyNumDoors_BodyStyleConfig` FOREIGN KEY (`BodyNumDoorsID`) REFERENCES `bodynumdoors` (`BodyNumDoorsID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_BodyType_BodyStyleConfig` FOREIGN KEY (`BodyTypeID`) REFERENCES `bodytype` (`BodyTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `BodyType`
-- ----------------------------
DROP TABLE IF EXISTS `BodyType`;
CREATE TABLE `BodyType` (
  `BodyTypeID` int(10) NOT NULL,
  `BodyTypeName` varchar(50) NOT NULL,
  PRIMARY KEY (`BodyTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `BrakeABS`
-- ----------------------------
DROP TABLE IF EXISTS `BrakeABS`;
CREATE TABLE `BrakeABS` (
  `BrakeABSID` int(10) NOT NULL,
  `BrakeABSName` varchar(30) NOT NULL,
  PRIMARY KEY (`BrakeABSID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `BrakeConfig`
-- ----------------------------
DROP TABLE IF EXISTS `BrakeConfig`;
CREATE TABLE `BrakeConfig` (
  `BrakeConfigID` int(10) NOT NULL,
  `FrontBrakeTypeID` int(10) NOT NULL,
  `RearBrakeTypeID` int(10) NOT NULL,
  `BrakeSystemID` int(10) NOT NULL,
  `BrakeABSID` int(10) NOT NULL,
  PRIMARY KEY (`BrakeConfigID`),
  KEY `IDX_BrakeConfig_BrakeABSID` (`BrakeABSID`),
  KEY `IDX_BrakeConfig_BrakeSystemID` (`BrakeSystemID`),
  KEY `IDX_BrakeConfig_FrontBrakeTypeID` (`FrontBrakeTypeID`),
  KEY `IDX_BrakeConfig_RearBrakeTypeID` (`RearBrakeTypeID`),
  CONSTRAINT `brakeabsbrakeconfig_fk` FOREIGN KEY (`BrakeABSID`) REFERENCES `brakeabs` (`BrakeABSID`),
  CONSTRAINT `brakesystembrakeconfig_fk` FOREIGN KEY (`BrakeSystemID`) REFERENCES `brakesystem` (`BrakeSystemID`),
  CONSTRAINT `braketypebrakeconfig1_fk` FOREIGN KEY (`RearBrakeTypeID`) REFERENCES `braketype` (`BrakeTypeID`),
  CONSTRAINT `braketypebrakeconfig_fk` FOREIGN KEY (`FrontBrakeTypeID`) REFERENCES `braketype` (`BrakeTypeID`),
  CONSTRAINT `FK_BrakeABS_BrakeConfig` FOREIGN KEY (`BrakeABSID`) REFERENCES `brakeabs` (`BrakeABSID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_BrakeSystem_BrakeConfig` FOREIGN KEY (`BrakeSystemID`) REFERENCES `brakesystem` (`BrakeSystemID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_FrontBrakeType_BrakeConfig` FOREIGN KEY (`FrontBrakeTypeID`) REFERENCES `braketype` (`BrakeTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_RearBrakeType_BrakeConfig` FOREIGN KEY (`RearBrakeTypeID`) REFERENCES `braketype` (`BrakeTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `BrakeSystem`
-- ----------------------------
DROP TABLE IF EXISTS `BrakeSystem`;
CREATE TABLE `BrakeSystem` (
  `BrakeSystemID` int(10) NOT NULL,
  `BrakeSystemName` varchar(30) NOT NULL,
  PRIMARY KEY (`BrakeSystemID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `BrakeType`
-- ----------------------------
DROP TABLE IF EXISTS `BrakeType`;
CREATE TABLE `BrakeType` (
  `BrakeTypeID` int(10) NOT NULL,
  `BrakeTypeName` varchar(30) NOT NULL,
  PRIMARY KEY (`BrakeTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `CylinderHeadType`
-- ----------------------------
DROP TABLE IF EXISTS `CylinderHeadType`;
CREATE TABLE `CylinderHeadType` (
  `CylinderHeadTypeID` int(10) NOT NULL,
  `CylinderHeadTypeName` varchar(30) NOT NULL,
  PRIMARY KEY (`CylinderHeadTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `DriveType`
-- ----------------------------
DROP TABLE IF EXISTS `DriveType`;
CREATE TABLE `DriveType` (
  `DriveTypeID` int(10) NOT NULL,
  `DriveTypeName` varchar(30) NOT NULL,
  PRIMARY KEY (`DriveTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `ElecControlled`
-- ----------------------------
DROP TABLE IF EXISTS `ElecControlled`;
CREATE TABLE `ElecControlled` (
  `ElecControlledID` int(10) NOT NULL,
  `ElecControlled` char(3) NOT NULL,
  PRIMARY KEY (`ElecControlledID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `EngineBase`
-- ----------------------------
DROP TABLE IF EXISTS `EngineBase`;
CREATE TABLE `EngineBase` (
  `EngineBaseID` int(10) NOT NULL,
  `Liter` char(6) NOT NULL,
  `CC` char(8) NOT NULL,
  `CID` char(7) NOT NULL,
  `Cylinders` char(2) NOT NULL,
  `BlockType` char(2) NOT NULL,
  `EngBoreIn` char(10) NOT NULL,
  `EngBoreMetric` char(10) NOT NULL,
  `EngStrokeIn` char(10) NOT NULL,
  `EngStrokeMetric` char(10) NOT NULL,
  PRIMARY KEY (`EngineBaseID`),
  KEY `IDX_EngineBase_BlockType` (`BlockType`),
  KEY `IDX_EngineBase_CC` (`CC`),
  KEY `IDX_EngineBase_CID` (`CID`),
  KEY `IDX_EngineBase_Cylinders` (`Cylinders`),
  KEY `IDX_EngineBase_EngBoreIn` (`EngBoreIn`),
  KEY `IDX_EngineBase_EngBoreMetric` (`EngBoreMetric`),
  KEY `IDX_EngineBase_EngStrokeIn` (`EngStrokeIn`),
  KEY `IDX_EngineBase_EngStrokeMetric` (`EngStrokeMetric`),
  KEY `IDX_EngineBase_Liter` (`Liter`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `EngineConfig`
-- ----------------------------
DROP TABLE IF EXISTS `EngineConfig`;
CREATE TABLE `EngineConfig` (
  `EngineConfigID` int(10) NOT NULL,
  `EngineDesignationID` int(10) NOT NULL,
  `EngineVINID` int(10) NOT NULL,
  `ValvesID` int(10) NOT NULL,
  `EngineBaseID` int(10) NOT NULL,
  `FuelDeliveryConfigID` int(10) NOT NULL,
  `AspirationID` int(10) NOT NULL,
  `CylinderHeadTypeID` int(10) NOT NULL,
  `FuelTypeID` int(10) NOT NULL,
  `IgnitionSystemTypeID` int(10) NOT NULL,
  `EngineMfrID` int(10) NOT NULL,
  `EngineVersionID` int(10) NOT NULL,
  `PowerOutputID` int(10) NOT NULL DEFAULT '1',
  PRIMARY KEY (`EngineConfigID`),
  KEY `IDX_EngineConfig_AspirationID` (`AspirationID`),
  KEY `IDX_EngineConfig_CylinderHeadTy` (`CylinderHeadTypeID`),
  KEY `IDX_EngineConfig_EngineBaseID` (`EngineBaseID`),
  KEY `IDX_EngineConfig_EngineDesignat` (`EngineDesignationID`),
  KEY `IDX_EngineConfig_EngineMfrID` (`EngineMfrID`),
  KEY `IDX_EngineConfig_EngineVersionID` (`EngineVersionID`),
  KEY `IDX_EngineConfig_EngineVINID` (`EngineVINID`),
  KEY `IDX_EngineConfig_FuelDeliveryCo` (`FuelDeliveryConfigID`),
  KEY `IDX_EngineConfig_FuelTypeID` (`FuelTypeID`),
  KEY `IDX_EngineConfig_IgnitionSystem` (`IgnitionSystemTypeID`),
  KEY `FK_EngineConfig_Valves1` (`ValvesID`),
  CONSTRAINT `aspirationengineconfig_fk` FOREIGN KEY (`AspirationID`) REFERENCES `aspiration` (`AspirationID`),
  CONSTRAINT `cylinderheadtypeengineconfig_fk` FOREIGN KEY (`CylinderHeadTypeID`) REFERENCES `cylinderheadtype` (`CylinderHeadTypeID`),
  CONSTRAINT `enginebaseengineconfig_fk` FOREIGN KEY (`EngineBaseID`) REFERENCES `enginebase` (`EngineBaseID`),
  CONSTRAINT `enginedesignationengineconfi_fk` FOREIGN KEY (`EngineDesignationID`) REFERENCES `enginedesignation` (`EngineDesignationID`),
  CONSTRAINT `engineversionengineconfig_fk` FOREIGN KEY (`EngineVersionID`) REFERENCES `engineversion` (`EngineVersionID`),
  CONSTRAINT `enginevinengineconfig_fk` FOREIGN KEY (`EngineVINID`) REFERENCES `enginevin` (`EngineVINID`),
  CONSTRAINT `FK_Aspiration_EngineConfig` FOREIGN KEY (`AspirationID`) REFERENCES `aspiration` (`AspirationID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_CylinderHeadType_EngineConfig` FOREIGN KEY (`CylinderHeadTypeID`) REFERENCES `cylinderheadtype` (`CylinderHeadTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_EngineBase_EngineConfig` FOREIGN KEY (`EngineBaseID`) REFERENCES `enginebase` (`EngineBaseID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_EngineConfig_Valves` FOREIGN KEY (`ValvesID`) REFERENCES `valves` (`ValvesID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_EngineConfig_Valves1` FOREIGN KEY (`ValvesID`) REFERENCES `valves` (`ValvesID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_EngineDesignation_EngineConfig` FOREIGN KEY (`EngineDesignationID`) REFERENCES `enginedesignation` (`EngineDesignationID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_EngineVersion_EngineConfig` FOREIGN KEY (`EngineVersionID`) REFERENCES `engineversion` (`EngineVersionID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_EngineVIN_EngineConfig` FOREIGN KEY (`EngineVINID`) REFERENCES `enginevin` (`EngineVINID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_FuelDeliveryConfig_EngineConfig` FOREIGN KEY (`FuelDeliveryConfigID`) REFERENCES `FuelDeliveryConfig` (`FuelDeliveryConfigID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_FuelType_EngineConfig` FOREIGN KEY (`FuelTypeID`) REFERENCES `fueltype` (`FuelTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_IgnitionSystemType_EngineConfig` FOREIGN KEY (`IgnitionSystemTypeID`) REFERENCES `ignitionsystemtype` (`IgnitionSystemTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_Mfr_EngineConfig` FOREIGN KEY (`EngineMfrID`) REFERENCES `mfr` (`MfrID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fueldeliveryconfigengineconf_fk` FOREIGN KEY (`FuelDeliveryConfigID`) REFERENCES `FuelDeliveryConfig` (`FuelDeliveryConfigID`),
  CONSTRAINT `fueltypeengineconfig_fk` FOREIGN KEY (`FuelTypeID`) REFERENCES `fueltype` (`FuelTypeID`),
  CONSTRAINT `ignitionsystemtypeengineconf_fk` FOREIGN KEY (`IgnitionSystemTypeID`) REFERENCES `ignitionsystemtype` (`IgnitionSystemTypeID`),
  CONSTRAINT `mfrengineconfig_fk` FOREIGN KEY (`EngineMfrID`) REFERENCES `mfr` (`MfrID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `EngineDesignation`
-- ----------------------------
DROP TABLE IF EXISTS `EngineDesignation`;
CREATE TABLE `EngineDesignation` (
  `EngineDesignationID` int(10) NOT NULL,
  `EngineDesignationName` varchar(30) NOT NULL,
  PRIMARY KEY (`EngineDesignationID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `EngineVIN`
-- ----------------------------
DROP TABLE IF EXISTS `EngineVIN`;
CREATE TABLE `EngineVIN` (
  `EngineVINID` int(10) NOT NULL,
  `EngineVINName` varchar(5) NOT NULL,
  PRIMARY KEY (`EngineVINID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `EngineVersion`
-- ----------------------------
DROP TABLE IF EXISTS `EngineVersion`;
CREATE TABLE `EngineVersion` (
  `EngineVersionID` int(10) NOT NULL,
  `EngineVersion` varchar(20) NOT NULL,
  PRIMARY KEY (`EngineVersionID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `EnglishPhrase`
-- ----------------------------
DROP TABLE IF EXISTS `EnglishPhrase`;
CREATE TABLE `EnglishPhrase` (
  `EnglishPhraseID` int(10) NOT NULL AUTO_INCREMENT,
  `EnglishPhrase` varchar(100) NOT NULL,
  PRIMARY KEY (`EnglishPhraseID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `FuelDeliveryConfig`
-- ----------------------------
DROP TABLE IF EXISTS `FuelDeliveryConfig`;
CREATE TABLE `FuelDeliveryConfig` (
  `FuelDeliveryConfigID` int(10) NOT NULL,
  `FuelDeliveryTypeID` int(10) NOT NULL,
  `FuelDeliverySubTypeID` int(10) NOT NULL,
  `FuelSystemControlTypeID` int(10) NOT NULL,
  `FuelSystemDesignID` int(10) NOT NULL,
  PRIMARY KEY (`FuelDeliveryConfigID`),
  KEY `IDX_FuelDeliveryConfig_FuelDel1` (`FuelDeliverySubTypeID`),
  KEY `IDX_FuelDeliveryConfig_FuelDel2` (`FuelDeliveryTypeID`),
  KEY `IDX_FuelDeliveryConfig_FuelSys3` (`FuelSystemControlTypeID`),
  KEY `IDX_FuelDeliveryConfig_FuelSys4` (`FuelSystemDesignID`),
  CONSTRAINT `FK_FuelDeliverySubType_FuelDeliveryConfig` FOREIGN KEY (`FuelDeliverySubTypeID`) REFERENCES `fueldeliverysubtype` (`FuelDeliverySubTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_FuelDeliveryType_FuelDeliveryConfig` FOREIGN KEY (`FuelDeliveryTypeID`) REFERENCES `fueldeliverytype` (`FuelDeliveryTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_FuelSystemControlType_FuelDeliveryConfig` FOREIGN KEY (`FuelSystemControlTypeID`) REFERENCES `fuelsystemcontroltype` (`FuelSystemControlTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_FuelSystemDesign_FuelDeliveryConfig` FOREIGN KEY (`FuelSystemDesignID`) REFERENCES `fuelsystemdesign` (`FuelSystemDesignID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fueldeliverysubtypefueldeliv_fk` FOREIGN KEY (`FuelDeliverySubTypeID`) REFERENCES `fueldeliverysubtype` (`FuelDeliverySubTypeID`),
  CONSTRAINT `fueldeliverytypefueldelivery_fk` FOREIGN KEY (`FuelDeliveryTypeID`) REFERENCES `fueldeliverytype` (`FuelDeliveryTypeID`),
  CONSTRAINT `fuelsystemcontroltypefueldel_fk` FOREIGN KEY (`FuelSystemControlTypeID`) REFERENCES `fuelsystemcontroltype` (`FuelSystemControlTypeID`),
  CONSTRAINT `fuelsystemdesignfueldelivery_fk` FOREIGN KEY (`FuelSystemDesignID`) REFERENCES `fuelsystemdesign` (`FuelSystemDesignID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `FuelDeliverySubType`
-- ----------------------------
DROP TABLE IF EXISTS `FuelDeliverySubType`;
CREATE TABLE `FuelDeliverySubType` (
  `FuelDeliverySubTypeID` int(10) NOT NULL,
  `FuelDeliverySubTypeName` varchar(50) NOT NULL,
  PRIMARY KEY (`FuelDeliverySubTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `FuelDeliveryType`
-- ----------------------------
DROP TABLE IF EXISTS `FuelDeliveryType`;
CREATE TABLE `FuelDeliveryType` (
  `FuelDeliveryTypeID` int(10) NOT NULL,
  `FuelDeliveryTypeName` varchar(50) NOT NULL,
  PRIMARY KEY (`FuelDeliveryTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `FuelSystemControlType`
-- ----------------------------
DROP TABLE IF EXISTS `FuelSystemControlType`;
CREATE TABLE `FuelSystemControlType` (
  `FuelSystemControlTypeID` int(10) NOT NULL,
  `FuelSystemControlTypeName` varchar(50) NOT NULL,
  PRIMARY KEY (`FuelSystemControlTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `FuelSystemDesign`
-- ----------------------------
DROP TABLE IF EXISTS `FuelSystemDesign`;
CREATE TABLE `FuelSystemDesign` (
  `FuelSystemDesignID` int(10) NOT NULL,
  `FuelSystemDesignName` varchar(50) NOT NULL,
  PRIMARY KEY (`FuelSystemDesignID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `FuelType`
-- ----------------------------
DROP TABLE IF EXISTS `FuelType`;
CREATE TABLE `FuelType` (
  `FuelTypeID` int(10) NOT NULL,
  `FuelTypeName` varchar(30) NOT NULL,
  PRIMARY KEY (`FuelTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `IgnitionSystemType`
-- ----------------------------
DROP TABLE IF EXISTS `IgnitionSystemType`;
CREATE TABLE `IgnitionSystemType` (
  `IgnitionSystemTypeID` int(10) NOT NULL,
  `IgnitionSystemTypeName` varchar(30) NOT NULL,
  PRIMARY KEY (`IgnitionSystemTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Language`
-- ----------------------------
DROP TABLE IF EXISTS `Language`;
CREATE TABLE `Language` (
  `LanguageID` int(10) NOT NULL AUTO_INCREMENT,
  `LanguageName` varchar(20) NOT NULL,
  `DialectName` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`LanguageID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `LanguageTranslation`
-- ----------------------------
DROP TABLE IF EXISTS `LanguageTranslation`;
CREATE TABLE `LanguageTranslation` (
  `LanguageTranslationID` int(10) NOT NULL AUTO_INCREMENT,
  `EnglishPhraseID` int(10) NOT NULL,
  `LanguageID` int(10) NOT NULL,
  `Translation` varchar(150) NOT NULL,
  PRIMARY KEY (`LanguageTranslationID`),
  KEY `IX_LanguageTranslation_EnglishP` (`EnglishPhraseID`),
  KEY `IX_LanguageTranslation_Language` (`LanguageID`),
  CONSTRAINT `englishphraselanguagetranslation_fk` FOREIGN KEY (`EnglishPhraseID`) REFERENCES `englishphrase` (`EnglishPhraseID`),
  CONSTRAINT `FK_EnglishPhrase_LanguageTranslation` FOREIGN KEY (`EnglishPhraseID`) REFERENCES `englishphrase` (`EnglishPhraseID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_Language_LanguageTranslation` FOREIGN KEY (`LanguageID`) REFERENCES `language` (`LanguageID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `languagelanguagetranslation_fk` FOREIGN KEY (`LanguageID`) REFERENCES `language` (`LanguageID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `LanguageTranslationAttachment`
-- ----------------------------
DROP TABLE IF EXISTS `LanguageTranslationAttachment`;
CREATE TABLE `LanguageTranslationAttachment` (
  `LanguageTranslationAttachmentID` int(10) NOT NULL AUTO_INCREMENT,
  `LanguageTranslationID` int(10) NOT NULL,
  `AttachmentID` int(10) NOT NULL,
  PRIMARY KEY (`LanguageTranslationAttachmentID`),
  KEY `IX_LanguageTranslationAttachme1` (`AttachmentID`),
  KEY `IX_LanguageTranslationAttachme2` (`LanguageTranslationID`),
  CONSTRAINT `attachmentlanguagetranslationattachment_fk` FOREIGN KEY (`AttachmentID`) REFERENCES `Attachment` (`AttachmentID`),
  CONSTRAINT `FK_Attachment_LanguageTranslationAttachment` FOREIGN KEY (`AttachmentID`) REFERENCES `Attachment` (`AttachmentID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_LanguageTranslation_LanguageTranslationAttachment` FOREIGN KEY (`LanguageTranslationID`) REFERENCES `LanguageTranslation` (`LanguageTranslationID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `languagetranslationlanguagetranslationattachment_fk` FOREIGN KEY (`LanguageTranslationID`) REFERENCES `LanguageTranslation` (`LanguageTranslationID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Make`
-- ----------------------------
DROP TABLE IF EXISTS `Make`;
CREATE TABLE `Make` (
  `MakeID` int(10) NOT NULL,
  `MakeName` varchar(50) NOT NULL,
  PRIMARY KEY (`MakeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Mfr`
-- ----------------------------
DROP TABLE IF EXISTS `Mfr`;
CREATE TABLE `Mfr` (
  `MfrID` int(10) NOT NULL,
  `MfrName` varchar(30) NOT NULL,
  PRIMARY KEY (`MfrID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `MfrBodyCode`
-- ----------------------------
DROP TABLE IF EXISTS `MfrBodyCode`;
CREATE TABLE `MfrBodyCode` (
  `MfrBodyCodeID` int(10) NOT NULL,
  `MfrBodyCodeName` varchar(10) NOT NULL,
  PRIMARY KEY (`MfrBodyCodeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Model`
-- ----------------------------
DROP TABLE IF EXISTS `Model`;
CREATE TABLE `Model` (
  `ModelID` int(10) NOT NULL,
  `ModelName` varchar(100) DEFAULT NULL,
  `VehicleTypeID` int(10) NOT NULL,
  PRIMARY KEY (`ModelID`),
  KEY `IDX_Model_VehicleTypeID` (`VehicleTypeID`),
  CONSTRAINT `FK_VehicleType_Model` FOREIGN KEY (`VehicleTypeID`) REFERENCES `vehicletype` (`VehicleTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `vehicletypemodel_fk` FOREIGN KEY (`VehicleTypeID`) REFERENCES `vehicletype` (`VehicleTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `PowerOutput`
-- ----------------------------
DROP TABLE IF EXISTS `PowerOutput`;
CREATE TABLE `PowerOutput` (
  `PowerOutputID` int(10) NOT NULL,
  `HorsePower` varchar(10) NOT NULL,
  `KilowattPower` varchar(10) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `PublicationStage`
-- ----------------------------
DROP TABLE IF EXISTS `PublicationStage`;
CREATE TABLE `PublicationStage` (
  `PublicationStageID` int(10) NOT NULL,
  `PublicationStageName` varchar(100) NOT NULL,
  PRIMARY KEY (`PublicationStageID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Region`
-- ----------------------------
DROP TABLE IF EXISTS `Region`;
CREATE TABLE `Region` (
  `RegionID` int(10) NOT NULL,
  `ParentID` int(10) DEFAULT NULL,
  `RegionAbbr` char(3) NOT NULL,
  `RegionName` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`RegionID`),
  KEY `IDX_Region_ParentID` (`ParentID`),
  CONSTRAINT `FK_Region_Parent` FOREIGN KEY (`ParentID`) REFERENCES `region` (`RegionID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `SpringType`
-- ----------------------------
DROP TABLE IF EXISTS `SpringType`;
CREATE TABLE `SpringType` (
  `SpringTypeID` int(10) NOT NULL,
  `SpringTypeName` varchar(50) NOT NULL,
  PRIMARY KEY (`SpringTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `SpringTypeConfig`
-- ----------------------------
DROP TABLE IF EXISTS `SpringTypeConfig`;
CREATE TABLE `SpringTypeConfig` (
  `SpringTypeConfigID` int(10) NOT NULL,
  `FrontSpringTypeID` int(10) NOT NULL,
  `RearSpringTypeID` int(10) NOT NULL,
  PRIMARY KEY (`SpringTypeConfigID`),
  KEY `IDX_SpringTypeConfig_FrontSprin` (`FrontSpringTypeID`),
  KEY `IDX_SpringTypeConfig_RearSpring` (`RearSpringTypeID`),
  CONSTRAINT `FK_FrontSpringType_SpringTypeConfig` FOREIGN KEY (`FrontSpringTypeID`) REFERENCES `springtype` (`SpringTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_RearSpringType_SpringTypeConfig` FOREIGN KEY (`RearSpringTypeID`) REFERENCES `springtype` (`SpringTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `springconfig1_fk` FOREIGN KEY (`FrontSpringTypeID`) REFERENCES `springtype` (`SpringTypeID`),
  CONSTRAINT `springconfig2_fk` FOREIGN KEY (`RearSpringTypeID`) REFERENCES `springtype` (`SpringTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `SteeringConfig`
-- ----------------------------
DROP TABLE IF EXISTS `SteeringConfig`;
CREATE TABLE `SteeringConfig` (
  `SteeringConfigID` int(10) NOT NULL,
  `SteeringTypeID` int(10) NOT NULL,
  `SteeringSystemID` int(10) NOT NULL,
  PRIMARY KEY (`SteeringConfigID`),
  KEY `IDX_SteeringConfig_SteeringSyst` (`SteeringSystemID`),
  KEY `IDX_SteeringConfig_SteeringType` (`SteeringTypeID`),
  CONSTRAINT `FK_SteeringSystem_SteeringConfig` FOREIGN KEY (`SteeringSystemID`) REFERENCES `steeringsystem` (`SteeringSystemID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_SteeringType_SteeringConfig` FOREIGN KEY (`SteeringTypeID`) REFERENCES `steeringtype` (`SteeringTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `steeringsystemsteeringconfig_fk` FOREIGN KEY (`SteeringSystemID`) REFERENCES `steeringsystem` (`SteeringSystemID`),
  CONSTRAINT `steeringtypesteeringconfig_fk` FOREIGN KEY (`SteeringTypeID`) REFERENCES `steeringtype` (`SteeringTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `SteeringSystem`
-- ----------------------------
DROP TABLE IF EXISTS `SteeringSystem`;
CREATE TABLE `SteeringSystem` (
  `SteeringSystemID` int(10) NOT NULL,
  `SteeringSystemName` varchar(30) NOT NULL,
  PRIMARY KEY (`SteeringSystemID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `SteeringType`
-- ----------------------------
DROP TABLE IF EXISTS `SteeringType`;
CREATE TABLE `SteeringType` (
  `SteeringTypeID` int(10) NOT NULL,
  `SteeringTypeName` varchar(30) NOT NULL,
  PRIMARY KEY (`SteeringTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Submodel`
-- ----------------------------
DROP TABLE IF EXISTS `Submodel`;
CREATE TABLE `Submodel` (
  `SubmodelID` int(10) NOT NULL,
  `SubmodelName` varchar(50) NOT NULL,
  PRIMARY KEY (`SubmodelID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Transmission`
-- ----------------------------
DROP TABLE IF EXISTS `Transmission`;
CREATE TABLE `Transmission` (
  `TransmissionID` int(10) NOT NULL,
  `TransmissionBaseID` int(10) NOT NULL,
  `TransmissionMfrCodeID` int(10) NOT NULL,
  `TransmissionElecControlledID` int(10) NOT NULL,
  `TransmissionMfrID` int(10) NOT NULL,
  PRIMARY KEY (`TransmissionID`),
  KEY `IDX_Transmission_TransmissionBa` (`TransmissionBaseID`),
  KEY `IDX_Transmission_TransmissionM1` (`TransmissionMfrCodeID`),
  KEY `IDX_Transmission_TransmissionM2` (`TransmissionMfrID`),
  KEY `FK_Transmission_ElecControlled` (`TransmissionElecControlledID`),
  CONSTRAINT `FK_Mfr_Transmission` FOREIGN KEY (`TransmissionMfrID`) REFERENCES `mfr` (`MfrID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_TransmissionBase_Transmission` FOREIGN KEY (`TransmissionBaseID`) REFERENCES `TransmissionBase` (`TransmissionBaseID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_TransmissionMfrCode_Transmission` FOREIGN KEY (`TransmissionMfrCodeID`) REFERENCES `transmissionmfrcode` (`TransmissionMfrCodeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_Transmission_ElecControlled` FOREIGN KEY (`TransmissionElecControlledID`) REFERENCES `eleccontrolled` (`ElecControlledID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `mfrtransmission_fk` FOREIGN KEY (`TransmissionMfrID`) REFERENCES `mfr` (`MfrID`),
  CONSTRAINT `transmissionbasetransmission_fk` FOREIGN KEY (`TransmissionBaseID`) REFERENCES `TransmissionBase` (`TransmissionBaseID`),
  CONSTRAINT `transmissionmfrcodetransmiss_fk` FOREIGN KEY (`TransmissionMfrCodeID`) REFERENCES `transmissionmfrcode` (`TransmissionMfrCodeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `TransmissionBase`
-- ----------------------------
DROP TABLE IF EXISTS `TransmissionBase`;
CREATE TABLE `TransmissionBase` (
  `TransmissionBaseID` int(10) NOT NULL,
  `TransmissionTypeID` int(10) NOT NULL,
  `TransmissionNumSpeedsID` int(10) NOT NULL,
  `TransmissionControlTypeID` int(10) NOT NULL,
  PRIMARY KEY (`TransmissionBaseID`),
  KEY `IDX_TransmissionBase_Transmiss1` (`TransmissionControlTypeID`),
  KEY `IDX_TransmissionBase_Transmiss2` (`TransmissionNumSpeedsID`),
  KEY `IDX_TransmissionBase_Transmiss3` (`TransmissionTypeID`),
  CONSTRAINT `FK_TransmissionControlType_TransmissionBase` FOREIGN KEY (`TransmissionControlTypeID`) REFERENCES `transmissioncontroltype` (`TransmissionControlTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_TransmissionNumSpeeds_TransmissionBase` FOREIGN KEY (`TransmissionNumSpeedsID`) REFERENCES `transmissionnumspeeds` (`TransmissionNumSpeedsID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_TransmissionType_TransmissionBase` FOREIGN KEY (`TransmissionTypeID`) REFERENCES `transmissiontype` (`TransmissionTypeID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `transmissioncontroltypetrans_fk` FOREIGN KEY (`TransmissionControlTypeID`) REFERENCES `transmissioncontroltype` (`TransmissionControlTypeID`),
  CONSTRAINT `transmissionnumspeedstransmi_fk` FOREIGN KEY (`TransmissionNumSpeedsID`) REFERENCES `transmissionnumspeeds` (`TransmissionNumSpeedsID`),
  CONSTRAINT `transmissiontypetransmission_fk` FOREIGN KEY (`TransmissionTypeID`) REFERENCES `transmissiontype` (`TransmissionTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `TransmissionControlType`
-- ----------------------------
DROP TABLE IF EXISTS `TransmissionControlType`;
CREATE TABLE `TransmissionControlType` (
  `TransmissionControlTypeID` int(10) NOT NULL,
  `TransmissionControlTypeName` varchar(30) NOT NULL,
  PRIMARY KEY (`TransmissionControlTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `TransmissionMfrCode`
-- ----------------------------
DROP TABLE IF EXISTS `TransmissionMfrCode`;
CREATE TABLE `TransmissionMfrCode` (
  `TransmissionMfrCodeID` int(10) NOT NULL,
  `TransmissionMfrCode` varchar(30) NOT NULL,
  PRIMARY KEY (`TransmissionMfrCodeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `TransmissionNumSpeeds`
-- ----------------------------
DROP TABLE IF EXISTS `TransmissionNumSpeeds`;
CREATE TABLE `TransmissionNumSpeeds` (
  `TransmissionNumSpeedsID` int(10) NOT NULL,
  `TransmissionNumSpeeds` char(3) NOT NULL,
  PRIMARY KEY (`TransmissionNumSpeedsID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `TransmissionType`
-- ----------------------------
DROP TABLE IF EXISTS `TransmissionType`;
CREATE TABLE `TransmissionType` (
  `TransmissionTypeID` int(10) NOT NULL,
  `TransmissionTypeName` varchar(30) NOT NULL,
  PRIMARY KEY (`TransmissionTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Valves`
-- ----------------------------
DROP TABLE IF EXISTS `Valves`;
CREATE TABLE `Valves` (
  `ValvesID` int(10) NOT NULL,
  `ValvesPerEngine` char(3) NOT NULL,
  PRIMARY KEY (`ValvesID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Vehicle`
-- ----------------------------
DROP TABLE IF EXISTS `Vehicle`;
CREATE TABLE `Vehicle` (
  `VehicleID` int(10) NOT NULL,
  `BaseVehicleID` int(10) NOT NULL,
  `SubmodelID` int(10) NOT NULL,
  `RegionID` int(10) NOT NULL,
  `Source` char(10) DEFAULT NULL,
  `PublicationStageID` int(10) NOT NULL DEFAULT '4',
  `PublicationStageSource` varchar(100) NOT NULL,
  `PublicationStageDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`VehicleID`),
  KEY `IDX_Vehicle_BaseVehicleID` (`BaseVehicleID`),
  KEY `IDX_Vehicle_PublicationStage` (`PublicationStageID`),
  KEY `IDX_Vehicle_RegionID` (`RegionID`),
  KEY `IDX_Vehicle_SubmodelID` (`SubmodelID`),
  CONSTRAINT `basevehiclevehicle_fk` FOREIGN KEY (`BaseVehicleID`) REFERENCES `BaseVehicle` (`BaseVehicleID`),
  CONSTRAINT `FK_BaseVehicle_Vehicle` FOREIGN KEY (`BaseVehicleID`) REFERENCES `BaseVehicle` (`BaseVehicleID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_PublicationStage_Vehicle` FOREIGN KEY (`PublicationStageID`) REFERENCES `publicationstage` (`PublicationStageID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_Region_Vehicle` FOREIGN KEY (`RegionID`) REFERENCES `region` (`RegionID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_SubModel_Vehicle` FOREIGN KEY (`SubmodelID`) REFERENCES `submodel` (`SubmodelID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `publicationstage_fk` FOREIGN KEY (`PublicationStageID`) REFERENCES `publicationstage` (`PublicationStageID`),
  CONSTRAINT `regionvehicle_fk` FOREIGN KEY (`RegionID`) REFERENCES `region` (`RegionID`),
  CONSTRAINT `submodelvehicle_fk` FOREIGN KEY (`SubmodelID`) REFERENCES `submodel` (`SubmodelID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `VehicleToBedConfig`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleToBedConfig`;
CREATE TABLE `VehicleToBedConfig` (
  `VehicleToBedConfigID` int(10) NOT NULL,
  `VehicleID` int(10) NOT NULL,
  `BedConfigID` int(10) NOT NULL,
  `Source` char(10) DEFAULT NULL,
  PRIMARY KEY (`VehicleToBedConfigID`),
  KEY `IDX_VehicleToBedConfig_BedConfi` (`BedConfigID`),
  KEY `IDX_VehicleToBedConfig_VehicleID` (`VehicleID`),
  CONSTRAINT `bedconfigvehicletobed_fk` FOREIGN KEY (`BedConfigID`) REFERENCES `BedConfig` (`BedConfigID`),
  CONSTRAINT `vehiclevehicle_fk` FOREIGN KEY (`VehicleID`) REFERENCES `Vehicle` (`VehicleID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `VehicleToBodyStyleConfig`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleToBodyStyleConfig`;
CREATE TABLE `VehicleToBodyStyleConfig` (
  `VehicleToBodyStyleConfigID` int(10) NOT NULL,
  `VehicleID` int(10) NOT NULL,
  `BodyStyleConfigID` int(10) NOT NULL,
  `Source` char(10) DEFAULT NULL,
  PRIMARY KEY (`VehicleToBodyStyleConfigID`),
  KEY `IDX_VehicleToBodyStyleConfig_Bo` (`BodyStyleConfigID`),
  KEY `IDX_VehicleToBodyStyleConfig_Ve` (`VehicleID`),
  CONSTRAINT `bodystyleconfigbasetosubmode_fk` FOREIGN KEY (`BodyStyleConfigID`) REFERENCES `BodyStyleConfig` (`BodyStyleConfigID`),
  CONSTRAINT `vehicletobodystyleconfigvehicle_fk` FOREIGN KEY (`VehicleID`) REFERENCES `Vehicle` (`VehicleID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `VehicleToBrakeConfig`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleToBrakeConfig`;
CREATE TABLE `VehicleToBrakeConfig` (
  `VehicleToBrakeConfigID` int(10) NOT NULL,
  `VehicleID` int(10) NOT NULL,
  `BrakeConfigID` int(10) NOT NULL,
  `Source` char(10) DEFAULT NULL,
  PRIMARY KEY (`VehicleToBrakeConfigID`),
  KEY `IDX_VehicleToBrakeConfig_BrakeC` (`BrakeConfigID`),
  KEY `IDX_VehicleToBrakeConfig_Vehicl` (`VehicleID`),
  CONSTRAINT `brakeconfigvehicle_fk` FOREIGN KEY (`BrakeConfigID`) REFERENCES `BrakeConfig` (`BrakeConfigID`),
  CONSTRAINT `vehicletobrakeconfigvehicle_fk` FOREIGN KEY (`VehicleID`) REFERENCES `Vehicle` (`VehicleID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `VehicleToDriveType`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleToDriveType`;
CREATE TABLE `VehicleToDriveType` (
  `VehicleToDriveTypeID` int(10) NOT NULL,
  `VehicleID` int(10) NOT NULL,
  `DriveTypeID` int(10) NOT NULL,
  `Source` char(10) DEFAULT NULL,
  PRIMARY KEY (`VehicleToDriveTypeID`),
  KEY `IDX_VehicleToDriveType_DriveTyp` (`DriveTypeID`),
  KEY `IDX_VehicleToDriveType_VehicleID` (`VehicleID`),
  CONSTRAINT `drivetypevehicletodri_fk` FOREIGN KEY (`DriveTypeID`) REFERENCES `drivetype` (`DriveTypeID`),
  CONSTRAINT `vehicletodrivetypevehicle_fk` FOREIGN KEY (`VehicleID`) REFERENCES `Vehicle` (`VehicleID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `VehicleToEngineConfig`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleToEngineConfig`;
CREATE TABLE `VehicleToEngineConfig` (
  `VehicleToEngineConfigID` int(10) NOT NULL,
  `VehicleID` int(10) NOT NULL,
  `EngineConfigID` int(10) NOT NULL,
  `Source` char(10) DEFAULT NULL,
  PRIMARY KEY (`VehicleToEngineConfigID`),
  KEY `IDX_VehicleToEngineConfig_Engin` (`EngineConfigID`),
  KEY `IDX_VehicleToEngineConfig_Vehic` (`VehicleID`),
  CONSTRAINT `engineconfigvehicleto_fk` FOREIGN KEY (`EngineConfigID`) REFERENCES `EngineConfig` (`EngineConfigID`),
  CONSTRAINT `vehicletoengineconfigvehicle_fk` FOREIGN KEY (`VehicleID`) REFERENCES `Vehicle` (`VehicleID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `VehicleToMfrBodyCode`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleToMfrBodyCode`;
CREATE TABLE `VehicleToMfrBodyCode` (
  `VehicleToMfrBodyCodeID` int(10) NOT NULL,
  `VehicleID` int(10) NOT NULL,
  `MfrBodyCodeID` int(10) NOT NULL,
  `Source` char(10) DEFAULT NULL,
  PRIMARY KEY (`VehicleToMfrBodyCodeID`),
  KEY `IDX_VehicleToMfrBodyCode_MfrBod` (`MfrBodyCodeID`),
  KEY `IDX_VehicleToMfrBodyCode_Vehicl` (`VehicleID`),
  CONSTRAINT `mfrbodycodevehicletom_fk` FOREIGN KEY (`MfrBodyCodeID`) REFERENCES `mfrbodycode` (`MfrBodyCodeID`),
  CONSTRAINT `vehicletomfrbodycodevehicle_fk` FOREIGN KEY (`VehicleID`) REFERENCES `Vehicle` (`VehicleID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `VehicleToSpringTypeConfig`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleToSpringTypeConfig`;
CREATE TABLE `VehicleToSpringTypeConfig` (
  `VehicleToSpringTypeConfigID` int(10) NOT NULL,
  `VehicleID` int(10) NOT NULL,
  `SpringTypeConfigID` int(10) NOT NULL,
  `Source` char(10) DEFAULT NULL,
  PRIMARY KEY (`VehicleToSpringTypeConfigID`),
  KEY `IDX_VehicleToSpringTypeConfig_S` (`SpringTypeConfigID`),
  KEY `IDX_VehicleToSpringTypeConfig_V` (`VehicleID`),
  CONSTRAINT `vehicletospringtype1_fk` FOREIGN KEY (`SpringTypeConfigID`) REFERENCES `SpringTypeConfig` (`SpringTypeConfigID`),
  CONSTRAINT `vehicletospringtypeconfigvehicle_fk` FOREIGN KEY (`VehicleID`) REFERENCES `Vehicle` (`VehicleID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `VehicleToSteeringConfig`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleToSteeringConfig`;
CREATE TABLE `VehicleToSteeringConfig` (
  `VehicleToSteeringConfigID` int(10) NOT NULL,
  `VehicleID` int(10) NOT NULL,
  `SteeringConfigID` int(10) NOT NULL,
  `Source` char(10) DEFAULT NULL,
  PRIMARY KEY (`VehicleToSteeringConfigID`),
  KEY `IDX_VehicleToSteeringConfig_Ste` (`SteeringConfigID`),
  KEY `IDX_VehicleToSteeringConfig_Veh` (`VehicleID`),
  CONSTRAINT `steeringconfigvehicle_fk` FOREIGN KEY (`SteeringConfigID`) REFERENCES `SteeringConfig` (`SteeringConfigID`),
  CONSTRAINT `vehicletosteeringconfigvehicle_fk` FOREIGN KEY (`VehicleID`) REFERENCES `Vehicle` (`VehicleID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `VehicleToTransmission`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleToTransmission`;
CREATE TABLE `VehicleToTransmission` (
  `VehicleToTransmissionID` int(10) NOT NULL,
  `VehicleID` int(10) NOT NULL,
  `TransmissionID` int(10) NOT NULL,
  `Source` char(10) DEFAULT NULL,
  PRIMARY KEY (`VehicleToTransmissionID`),
  KEY `IDX_VehicleToTransmission_Trans` (`TransmissionID`),
  KEY `IDX_VehicleToTransmission_Vehic` (`VehicleID`),
  CONSTRAINT `transmissionvehicleto_fk` FOREIGN KEY (`TransmissionID`) REFERENCES `Transmission` (`TransmissionID`),
  CONSTRAINT `vehicletotransmissionvehicle_fk` FOREIGN KEY (`VehicleID`) REFERENCES `Vehicle` (`VehicleID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `VehicleToWheelbase`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleToWheelbase`;
CREATE TABLE `VehicleToWheelbase` (
  `VehicleToWheelbaseID` int(10) NOT NULL,
  `VehicleID` int(10) NOT NULL,
  `WheelbaseID` int(10) NOT NULL,
  `Source` char(10) DEFAULT NULL,
  PRIMARY KEY (`VehicleToWheelbaseID`),
  KEY `IDX_VehicleToWheelbase_VehicleID` (`VehicleID`),
  KEY `IDX_VehicleToWheelbase_Wheelbas` (`WheelbaseID`),
  CONSTRAINT `vehicletowheelbasevehicle_fk` FOREIGN KEY (`VehicleID`) REFERENCES `Vehicle` (`VehicleID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `VehicleType`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleType`;
CREATE TABLE `VehicleType` (
  `VehicleTypeID` int(10) NOT NULL,
  `VehicleTypeName` varchar(50) NOT NULL,
  `VehicleTypeGroupID` int(10) DEFAULT NULL,
  PRIMARY KEY (`VehicleTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `VehicleTypeGroup`
-- ----------------------------
DROP TABLE IF EXISTS `VehicleTypeGroup`;
CREATE TABLE `VehicleTypeGroup` (
  `VehicleTypeGroupID` int(10) NOT NULL,
  `VehicleTypeGroupName` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Version`
-- ----------------------------
DROP TABLE IF EXISTS `Version`;
CREATE TABLE `Version` (
  `VersionDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `WheelBase`
-- ----------------------------
DROP TABLE IF EXISTS `WheelBase`;
CREATE TABLE `WheelBase` (
  `WheelBaseID` int(10) NOT NULL,
  `WheelBase` varchar(10) NOT NULL,
  `WheelBaseMetric` varchar(10) NOT NULL,
  PRIMARY KEY (`WheelBaseID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `Year`
-- ----------------------------
DROP TABLE IF EXISTS `Year`;
CREATE TABLE `Year` (
  `YearID` int(10) NOT NULL,
  PRIMARY KEY (`YearID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
--  Table structure for `vcdbchanges`
-- ----------------------------
DROP TABLE IF EXISTS `vcdbchanges`;
CREATE TABLE `vcdbchanges` (
  `versiondate` datetime NOT NULL,
  `tablename` varchar(30) NOT NULL,
  `id` int(10) NOT NULL,
  `action` char(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

SET FOREIGN_KEY_CHECKS = 1;
