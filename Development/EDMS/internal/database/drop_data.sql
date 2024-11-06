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
-- Then reset all sequences
ALTER SEQUENCE buildingt_buildingid_seq RESTART WITH 1;
ALTER SEQUENCE emergency_device_inspectiont_emergencydeviceinspectionid_seq RESTART WITH 1;
ALTER SEQUENCE emergency_device_typet_emergencydevicetypeid_seq RESTART WITH 1;
ALTER SEQUENCE emergency_devicet_emergencydeviceid_seq RESTART WITH 1;
ALTER SEQUENCE extinguisher_typet_extinguishertypeid_seq RESTART WITH 1;
ALTER SEQUENCE roomt_roomid_seq RESTART WITH 1;
ALTER SEQUENCE sitet_siteid_seq RESTART WITH 1;
ALTER SEQUENCE usert_userid_seq RESTART WITH 1;
-- Generate select script for all tables and data
