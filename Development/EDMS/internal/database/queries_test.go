package database_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/AlexGithub777/safety-device-app/internal/database"
	"github.com/AlexGithub777/safety-device-app/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllDevices(t *testing.T) {
	testCases := []struct {
		name            string
		expectedDevices []models.EmergencyDevice
		mockSetup       func(mock sqlmock.Sqlmock)
		expectedError   error
	}{
		{
			name: "TestFetchAllDevices with valid non-null fields",
			expectedDevices: []models.EmergencyDevice{
				{
					EmergencyDeviceID:       1,
					EmergencyDeviceTypeName: "TypeA",
					ExtinguisherTypeName:    sql.NullString{String: "ExtinguisherA", Valid: true},
					RoomCode:                "Room101",
					SerialNumber:            sql.NullString{String: "SN123", Valid: true},
					ManufactureDate:         sql.NullTime{Time: time.Date(2024, time.August, 1, 0, 0, 0, 0, time.UTC), Valid: true},
					LastInspectionDate:      sql.NullTime{Time: time.Date(2024, time.July, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					Description:             sql.NullString{String: "Description of device", Valid: true},
					Size:                    sql.NullString{String: "10kg", Valid: true},
					Status:                  sql.NullString{String: "Active", Valid: true},
					NextInspectionDate:      sql.NullTime{Time: time.Date(2024, time.October, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					ExpireDate:              sql.NullTime{Time: time.Date(2029, time.August, 1, 0, 0, 0, 0, time.UTC), Valid: true},
				},
			},
		},
		{
			name: "TestFetchAllDevices with valid null fields",
			expectedDevices: []models.EmergencyDevice{
				{
					EmergencyDeviceID:       2,
					EmergencyDeviceTypeName: "TypeB",
					ExtinguisherTypeName:    sql.NullString{Valid: false, String: "N/A"},
					RoomCode:                "Room102",
					SerialNumber:            sql.NullString{Valid: false, String: "N/A"},
					ManufactureDate:         sql.NullTime{Valid: false},
					LastInspectionDate:      sql.NullTime{Valid: false},
					Description:             sql.NullString{Valid: false, String: "N/A"},
					Size:                    sql.NullString{Valid: false, String: "N/A"},
					Status:                  sql.NullString{Valid: false, String: "N/A"},
					NextInspectionDate:      sql.NullTime{Valid: false},
					ExpireDate:              sql.NullTime{Valid: false},
				},
			},
		},
		{
			name: "TestFetchAllDevices with invalid manufacture date",
			expectedDevices: []models.EmergencyDevice{
				{
					EmergencyDeviceID:       3,
					EmergencyDeviceTypeName: "TypeC",
					ExtinguisherTypeName:    sql.NullString{String: "ExtinguisherD", Valid: true},
					RoomCode:                "C103",
					SerialNumber:            sql.NullString{String: "SN456", Valid: true},
					ManufactureDate:         sql.NullTime{Valid: false},
					LastInspectionDate:      sql.NullTime{Time: time.Date(2024, time.July, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					Description:             sql.NullString{String: "Device with invalid manufacture date", Valid: true},
					Size:                    sql.NullString{String: "5kg", Valid: true},
					Status:                  sql.NullString{String: "Active", Valid: true},
					NextInspectionDate:      sql.NullTime{Time: time.Date(2024, time.October, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					ExpireDate:              sql.NullTime{Valid: false},
				},
			},
		},
		{
			name: "TestFetchAllDevices with invalid inspection date",
			expectedDevices: []models.EmergencyDevice{
				{
					EmergencyDeviceID:       4,
					EmergencyDeviceTypeName: "TypeD",
					ExtinguisherTypeName:    sql.NullString{String: "ExtinguisherE", Valid: true},
					RoomCode:                "D104",
					SerialNumber:            sql.NullString{String: "SN789", Valid: true},
					ManufactureDate:         sql.NullTime{Time: time.Date(2024, time.August, 1, 0, 0, 0, 0, time.UTC), Valid: true},
					LastInspectionDate:      sql.NullTime{Valid: false},
					Description:             sql.NullString{String: "Device with invalid inspection date", Valid: true},
					Size:                    sql.NullString{String: "8kg", Valid: true},
					Status:                  sql.NullString{String: "Active", Valid: true},
					NextInspectionDate:      sql.NullTime{Valid: false},
					ExpireDate:              sql.NullTime{Time: time.Date(2029, time.August, 1, 0, 0, 0, 0, time.UTC), Valid: true},
				},
			},
		},

		{
			name: "TestFetchAllDevices with query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM emergency_deviceT").WillReturnError(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name: "TestFetchAllDevices with row scan error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"emergencydeviceid",
					"emergencydevicetypename",
					"extinguishertypename",
					"roomname",
					"serialnumber",
					"manufacturedate",
					"lastinspectiondate",
					"description",
					"size",
					"status",
				}).AddRow(
					"invalid", // This will cause a scan error as it's not an int
					"TypeA",
					sql.NullString{String: "ExtinguisherA", Valid: true},
					"Room101",
					sql.NullString{String: "SN123", Valid: true},
					sql.NullTime{Time: time.Now(), Valid: true},
					sql.NullTime{Time: time.Now(), Valid: true},
					sql.NullString{String: "Description", Valid: true},
					sql.NullString{String: "10kg", Valid: true},
					sql.NullString{String: "Active", Valid: true},
				)
				mock.ExpectQuery("^SELECT (.+) FROM emergency_deviceT").WillReturnRows(rows)
			},
			expectedError: errors.New("sql: Scan error on column index 0, name \"emergencydeviceid\": converting driver.Value type string (\"invalid\") to a int: invalid syntax"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create SQL mock: %v", err)
			}
			defer db.Close()

			dbInstance := &database.DB{DB: db}

			if tc.mockSetup != nil {
				tc.mockSetup(mock)
			} else {
				rows := sqlmock.NewRows([]string{
					"emergencydeviceid",
					"emergencydevicetypename",
					"extinguishertypename",
					"roomname",
					"serialnumber",
					"manufacturedate",
					"lastinspectiondate",
					"description",
					"size",
					"status",
				})

				for _, device := range tc.expectedDevices {
					rows.AddRow(
						device.EmergencyDeviceID,
						device.EmergencyDeviceTypeName,
						device.ExtinguisherTypeName,
						device.RoomCode,
						device.SerialNumber,
						device.ManufactureDate,
						device.LastInspectionDate,
						device.Description,
						device.Size,
						device.Status,
					)
				}

				mock.ExpectQuery("^SELECT (.+) FROM emergency_deviceT").WillReturnRows(rows)
			}

			actualDevices, err := dbInstance.GetAllDevices("your_building_code", "your_site_name")

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
				assert.Nil(t, actualDevices)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.expectedDevices), len(actualDevices), "Number of devices does not match")

				for i, expectedDevice := range tc.expectedDevices {
					actualDevice := actualDevices[i]
					assert.Equal(t, expectedDevice.EmergencyDeviceID, actualDevice.EmergencyDeviceID, "EmergencyDeviceID does not match")
					assert.Equal(t, expectedDevice.EmergencyDeviceTypeName, actualDevice.EmergencyDeviceTypeName, "EmergencyDeviceTypeName does not match")
					assert.Equal(t, expectedDevice.ExtinguisherTypeName, actualDevice.ExtinguisherTypeName, "ExtinguisherTypeName does not match")
					assert.Equal(t, expectedDevice.RoomCode, actualDevice.RoomCode, "RoomCode does not match")
					assert.Equal(t, expectedDevice.SerialNumber, actualDevice.SerialNumber, "SerialNumber does not match")
					assert.Equal(t, expectedDevice.ManufactureDate, actualDevice.ManufactureDate, "ManufactureDate does not match")
					assert.Equal(t, expectedDevice.LastInspectionDate, actualDevice.LastInspectionDate, "LastInspectionDate does not match")
					assert.Equal(t, expectedDevice.Description, actualDevice.Description, "Description does not match")
					assert.Equal(t, expectedDevice.Size, actualDevice.Size, "Size does not match")
					assert.Equal(t, expectedDevice.Status, actualDevice.Status, "Status does not match")
					assert.Equal(t, expectedDevice.NextInspectionDate, actualDevice.NextInspectionDate, "NextInspectionDate does not match")
					assert.Equal(t, expectedDevice.ExpireDate, actualDevice.ExpireDate, "ExpireDate does not match")
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
