-- +goose Up

-- User table to store information about system users
CREATE TABLE UserT (
    UserID SERIAL PRIMARY KEY,
    Username VARCHAR(50) NOT NULL UNIQUE,
    Password VARCHAR(255) NOT NULL,
    Email VARCHAR(255) NOT NULL UNIQUE,
    Role VARCHAR(20) NOT NULL DEFAULT 'User',
    DefaultAdmin BOOLEAN NOT NULL DEFAULT FALSE
);

-- Site table to store information about different locations
CREATE TABLE SiteT (
    SiteID SERIAL PRIMARY KEY,
    SiteName VARCHAR(100) NOT NULL UNIQUE,
    SiteAddress VARCHAR(255),
    SiteMapImagePath VARCHAR(150)
);

-- Building table to store information about buildings within sites
CREATE TABLE BuildingT (
    BuildingID SERIAL PRIMARY KEY,
    SiteID INT NOT NULL,
    BuildingCode VARCHAR(100) NOT NULL,
    FOREIGN KEY (SiteID) REFERENCES SiteT(SiteID)
        ON UPDATE CASCADE  -- If a SiteID changes, update it in BuildingT
        ON DELETE RESTRICT, -- Prevent deletion of a Site if it has associated Buildings
    UNIQUE (SiteID, BuildingCode)  -- Ensure unique building names within each site
);

-- Room table to store information about rooms within buildings
CREATE TABLE RoomT (
    RoomID SERIAL PRIMARY KEY,
    BuildingID INT NOT NULL,
    RoomCode VARCHAR(100) NOT NULL,
    FOREIGN KEY (BuildingID) REFERENCES BuildingT(BuildingID)
        ON UPDATE CASCADE  -- If a BuildingID changes, update it in RoomT
        ON DELETE RESTRICT, -- Prevent deletion of a Building if it has associated Rooms
    UNIQUE (BuildingID, RoomCode)  -- Ensure unique room names within each building
);

-- Emergency Device Type table to categorize different types of emergency devices
CREATE TABLE Emergency_Device_TypeT (
    EmergencyDeviceTypeID SERIAL PRIMARY KEY,
    EmergencyDeviceTypeName VARCHAR(50) NOT NULL UNIQUE
);

-- Extinguisher Type table to categorize different types of fire extinguishers
CREATE TABLE Extinguisher_TypeT (
    ExtinguisherTypeID SERIAL PRIMARY KEY,
    ExtinguisherTypeName VARCHAR(50) NOT NULL UNIQUE
);

-- Emergency Device table to store information about individual emergency devices
CREATE TABLE Emergency_DeviceT (
    EmergencyDeviceID SERIAL PRIMARY KEY,
    EmergencyDeviceTypeID INT NOT NULL,
    RoomID INT NOT NULL,
    ExtinguisherTypeID INT NULL,
    ManufactureDate DATE NULL,
    SerialNumber VARCHAR(50) NULL,
    Description VARCHAR(255) NULL,
    Size VARCHAR(50) NULL,
    LastInspectionDateTime TIMESTAMP NULL,
    Status VARCHAR(50) NULL,
    FOREIGN KEY (EmergencyDeviceTypeID) REFERENCES Emergency_Device_TypeT(EmergencyDeviceTypeID)
        ON UPDATE CASCADE  -- If an EmergencyDeviceTypeID changes, update it in Emergency_DeviceT
        ON DELETE RESTRICT, -- Prevent deletion of an Emergency Device Type if it's associated with any devices
    FOREIGN KEY (RoomID) REFERENCES RoomT(RoomID)
        ON UPDATE CASCADE  -- If a RoomID changes, update it in Emergency_DeviceT
        ON DELETE RESTRICT, -- Prevent deletion of a Room if it has associated Emergency Devices
    FOREIGN KEY (ExtinguisherTypeID) REFERENCES Extinguisher_TypeT(ExtinguisherTypeID)
        ON UPDATE CASCADE  -- If an ExtinguisherTypeID changes, update it in Emergency_DeviceT
        ON DELETE SET NULL -- If an Extinguisher Type is deleted, set the ExtinguisherTypeID to NULL in Emergency_DeviceT
);
 
-- Emergency Device Inspection table to store inspection records for emergency devices
CREATE TABLE Emergency_Device_InspectionT (
    EmergencyDeviceInspectionID SERIAL PRIMARY KEY,
    EmergencyDeviceID INT NOT NULL,
    UserID INT NOT NULL,
    InspectionDateTime TIMESTAMP NOT NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    IsConspicuous BOOLEAN NULL,
    IsAccessible BOOLEAN NULL,
    IsAssignedLocation BOOLEAN NULL,
    IsSignVisible BOOLEAN NULL,
    IsAntiTamperDeviceIntact BOOLEAN NULL,
    IsSupportBracketSecure BOOLEAN NULL,
    AreOperatingInstructionsClear BOOLEAN NULL,
    IsMaintenanceTagAttached BOOLEAN NULL,
    IsNoExternalDamage BOOLEAN NULL,
    IsChargeGaugeNormal BOOLEAN NULL,
    IsReplaced BOOLEAN NULL,
    AreMaintenanceRecordsComplete BOOLEAN NULL,
    WorkOrderRequired BOOLEAN NULL,
    InspectionStatus VARCHAR(20) NOT NULL,
    Notes VARCHAR(255) NULL,
    FOREIGN KEY (EmergencyDeviceID) REFERENCES Emergency_DeviceT(EmergencyDeviceID)
        ON UPDATE CASCADE  -- If an EmergencyDeviceID changes, update it in Emergency_Device_InspectionT
        ON DELETE DELETE, -- Delete associated Inspection records if an Emergency Device is deleted
    FOREIGN KEY (UserID) REFERENCES UserT(UserID)
        ON UPDATE CASCADE  -- If a UserID changes, update it in Emergency_Device_InspectionT
        ON DELETE RESTRICT -- Prevent deletion of a User if they have associated Inspection records
);

-- +goose Down
-- Drop tables in reverse order of creation to avoid foreign key constraint violations
DROP TABLE IF EXISTS Emergency_Device_InspectionT;
DROP TABLE IF EXISTS Emergency_DeviceT;
DROP TABLE IF EXISTS Extinguisher_TypeT; 
DROP TABLE IF EXISTS Emergency_Device_TypeT;
DROP TABLE IF EXISTS RoomT;
DROP TABLE IF EXISTS BuildingT;
DROP TABLE IF EXISTS SiteT;
DROP TABLE IF EXISTS UserT;
DROP TRIGGER IF EXISTS trg_update_device_status ON Emergency_Device_InspectionT;
DROP FUNCTION IF EXISTS update_device_status_on_inspection;