package models

import (
	"database/sql"
)

type EmergencyDevice struct {
	EmergencyDeviceID       int            `json:"emergency_device_id"`        // From emergency_deviceT table
	EmergencyDeviceTypeID   int            `json:"emergency_device_type_id"`   // From emergency_deviceT table (FK)
	EmergencyDeviceTypeName string         `json:"emergency_device_type_name"` // From emergency_device_typeT table
	ExtinguisherTypeName    sql.NullString `json:"extinguisher_type_name"`     // From Extinguisher_TypeT table
	ExtinguisherTypeID      sql.NullInt64  `json:"extinguisher_type_id"`       // From Extinguisher_TypeT table
	RoomID                  int            `json:"room_id"`                    // From emergency_deviceT table (FK)
	RoomCode                string         `json:"room_code"`                  // From roomT table
	BuildingID              int            `json:"building_id"`                // From buildingT table
	BuildingCode            string         `json:"building_code"`              // From buildingT table
	SiteID                  int            `json:"site_id"`                    // From siteT table
	SiteName                string         `json:"site_name"`                  // From siteT table
	SerialNumber            sql.NullString `json:"serial_number"`              // From emergency_deviceT table
	ManufactureDate         sql.NullTime   `json:"manufacture_date"`           // From emergency_deviceT table
	ExpireDate              sql.NullTime   `json:"expire_date"`                // Calculated
	LastInspectionDate      sql.NullTime   `json:"last_inspection_date"`       // From emergency_deviceT table
	NextInspectionDate      sql.NullTime   `json:"next_inspection_date"`       // Calculated
	Description             sql.NullString `json:"description"`                // From emergency_deviceT table
	Size                    sql.NullString `json:"size"`                       // From emergency_deviceT table
	Status                  sql.NullString `json:"status"`                     // From emergency_deviceT table
}
