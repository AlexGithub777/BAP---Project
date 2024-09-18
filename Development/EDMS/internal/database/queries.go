package database

import (
	"database/sql"
	"time"

	"github.com/AlexGithub777/safety-device-app/internal/models"
)

// GetAllUsers function
func (db *DB) GetAllUsers() ([]models.User, error) {
	query := `SELECT userid, username, email, role FROM userT`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.Email,
			&user.Role,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Create user function
func (db *DB) CreateUser(user *models.User) error {
	query := `
		INSERT INTO userT (username, password, email)
		VALUES ($1, $2, $3)
		`
	insertStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, err = insertStmt.Exec(user.Username, user.Password, user.Email)

	if err != nil {
		return err
	}

	return nil
}

// Get user by username function
func (db *DB) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT userid, username, password, email, role
		FROM userT
		WHERE username = $1
		`
	var user models.User
	err := db.QueryRow(query, username).Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update password function
func (db *DB) UpdatePassword(userid int, password string) error {
	query := `
		UPDATE userT
		SET password = $1
		WHERE userid = $2
		`
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(password, userid)

	if err != nil {
		return err
	}

	return nil
}

// Get user by email function
func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT userid, username, password, email, role
		FROM userT
		WHERE email = $1
		`
	var user models.User
	err := db.QueryRow(query, email).Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Refactor/add new function to GetAllDevices by Site
func (db *DB) GetAllDevices(siteId string, buildingCode string) ([]models.EmergencyDevice, error) {
	var query string
	var args []interface{}

	// Define the base query
	query = `
	SELECT 
		ed.emergencydeviceid, 
		edt.emergencydevicetypename,
		et.extinguishertypename AS ExtinguisherTypeName,
		r.roomcode,
		ed.serialnumber,
		ed.manufacturedate,
		ed.lastinspectiondate,
		ed.description,
		ed.size,
		ed.status 
	FROM emergency_deviceT ed
	JOIN roomT r ON ed.roomid = r.roomid
	LEFT JOIN emergency_device_typeT edt ON ed.emergencydevicetypeid = edt.emergencydevicetypeid
	LEFT JOIN Extinguisher_TypeT et ON ed.extinguishertypeid = et.extinguishertypeid
	`

	// Add filtering by site name if provided
	if siteId != "" {
		query += `
		JOIN buildingT b ON r.buildingid = b.buildingid
		JOIN siteT s ON b.siteid = s.siteid
		WHERE s.siteid = $1
		`
		args = append(args, siteId)
	} else if buildingCode != "" {
		// Add filtering by building code if provided
		query += `
		JOIN buildingT b ON r.buildingid = b.buildingid
		WHERE b.buildingcode = $1
		`
		args = append(args, buildingCode)
	}

	// Prepare and execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define the result slice
	var emergencyDevices []models.EmergencyDevice

	// Scan the results
	for rows.Next() {
		var device models.EmergencyDevice
		err := rows.Scan(
			&device.EmergencyDeviceID,
			&device.EmergencyDeviceTypeName,
			&device.ExtinguisherTypeName,
			&device.RoomCode,
			&device.SerialNumber,
			&device.ManufactureDate,
			&device.LastInspectionDate,
			&device.Description,
			&device.Size,
			&device.Status,
		)
		if err != nil {
			return nil, err
		}

		// If any of the following fields are null, replace them with a default value
		if !device.ExtinguisherTypeName.Valid {
			device.ExtinguisherTypeName.String = "N/A"
			device.ExtinguisherTypeName.Valid = false
		}
		if !device.SerialNumber.Valid {
			device.SerialNumber.String = "N/A"
			device.SerialNumber.Valid = false
		}
		if !device.Description.Valid {
			device.Description.String = "N/A"
			device.Description.Valid = false
		}
		if !device.Size.Valid {
			device.Size.String = "N/A"
			device.Size.Valid = false
		}
		if !device.Status.Valid {
			device.Status.String = "N/A"
			device.Status.Valid = false
		}

		// Handle dates and calculate expiry and next inspection dates
		if device.ManufactureDate.Valid {
			expiryDate := device.ManufactureDate.Time.AddDate(5, 0, 0)
			device.ExpireDate = sql.NullTime{
				Time:  expiryDate,
				Valid: true,
			}
		} else {
			device.ManufactureDate = sql.NullTime{
				Time:  time.Time{}, // Zero value of time.Time
				Valid: false,
			}
			device.ExpireDate = sql.NullTime{
				Time:  time.Time{}, // Zero value of time.Time
				Valid: false,
			}
		}

		if device.LastInspectionDate.Valid {
			nextInspectionDate := device.LastInspectionDate.Time.AddDate(0, 3, 0)
			device.NextInspectionDate = sql.NullTime{
				Time:  nextInspectionDate,
				Valid: true,
			}
		} else {
			device.NextInspectionDate = sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			}
		}

		emergencyDevices = append(emergencyDevices, device)
	}

	return emergencyDevices, nil
}

