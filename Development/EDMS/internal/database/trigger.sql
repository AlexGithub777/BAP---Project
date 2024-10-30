CREATE OR REPLACE FUNCTION update_device_status_on_inspection() 
RETURNS TRIGGER AS $$
DECLARE
    current_last_inspection_date DATE;
BEGIN
    -- Retrieve the current last inspection date for the device
    SELECT LastInspectionDate INTO current_last_inspection_date
    FROM Emergency_DeviceT
    WHERE EmergencyDeviceID = NEW.EmergencyDeviceID;

    -- Check if the new inspection date is more recent than the current last inspection date
    IF current_last_inspection_date IS NULL OR NEW.InspectionDate > current_last_inspection_date THEN
        -- Update the last inspection date and status in Emergency_DeviceT
        UPDATE Emergency_DeviceT
        SET 
            LastInspectionDate = NEW.InspectionDate,
            Status = CASE 
                        WHEN NEW.InspectionStatus = 'Failed' THEN 'Inspection Failed'
                        WHEN NEW.InspectionStatus = 'Passed' THEN 'Active'
                        ELSE Status
                     END
        WHERE EmergencyDeviceID = NEW.EmergencyDeviceID;

        RAISE NOTICE 'Device % status updated. Last inspection date set to %.', 
            NEW.EmergencyDeviceID, NEW.InspectionDate;
    ELSE
        RAISE NOTICE 'New inspection date % is not more recent than the current last inspection date % for device %. No update performed.', 
            NEW.InspectionDate, current_last_inspection_date, NEW.EmergencyDeviceID;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to call the function after insert on Emergency_Device_InspectionT
CREATE TRIGGER trg_update_device_status
AFTER INSERT ON Emergency_Device_InspectionT
FOR EACH ROW
EXECUTE FUNCTION update_device_status_on_inspection();
