CREATE OR REPLACE FUNCTION update_device_status_on_inspection() 
RETURNS TRIGGER AS $$
DECLARE
    current_last_inspection_timestamp TIMESTAMP;
BEGIN
    -- Retrieve the current last inspection timestamp for the device
    SELECT LastInspectionDateTime INTO current_last_inspection_timestamp
    FROM Emergency_DeviceT
    WHERE EmergencyDeviceID = NEW.EmergencyDeviceID;

    -- Check if the new inspection timestamp is more recent than the current last inspection timestamp
    IF current_last_inspection_timestamp IS NULL OR NEW.InspectionDateTime > current_last_inspection_timestamp THEN
        -- Update LastInspectionDate with the new inspection timestamp
        UPDATE Emergency_DeviceT
        SET LastInspectionDateTime = NEW.InspectionDateTime
        Status = CASE 
                        WHEN NEW.InspectionStatus = 'Failed' THEN 'Inspection Failed'
                        WHEN NEW.InspectionStatus = 'Passed' THEN 'Active'
                        ELSE Status
                     END
        WHERE EmergencyDeviceID = NEW.EmergencyDeviceID;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to call the function after insert on Emergency_Device_InspectionT
CREATE TRIGGER trg_update_device_status
AFTER INSERT ON Emergency_Device_InspectionT
FOR EACH ROW
EXECUTE FUNCTION update_device_status_on_inspection();