func (db *DB) GetAllDeviceTypes() ([]models.EmergencyDeviceType, error) {
	query := `
	SELECT emergencydevicetypeid, emergencydevicetypename
	FROM emergency_device_typeT
	ORDER BY emergencydevicetypename
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deviceTypes []models.EmergencyDeviceType

	// Scan the results
	for rows.Next() {
		var deviceType models.EmergencyDeviceType
		err := rows.Scan(
			&deviceType.EmergencyDeviceTypeID,
			&deviceType.EmergencyDeviceTypeName,
		)
		if err != nil {
			return nil, err
		}

		deviceTypes = append(deviceTypes, deviceType)
	}

	return deviceTypes, nil
}

func (db *DB) GetAllExtinguisherTypes() ([]models.ExtinguisherType, error) {
	query := `
	SELECT extinguishertypeid, extinguishertypename
	FROM Extinguisher_TypeT
	ORDER BY extinguishertypename
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var extinguisherTypes []models.ExtinguisherType

	// Scan the results
	for rows.Next() {
		var extinguisherType models.ExtinguisherType
		err := rows.Scan(
			&extinguisherType.ExtinguisherTypeID,
			&extinguisherType.ExtinguisherTypeName,
		)
		if err != nil {
			return nil, err
		}

		extinguisherTypes = append(extinguisherTypes, extinguisherType)
	}

	return extinguisherTypes, nil
}

func (db *DB) GetAllRooms(buildingId string) ([]models.Room, error) {
	var query string
	var args []interface{}

	// Define the base query
	query = ` SELECT r.roomid, r.buildingid, r.roomcode, b.buildingcode, s.sitename
              FROM roomT r
              JOIN buildingT b ON r.buildingid = b.buildingid
              JOIN siteT s ON b.siteid = s.siteid`

	// Add filtering by building code if provided
	if buildingId != "" {
		query += ` WHERE b.buildingId = $1`
		args = append(args, buildingId)
	}

	// Prepare and execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define the result slice
	var rooms []models.Room

	// Scan the results
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.RoomID,
			&room.BuildingID,
			&room.RoomCode,
			&room.BuildingCode, // Assuming you have added this field to the Room model
			&room.SiteName,     // Assuming you have added this field to the Room model
		)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (db *DB) GetAllBuildings(siteId string) ([]models.Building, error) {
	var args []interface{}
	query := `
    SELECT b.buildingid, b.buildingcode, b.siteid, s.sitename
    FROM buildingT b
    JOIN siteT s ON b.siteid = s.siteid
    `

	// Add filtering by site name if provided
	if siteId != "" {
		query += ` WHERE s.siteid = $1`
		args = append(args, siteId)
	}

	query += ` ORDER BY b.buildingcode`

	// Prepare and execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Define the result slice

	var buildings []models.Building

	// Scan the results
	for rows.Next() {
		var building models.Building
		err := rows.Scan(
			&building.BuildingID,
			&building.BuildingCode,
			&building.SiteID,
			&building.SiteName, // Assuming you have added this field to the Building model
		)
		if err != nil {
			return nil, err
		}

		buildings = append(buildings, building)
	}

	return buildings, nil
}

func (db *DB) GetAllSites() ([]models.Site, error) {
	query := `
	SELECT siteid, sitename, siteaddress
	FROM siteT
	ORDER BY sitename
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sites []models.Site

	// Scan the results
	for rows.Next() {
		var site models.Site
		err := rows.Scan(
			&site.SiteID,
			&site.SiteName,
			&site.SiteAddress,
		)
		if err != nil {
			return nil, err
		}

		sites = append(sites, site)
	}

	return sites, nil

}
