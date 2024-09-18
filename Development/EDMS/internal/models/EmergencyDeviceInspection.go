package models

import "time"

// Emergency_Device_InspectionT represents the inspections done on emergency devices
type EmergencyDeviceInspection struct {
	EmergencyDeviceInspectionID    int       `json:"emergency_device_inspection_id"`
	EmergencyDeviceID              int       `json:"emergency_device_id"`
	UserID                         int       `json:"user_id"`
	InspectionDate                 time.Time `json:"inspection_date"`
	Notes                          string    `json:"notes,omitempty"`
	CreatedAt                      time.Time `json:"created_at"`
	IsConspicuous                  *bool     `json:"is_conspicuous,omitempty"`
	IsAccessible                   *bool     `json:"is_accessible,omitempty"`
	IsAssignedLocation             *bool     `json:"is_assigned_location,omitempty"`
	IsSignVisible                  *bool     `json:"is_sign_visible,omitempty"`
	IsAntiTamperDeviceIntact       *bool     `json:"is_anti_tamper_device_intact,omitempty"`
	IsSupportBracketSecure         *bool     `json:"is_support_bracket_secure,omitempty"`
	AreOperatingInstructionsClear  *bool     `json:"are_operating_instructions_clear,omitempty"`
	IsMaintenanceTagAttached       *bool     `json:"is_maintenance_tag_attached,omitempty"`
	IsExternalDamagePresent        *bool     `json:"is_external_damage_present,omitempty"`
	IsChargeGaugeNormal            *bool     `json:"is_charge_gauge_normal,omitempty"`
	IsReplaced                     *bool     `json:"is_replaced,omitempty"`
	AreMaintenanceRecordsComplete  *bool     `json:"are_maintenance_records_complete,omitempty"`
	WorkOrderRequired              *bool     `json:"work_order_required,omitempty"`
	IsInspectionComplete           *bool     `json:"is_inspection_complete,omitempty"`
}