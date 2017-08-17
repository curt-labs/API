CREATE DATABASE CurtData;

-- Brands --
CREATE TABLE Brand (
  ID           INT AUTO_INCREMENT PRIMARY KEY,
  name         VARCHAR(255) NOT NULL,
  code         VARCHAR(255) NOT NULL,
  logo         VARCHAR(255) NULL,
  logoAlt      VARCHAR(255) NULL,
  formalName   VARCHAR(255) NULL,
  longName     VARCHAR(255) NULL,
  primaryColor VARCHAR(10)  NULL,
  autocareID   VARCHAR(4)   NULL
);

INSERT INTO CurtData.Brand (ID, name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID) VALUES (1,'CURT',    'CURT',    'http://www.curtmfg.com/Content/img/logo.png', 'https://storage.googleapis.com/curt-logos/logo.png',               'CURT Manufacturing, LLC', 'CURT Manufacturing', '#e64d2c', 'BKDK');
INSERT INTO CurtData.Brand (ID, name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID) VALUES (3, 'ARIES',   'ARIES',   'https://storage.googleapis.com/aries-logo/SVG_Logo%20(2c_white%20with%20black%20outline%20on%20transparent).svg', 'https://storage.googleapis.com/aries-logo/ARIES%20Logo%20(1c_red%20on%20transparent).png', 'Aries Automotive', 'Aries Automotive', '#57111A', 'BBRD');
INSERT INTO CurtData.Brand (ID, name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID) VALUES (4, 'Luverne', 'Luverne', null,                                                                                                              null, 'Luverne Truck', 'Luverne Truck Equipment', null, 'FTNF');
INSERT INTO CurtData.Brand (ID, name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID) VALUES (5, 'Retrac',  'Retrac',  'https://storage.googleapis.com/curt-groups/Brand-Logos%20-%20RETRAC.png',                                         null, 'Retrac Mirrors', 'Retrac Mirrors', null, 'HCSF');
INSERT INTO CurtData.Brand (ID, name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID) VALUES (6, 'UWS',     'UWS',     'https://storage.googleapis.com/curt-groups/Brand-Logos%20-%20UWS.png',                                            null, 'UWS', 'UWS', null, 'BHSG');




-- API Key Types --

CREATE TABLE ApiKeyType (
  id         VARCHAR(64)                         NOT NULL PRIMARY KEY,
  type       VARCHAR(500)                        NULL,
  date_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  CONSTRAINT id UNIQUE (id)
);

-- Mock API Key Types --
INSERT INTO CurtData.ApiKeyType (id, type, date_added) VALUES ('EA181F86-3F74-4AD6-8884-829B4558B99D', 'Authentication', '2013-01-02 09:53:21');
INSERT INTO CurtData.ApiKeyType (id, type, date_added) VALUES ('CCDF2BD3-3123-4E54-9E45-7932BAFC8B4D', 'Custom',         '2013-01-14 07:54:50');
INSERT INTO CurtData.ApiKeyType (id, type, date_added) VALUES ('92ff1833-2ca6-11e4-8758-42010af0fd79', 'Internal',       '2015-07-28 05:33:27');
INSERT INTO CurtData.ApiKeyType (id, type, date_added) VALUES ('2922D5BF-6F81-4E9F-9910-C72426F728A1', 'Private',        '2013-01-02 09:53:21');
INSERT INTO CurtData.ApiKeyType (id, type, date_added) VALUES ('209A05AD-7D42-4C88-B5FA-FEEACDD19AC2', 'Public',         '2013-01-02 09:53:21');



-- Customer Users --
CREATE TABLE CustomerUser (
  id                VARCHAR(64)                         NOT NULL PRIMARY KEY,
  name              VARCHAR(255)                        NULL,
  email             VARCHAR(255)                        NOT NULL,
  password          VARCHAR(255)                        NOT NULL,
  customerID        INT                                 NULL,
  date_added        TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  active            TINYINT(1) DEFAULT '0'              NOT NULL,
  locationID        INT DEFAULT '0'                     NOT NULL,
  isSudo            TINYINT(1) DEFAULT '0'              NOT NULL,
  cust_ID           INT                                 NOT NULL,
  NotCustomer       TINYINT(1)                          NULL,
  passwordConverted TINYINT(1)                          NOT NULL,
  CONSTRAINT id UNIQUE (id)
);

-- Mock Customer Users --
-- No Brands --
INSERT INTO CustomerUser (
  id,
  name,
  email,
  password,
  customerID,
  date_added,
  active,
  locationID,
  isSudo,
  cust_ID,
  NotCustomer,
  passwordConverted
)
VALUES (
  '100000000000-0000-4000-1000-00000001',
  'Example Customer User 1',
  'example@example.com',
  'should be salted',
  10000001,
  '2013-01-02 09:53:21',
  1,
  1,
  1,
  11000001,
  0,
  1
);

