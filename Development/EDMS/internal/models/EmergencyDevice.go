package models

import (
	"database/sql"
)

type EmergencyDevice struct {
	EmergencyDeviceID       int            `json:"emergency_device_id"`        // From emergency_deviceT table
	EmergencyDeviceTypeName string         `json:"emergency_device_type_name"` // From emergency_device_typeT table
	ExtinguisherTypeName    sql.NullString `json:"extinguisher_type_name"`     // From Extinguisher_TypeT table
	RoomCode                string         `json:"room_code"`                  // From roomT table
	SerialNumber            sql.NullString `json:"serial_number"`              // From emergency_deviceT table
	ManufactureDate         sql.NullTime   `json:"manufacture_date"`           // From emergency_deviceT table
	ExpireDate              sql.NullTime   `json:"expire_date"`                // Calculated
	LastInspectionDate      sql.NullTime   `json:"last_inspection_date"`       // From emergency_deviceT table
	NextInspectionDate      sql.NullTime   `json:"next_inspection_date"`       // Calculated
	Description             sql.NullString `json:"description"`                // From emergency_deviceT table
	Size                    sql.NullString `json:"size"`                       // From emergency_deviceT table
	Status                  sql.NullString `json:"status"`                     // From emergency_deviceT table
}
