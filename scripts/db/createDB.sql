CREATE TABLE IF NOT EXISTS USERS (
  UserID CHAR(128) PRIMARY KEY NOT NULL,
  CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  LastUpdate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  Name VARCHAR(512) NOT NULL,
  Description TEXT DEFAULT NULL,
  DeviceID VARCHAR(64) DEFAULT NULL,
  AllowShareData BOOLEAN DEFAULT FALSE,
  Tags TEXT [] DEFAULT NULL,
  Images TEXT [] DEFAULT NULL,
  --Reputation
  Giver REAL DEFAULT 0,
  Taker REAL DEFAULT 0,
  --Contatct
  URL VARCHAR(512) NOT NULL,
  Email VARCHAR(512) DEFAULT NULL,
  Facebook VARCHAR(512) DEFAULT NULL,
  Instagram VARCHAR(512) DEFAULT NULL,
  Google VARCHAR(512) DEFAULT NULL,
  AdditionalData TEXT DEFAULT NULL,
  --Contact Address
  Address TEXT DEFAULT NULL,
  City VARCHAR(64) DEFAULT NULL,
  State VARCHAR(64) DEFAULT NULL,
  ZipCode INTEGER DEFAULT NULL,
  Country VARCHAR(64) DEFAULT NULL,
  Lat NUMERIC(12, 4) DEFAULT NULL,
  Long NUMERIC(12, 4) DEFAULT NULL,
  RegisterFrom TEXT DEFAULT NULL
);
CREATE INDEX IDX_USERS_USERID ON USERS (UserID);
CREATE INDEX IDX_USERS_EMAIL ON USERS (Email);
CREATE TABLE IF NOT EXISTS PHONES (
  PhoneID INTEGER PRIMARY KEY NOT NULL,
  CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UserID CHAR(128) REFERENCES USERS(UserID),
  CountryCode VARCHAR(4) DEFAULT '+55',
  IsDefault BOOLEAN DEFAULT NULL,
  PhoneNumber VARCHAR(13) NOT NULL,
  Region CHAR(2) DEFAULT NULL,
  WhatsApp BOOLEAN DEFAULT NULL
);
CREATE INDEX IDX_PHONES_USERID ON PHONES (UserID);
CREATE UNIQUE INDEX IDX_PHONES_PHONEID_USERID ON PHONES (PhoneID, UserID);
CREATE TABLE IF NOT EXISTS CATEGORIES (
  CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  Name TEXT PRIMARY KEY NOT NULL
);
CREATE INDEX IDX_CATEGORIES_NAME ON CATEGORIES (Name);
CREATE TABLE IF NOT EXISTS PROPOSALS (
  ProposalID CHAR(26) PRIMARY KEY NOT NULL,
  CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  LastUpdate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UserID CHAR(128) REFERENCES USERS(UserID),
  --Side: offer|request|local-business
  Side VARCHAR(64) NOT NULL,
  --ProposalType: job|service|product|finance
  ProposalType VARCHAR(64) NOT NULL,
  --Tags: categories
  Tags TEXT [] DEFAULT NULL,
  Title TEXT NOT NULL,
  Description TEXT NOT NULL,
  ProposalValidate TIMESTAMP DEFAULT NULL,
  Lat NUMERIC(12, 4) DEFAULT NULL,
  Long NUMERIC(12, 4) DEFAULT NULL,
  Range REAL DEFAULT NULL,
  AreaTags TEXT [] DEFAULT NULL,
  IsActive BOOLEAN DEFAULT TRUE,
  Images TEXT [] DEFAULT NULL,
  ExposeUserData BOOLEAN DEFAULT FALSE,
  DataToShare TEXT [] DEFAULT NULL,
  EstimatedValue Money DEFAULT 0
);
CREATE INDEX IDX_PROPOSALS_PROPOSALID ON PROPOSALS (ProposalID);
CREATE INDEX IDX_PROPOSALS_USERID ON PROPOSALS (UserID);
CREATE INDEX IDX_PROPOSALS_SIDE ON PROPOSALS (Side);
CREATE INDEX IDX_PROPOSALS_TYPE ON PROPOSALS (ProposalType);
CREATE TABLE IF NOT EXISTS TRANSACTIONS (
  TransactionID CHAR(26) PRIMARY KEY NOT NULL,
  ProposalID CHAR(26) REFERENCES PROPOSALS(ProposalID),
  GiverID CHAR(128) REFERENCES USERS(UserID),
  TakerID CHAR(128) REFERENCES USERS(UserID),
  CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  LastUpdate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  GiverRating REAL DEFAULT 0,
  GiverReviewComment TEXT DEFAULT NULL,
  TakerRating REAL DEFAULT 0,
  TakerReviewComment TEXT DEFAULT NULL,
  Status VARCHAR(32) NOT NULL
);
CREATE INDEX IDX_TRANSACTIONS_TRANSACTIONID ON TRANSACTIONS (TransactionID);
CREATE INDEX IDX_TRANSACTIONS_PROPOSALID ON TRANSACTIONS (ProposalID);
CREATE INDEX IDX_TRANSACTIONS_GIVERID ON TRANSACTIONS (GiverID);
CREATE INDEX IDX_TRANSACTIONS_TAKERID ON TRANSACTIONS (TakerID);
--NEW
CREATE TABLE IF NOT EXISTS TERMS (
  TermID INTEGER PRIMARY KEY NOT NULL,
  CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  LastUpdate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  Title TEXT NOT NULL,
  Description TEXT NOT NULL,
  IsActive BOOLEAN DEFAULT TRUE
);
CREATE TABLE IF NOT EXISTS TERMS_ACCEPTEDS (
  UserID CHAR(128) PRIMARY KEY REFERENCES USERS(UserID),
  TermID INTEGER REFERENCES TERMS(TermID),
  CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  LastUpdate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  Accepted BOOLEAN DEFAULT TRUE
);
CREATE TABLE IF NOT EXISTS COMPLAINTS (
  ComplaintID INTEGER PRIMARY KEY NOT NULL,
  ComplainerID CHAR(128) REFERENCES USERS(UserID),
  AccusedID CHAR(128) REFERENCES USERS(UserID),
  CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  LastUpdate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  Comment TEXT NOT NULL,
  Accepted BOOLEAN DEFAULT TRUE
);