-- One Brand --
INSERT INTO CustomerUser (
  id,
  name,
  email,
  password,
  customerID,
  date_added,
  active,
  locationID,
  isSudo,
  cust_ID,
  NotCustomer,
  passwordConverted
)
VALUES (
  '100000000000-0000-4000-1000-00000002',
  'Example Customer User 1',
  'example@example.com',
  'peanuts, like passwords, should be salted',
  10000002,
  '2013-01-02 09:53:21',
  1,
  1,
  1,
  11000002,
  0,
  1
);


CREATE TABLE CustomerToBrand (
  ID         INT AUTO_INCREMENT PRIMARY KEY,
  cust_id    INT                    NOT NULL,
  brandID    INT                    NOT NULL,
  dealerTier INT DEFAULT '0'        NOT NULL,
  dealerType INT DEFAULT '0'        NOT NULL,
  salesRepID INT DEFAULT '0'        NOT NULL,
  websiteID  INT DEFAULT '0'        NOT NULL,
  enableWTB  TINYINT(1) DEFAULT '0' NOT NULL
);




-- API Keys --
CREATE TABLE ApiKey
(
  id         INT AUTO_INCREMENT,
  api_key    VARCHAR(64)                         NOT NULL,
  type_id    VARCHAR(64)                         NOT NULL,
  user_id    VARCHAR(64)                         NOT NULL,
  date_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  CONSTRAINT id
  UNIQUE (id),
  CONSTRAINT FK__ApiKey__type_id__5AEE1AF6
  FOREIGN KEY (type_id) REFERENCES ApiKeyType (id),
  CONSTRAINT FK__ApiKey__user_id__5BE23F2F
  FOREIGN KEY (user_id) REFERENCES CustomerUser (id)
);

-- Mock API Keys --
-- No Brands --
INSERT INTO ApiKey (id, api_key, type_id, user_id) VALUES (
  1, '20000000-0000-4000-1000-000000000001', '209A05AD-7D42-4C88-B5FA-FEEACDD19AC2', '100000000000-0000-4000-1000-00000001'
);

-- One Brand --
INSERT INTO ApiKey (id, api_key, type_id, user_id) VALUES (
  2, '20000000-0000-4000-1000-000000000002', '209A05AD-7D42-4C88-B5FA-FEEACDD19AC2', '100000000000-0000-4000-1000-00000001'
);

-- Multiple Brands --
INSERT INTO ApiKey (id, api_key, type_id, user_id) VALUES (
  3, '20000000-0000-4000-1000-000000000003', '209A05AD-7D42-4C88-B5FA-FEEACDD19AC2', '100000000000-0000-4000-1000-00000001'
);

-- All Brands --
INSERT INTO ApiKey (id, api_key, type_id, user_id) VALUES (
  4, '20000000-0000-4000-1000-000000000004', '209A05AD-7D42-4C88-B5FA-FEEACDD19AC2', '100000000000-0000-4000-1000-00000001'
);


-- API Key Relation to Brands --
CREATE TABLE ApiKeyToBrand (
  ID      INT AUTO_INCREMENT PRIMARY KEY,
  keyID   INT NOT NULL,
  brandID INT NOT NULL,
  CONSTRAINT FK_ApiKeyToBrand_ApiKey FOREIGN KEY (keyID) REFERENCES ApiKey (id),
  CONSTRAINT FK_ApiKeyToBrand_Brand FOREIGN KEY (brandID) REFERENCES Brand (ID)
);

-- Api Key ID = 1 skipped intentionally --
-- One Brand --
INSERT INTO ApiKeyToBrand (ID, keyID, BrandID) VALUES (1, 2, 1);

-- Multiple Brands --
INSERT INTO ApiKeyToBrand (ID, keyID, BrandID) VALUES (2, 3, 1);
INSERT INTO ApiKeyToBrand (ID, keyID, BrandID) VALUES (3, 3, 3);
INSERT INTO ApiKeyToBrand (ID, keyID, BrandID) VALUES (4, 3, 5);

-- All Brands --
INSERT INTO ApiKeyToBrand (ID, keyID, BrandID) VALUES (5, 4, 1);
INSERT INTO ApiKeyToBrand (ID, keyID, BrandID) VALUES (6, 4, 3);
INSERT INTO ApiKeyToBrand (ID, keyID, BrandID) VALUES (7, 4, 4);
INSERT INTO ApiKeyToBrand (ID, keyID, BrandID) VALUES (8, 4, 5);
INSERT INTO ApiKeyToBrand (ID, keyID, BrandID) VALUES (9, 4, 6);
