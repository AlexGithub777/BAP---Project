-- First truncate all tables (in correct order due to foreign key constraints)
TRUNCATE TABLE 
    emergency_device_inspectiont,
    emergency_devicet,
    roomt,
    buildingt,
    sitet,
    usert,
    emergency_device_typet,
    extinguisher_typet
CASCADE;

-- Delete all data from the tables
DELETE FROM emergency_device_inspectiont;
DELETE FROM emergency_devicet;
DELETE FROM roomt;
DELETE FROM buildingt;
DELETE FROM sitet;
DELETE FROM usert;
DELETE FROM emergency_device_typet;
DELETE FROM extinguisher_typet;

-- Then reset all sequences
ALTER SEQUENCE buildingt_buildingid_seq RESTART WITH 1;
ALTER SEQUENCE emergency_device_inspectiont_emergencydeviceinspectionid_seq RESTART WITH 1;
ALTER SEQUENCE emergency_device_typet_emergencydevicetypeid_seq RESTART WITH 1;
ALTER SEQUENCE emergency_devicet_emergencydeviceid_seq RESTART WITH 1;
ALTER SEQUENCE extinguisher_typet_extinguishertypeid_seq RESTART WITH 1;
ALTER SEQUENCE roomt_roomid_seq RESTART WITH 1;
ALTER SEQUENCE sitet_siteid_seq RESTART WITH 1;
ALTER SEQUENCE usert_userid_seq RESTART WITH 1;