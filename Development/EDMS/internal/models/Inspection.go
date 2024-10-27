package models

import "database/sql"

// Inspection represents the inspection of a device
type Inspection struct {
	EmergencyDeviceInspectionID   int          `json:"emergency_device_inspection_id"`
	EmergencyDeviceID             int          `json:"emergency_device_id"`
	UserID                        int          `json:"user_id"`
	InspectorName                 string       `json:"inspector_name"`
	InspectionDate                sql.NullTime `json:"inspection_date"`
	CreatedAt                     sql.NullTime `json:"created_at"`
	IsConspicuous                 bool         `json:"is_conspicuous"`
	IsAccessible                  bool         `json:"is_accessible"`
	IsAssignedLocation            bool         `json:"is_assigned_location"`
	IsSignVisible                 bool         `json:"is_sign_visible"`
	IsAntiTamperDeviceIntact      bool         `json:"is_anti_tamper_device_intact"`
	IsSupportBracketSecure        bool         `json:"is_support_bracket_secure"`
	AreOperatingInstructionsClear bool         `json:"are_operating_instructions_clear"`
	IsMaintainceTagAttached       bool         `json:"is_maintaince_tag_attached"`
	IsExternalDamagePresent       bool         `json:"is_external_damage_present"`
	IsReplaced                    bool         `json:"is_replaced"`
	AreMaintenanceRecordsComplete bool         `json:"are_maintenance_records_complete"`
	WorkOrderRequired             bool         `json:"work_order_required"`
	InspectionStatus              string       `json:"inspection_status"`
	Notes                         string       `json:"notes"`
}
