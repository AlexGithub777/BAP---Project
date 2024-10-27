package models

import "database/sql"

// Inspection represents the inspection of a device
type Inspection struct {
	EmergencyDeviceInspectionID   int            `json:"emergency_device_inspection_id"`
	EmergencyDeviceID             int            `json:"emergency_device_id"`
	UserID                        int            `json:"user_id"`
	InspectorName                 string         `json:"inspector_name"`
	InspectionDate                sql.NullTime   `json:"inspection_date"`
	CreatedAt                     sql.NullTime   `json:"created_at"`
	IsConspicuous                 sql.NullBool   `json:"is_conspicuous"`
	IsAccessible                  sql.NullBool   `json:"is_accessible"`
	IsAssignedLocation            sql.NullBool   `json:"is_assigned_location"`
	IsSignVisible                 sql.NullBool   `json:"is_sign_visible"`
	IsAntiTamperDeviceIntact      sql.NullBool   `json:"is_anti_tamper_device_intact"`
	IsSupportBracketSecure        sql.NullBool   `json:"is_support_bracket_secure"`
	AreOperatingInstructionsClear sql.NullBool   `json:"are_operating_instructions_clear"`
	IsMaintenanceTagAttached      sql.NullBool   `json:"is_maintenance_tag_attached"`
	IsExternalDamagePresent       sql.NullBool   `json:"is_external_damage_present"`
	IsChargeGaugeNormal           sql.NullBool   `json:"is_charge_gauge_normal"`
	IsReplaced                    sql.NullBool   `json:"is_replaced"`
	AreMaintenanceRecordsComplete sql.NullBool   `json:"are_maintenance_records_complete"`
	WorkOrderRequired             sql.NullBool   `json:"work_order_required"`
	InspectionStatus              string         `json:"inspection_status"`
	Notes                         sql.NullString `json:"notes"`
}
