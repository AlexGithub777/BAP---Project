package database

import (
	"database/sql"
	"time"

	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/models"
)

// GetAllUsers function
func (db *DB) GetAllUsers() ([]models.User, error) {
	query := `SELECT userid, username, email, role, defaultadmin FROM userT`
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
			&user.DefaultAdmin,
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

// Update user function
func (db *DB) UpdateUserWithPassword(user *models.User) error {
	query := `
        UPDATE userT
        SET username = $1, email = $2, role = $3, password = $4
        WHERE userid = $5
        `
	args := []interface{}{user.Username, user.Email, user.Role, user.Password, user.UserID}

	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}

// Update user function
func (db *DB) UpdateUser(user *models.User) error {
	query := `
		UPDATE userT
		SET username = $1, email = $2, role = $3
		WHERE userid = $4
		`

	args := []interface{}{user.Username, user.Email, user.Role, user.UserID}

	updateStmt, err := db.Prepare(query)

	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(args...)

	if err != nil {
		return err
	}

	return nil

}

// Get user by username function
func (db *DB) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT userid, username, password, email, role, defaultadmin
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
		&user.DefaultAdmin,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get user by ID function
func (db *DB) GetUserByID(userid int) (*models.User, error) {
	query := `
		SELECT userid, username, password, email, role, defaultadmin
		FROM userT
		WHERE userid = $1
		`

	var user models.User
	err := db.QueryRow(query, userid).Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
		&user.DefaultAdmin,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Delete user function
func (db *DB) DeleteUser(userid int) error {
	query := `DELETE FROM userT WHERE userid = $1`
	deleteStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(userid)

	if err != nil {
		return err
	}

	return nil
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
		ed.LastInspectionDateTime AT TIME ZONE 'Pacific/Auckland' AS lastinspectiondatetime_nzdt,
		ed.description,
		ed.size,
		ed.status 
	FROM emergency_deviceT ed
	JOIN roomT r ON ed.roomid = r.roomid
	LEFT JOIN emergency_device_typeT edt ON ed.emergencydevicetypeid = edt.emergencydevicetypeid
	LEFT JOIN Extinguisher_TypeT et ON ed.extinguishertypeid = et.extinguishertypeid
	`

	// Add filtering by site name and building code if provided
	if siteId != "" && buildingCode != "" {
		query += `
		JOIN buildingT b ON r.buildingid = b.buildingid
		JOIN siteT s ON b.siteid = s.siteid
		WHERE s.siteid = $1 AND b.buildingcode = $2
		`
		args = append(args, siteId, buildingCode)
	} else if siteId != "" {
		// Add filtering by site name if provided
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
			&device.LastInspectionDateTime,
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

		if device.LastInspectionDateTime.Valid {
			nextInspectionDate := device.LastInspectionDateTime.Time.AddDate(0, 3, 0)
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

// GetDeviceByID function
func (db *DB) GetDeviceByID(deviceID int) (*models.EmergencyDevice, error) {
	query := `
	SELECT
		ed.emergencydeviceid,
		ed.emergencydevicetypeid,
		edt.emergencydevicetypename,
		et.extinguishertypename,
		ed.extinguishertypeid,
		ed.roomid,
		r.roomcode,
		b.buildingid,
		b.buildingcode,
		s.siteid,
		s.sitename,
		ed.serialnumber,
		ed.manufacturedate,
		ed.LastInspectionDateTime AT TIME ZONE 'Pacific/Auckland' AS lastinspectiondatetime_nzdt,
		ed.description,
		ed.size,
		ed.status
	FROM emergency_deviceT ed
	JOIN emergency_device_typeT edt ON ed.emergencydevicetypeid = edt.emergencydevicetypeid
	LEFT JOIN Extinguisher_TypeT et ON ed.extinguishertypeid = et.extinguishertypeid
	JOIN roomT r ON ed.roomid = r.roomid
	JOIN buildingT b ON r.buildingid = b.buildingid
	JOIN siteT s ON b.siteid = s.siteid
	WHERE ed.emergencydeviceid = $1
	`
	var device models.EmergencyDevice
	err := db.QueryRow(query, deviceID).Scan(
		&device.EmergencyDeviceID,
		&device.EmergencyDeviceTypeID,
		&device.EmergencyDeviceTypeName,
		&device.ExtinguisherTypeName, // Emergency device type name (string)
		&device.ExtinguisherTypeID,   // Extinguisher type ID (int, can be NULL)
		&device.RoomID,
		&device.RoomCode,
		&device.BuildingID,
		&device.BuildingCode,
		&device.SiteID,
		&device.SiteName,
		&device.SerialNumber,
		&device.ManufactureDate,
		&device.LastInspectionDateTime,
		&device.Description,
		&device.Size,
		&device.Status,
	)

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (db *DB) GetDevicesByRoomID(roomID int) ([]models.EmergencyDevice, error) {
	query := `
	SELECT
		ed.emergencydeviceid,
		ed.emergencydevicetypeid,
		edt.emergencydevicetypename,
		et.extinguishertypename,
		ed.extinguishertypeid,
		ed.roomid,
		r.roomcode,
		b.buildingid,
		b.buildingcode,
		s.siteid,
		s.sitename,
		ed.serialnumber,
		ed.manufacturedate,
		ed.LastInspectionDateTime AT TIME ZONE 'Pacific/Auckland' AS lastinspectiondatetime_nzdt,
		ed.description,
		ed.size,
		ed.status
	FROM emergency_deviceT ed
	JOIN emergency_device_typeT edt ON ed.emergencydevicetypeid = edt.emergencydevicetypeid
	LEFT JOIN Extinguisher_TypeT et ON ed.extinguishertypeid = et.extinguishertypeid
	JOIN roomT r ON ed.roomid = r.roomid
	JOIN buildingT b ON r.buildingid = b.buildingid
	JOIN siteT s ON b.siteid = s.siteid
	WHERE ed.roomid = $1
	`
	rows, err := db.Query(query, roomID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var emergencyDevices []models.EmergencyDevice

	// Scan the results
	for rows.Next() {
		var device models.EmergencyDevice
		err := rows.Scan(
			&device.EmergencyDeviceID,
			&device.EmergencyDeviceTypeID,
			&device.EmergencyDeviceTypeName,
			&device.ExtinguisherTypeName,
			&device.ExtinguisherTypeID,
			&device.RoomID,
			&device.RoomCode,
			&device.BuildingID,
			&device.BuildingCode,
			&device.SiteID,
			&device.SiteName,
			&device.SerialNumber,
			&device.ManufactureDate,
			&device.LastInspectionDateTime,
			&device.Description,
			&device.Size,
			&device.Status,
		)
		if err != nil {
			return nil, err
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

func (db *DB) GetBuildingById(buildingID string) (*models.Building, error) {
	query := `
	SELECT buildingid, siteid, buildingcode
	FROM buildingT
	WHERE buildingid = $1
	`
	var building models.Building
	err := db.QueryRow(query, buildingID).Scan(
		&building.BuildingID,
		&building.SiteID,
		&building.BuildingCode,
	)

	if err != nil {
		return nil, err
	}

	return &building, nil
}

func (db *DB) AddBuilding(building *models.Building) error {
	query := "INSERT INTO buildingT (siteId, buildingCode) VALUES ($1, $2)"
	insertStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, err = insertStmt.Exec(building.SiteID, building.BuildingCode)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateBuilding(building *models.Building) error {
	query := "UPDATE BuildingT SET siteId = $1, buildingCode = $2 WHERE buildingID = $3"
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(building.SiteID, building.BuildingCode, building.BuildingID)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteBuilding(buildingID string) error {
	query := "DELETE FROM BuildingT WHERE buildingID = $1"
	deleteStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(buildingID)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetRoomsByBuildingID(buildingID string) ([]models.Room, error) {
	query := `
	SELECT r.roomid
	FROM roomT r
	WHERE r.buildingid = $1
	ORDER BY r.roomcode
	`

	rows, err := db.Query(query, buildingID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rooms []models.Room

	// Scan the results
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.RoomID,
		)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
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

func (db *DB) GetSiteByID(siteID string) (*models.Site, error) {
	query := `
	SELECT siteid, sitename, siteaddress, sitemapimagepath
	FROM siteT
	WHERE siteid = $1
	`
	var site models.Site
	err := db.QueryRow(query, siteID).Scan(
		&site.SiteID,
		&site.SiteName,
		&site.SiteAddress,
		&site.SiteMapImagePath,
	)

	if err != nil {
		return nil, err
	}

	return &site, nil
}

// Get site by name function
func (db *DB) GetSiteByName(siteName string) (*models.Site, error) {
	query := `
	SELECT siteid, sitename, siteaddress, sitemapimagepath
	FROM siteT
	WHERE sitename = $1
	`

	var site models.Site
	err := db.QueryRow(query, siteName).Scan(
		&site.SiteID,
		&site.SiteName,
		&site.SiteAddress,
		&site.SiteMapImagePath,
	)

	if err != nil {
		return nil, err
	}

	return &site, nil
}

func (db *DB) AddSite(site *models.Site) error {
	query := "INSERT INTO SiteT (siteName, siteAddress, siteMapImagePath) VALUES ($1, $2, $3)"
	insertStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, err = insertStmt.Exec(site.SiteName, site.SiteAddress, site.SiteMapImagePath)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateSite(site *models.Site) error {
	query := "UPDATE SiteT SET siteName = $1, siteAddress = $2, siteMapImagePath = $3 WHERE siteID = $4"
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(site.SiteName, site.SiteAddress, site.SiteMapImagePath, site.SiteID)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteSite(siteID string) error {
	query := "DELETE FROM SiteT WHERE siteID = $1"
	deleteStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(siteID)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetRoomsBySiteID(siteID string) ([]models.Room, error) {
	query := `
	SELECT r.roomid, r.roomcode, b.buildingcode, s.sitename
	FROM roomT r
	JOIN buildingT b ON r.buildingid = b.buildingid
	JOIN siteT s ON b.siteid = s.siteid
	WHERE s.siteid = $1
	ORDER BY r.roomcode
	`

	rows, err := db.Query(query, siteID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rooms []models.Room

	// Scan the results
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.RoomID,
			&room.RoomCode,
			&room.BuildingCode,
			&room.SiteName,
		)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (db *DB) GetRoomByID(roomID int) (*models.Room, error) {
	query := `
	SELECT r.roomid, r.roomcode, b.buildingid,b.buildingcode, s.sitename, s.siteid
	FROM roomT r
	JOIN buildingT b ON r.buildingid = b.buildingid
	JOIN siteT s ON b.siteid = s.siteid
	WHERE r.roomid = $1
	`

	var room models.Room
	err := db.QueryRow(query, roomID).Scan(
		&room.RoomID,
		&room.RoomCode,
		&room.BuildingID,
		&room.BuildingCode,
		&room.SiteName,
		&room.SiteID,
	)

	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (db *DB) AddRoom(room *models.Room) error {
	query := "INSERT INTO RoomT (buildingId, roomCode) VALUES ($1, $2)"
	insertStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, err = insertStmt.Exec(room.BuildingID, room.RoomCode)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateRoom(room *models.Room) error {
	query := "UPDATE RoomT SET buildingId = $1, roomCode = $2 WHERE roomID = $3"
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(room.BuildingID, room.RoomCode, room.RoomID)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteRoom(roomID int) error {
	query := "DELETE FROM RoomT WHERE roomID = $1"
	deleteStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(roomID)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetEmergencyDeviceTypeByID(emergencyDeviceTypeID int) (*models.EmergencyDeviceType, error) {
	query := `
	SELECT emergencydevicetypeid, emergencydevicetypename
	FROM emergency_device_typeT
	WHERE emergencydevicetypeid = $1
	`

	var deviceType models.EmergencyDeviceType
	err := db.QueryRow(query, emergencyDeviceTypeID).Scan(
		&deviceType.EmergencyDeviceTypeID,
		&deviceType.EmergencyDeviceTypeName,
	)

	if err != nil {
		return nil, err
	}

	return &deviceType, nil
}

func (db *DB) GetDeviceTypeByName(emergencyDeviceTypeName string) (*models.EmergencyDeviceType, error) {
	query := `
	SELECT emergencydevicetypeid, emergencydevicetypename
	FROM emergency_device_typeT
	WHERE emergencydevicetypename = $1
	`

	var deviceType models.EmergencyDeviceType
	err := db.QueryRow(query, emergencyDeviceTypeName).Scan(
		&deviceType.EmergencyDeviceTypeID,
		&deviceType.EmergencyDeviceTypeName,
	)

	if err != nil {
		return nil, err
	}

	return &deviceType, nil
}

func (db *DB) AddEmergencyDeviceType(emergencyDeviceTypeName string) error {
	query := `
	INSERT INTO emergency_device_typeT (emergencydevicetypename)
	VALUES ($1)
	`
	insertStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, err = insertStmt.Exec(
		emergencyDeviceTypeName,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateEmergencyDeviceType(emergencyDeviceType *models.EmergencyDeviceType) error {
	query := `
	UPDATE emergency_device_typeT
	SET emergencydevicetypename = $1
	WHERE emergencydevicetypeid = $2
	`

	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(
		emergencyDeviceType.EmergencyDeviceTypeName,
		emergencyDeviceType.EmergencyDeviceTypeID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteEmergencyDeviceType(emergencyDeviceTypeID int) error {
	query := "DELETE FROM Emergency_Device_TypeT WHERE EmergencyDeviceTypeID = $1"
	deleteStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(emergencyDeviceTypeID)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetDevicesByTypeID(emergencyDeviceTypeID int) ([]models.EmergencyDevice, error) {
	query := `
    SELECT
        ed.emergencydeviceid,
        ed.emergencydevicetypeid,
        edt.emergencydevicetypename,
        et.extinguishertypename,
        ed.extinguishertypeid,
        ed.roomid,
        r.roomcode,
        b.buildingid,
        b.buildingcode,
        s.siteid,
        s.sitename,
        ed.serialnumber,
        ed.manufacturedate,
        ed.LastInspectionDateTime AT TIME ZONE 'Pacific/Auckland' AS lastinspectiondatetime_nzdt,
        ed.description,
        ed.size,
        ed.status
    FROM emergency_deviceT ed
    JOIN emergency_device_typeT edt ON ed.emergencydevicetypeid = edt.emergencydevicetypeid
    LEFT JOIN Extinguisher_TypeT et ON ed.extinguishertypeid = et.extinguishertypeid
    JOIN roomT r ON ed.roomid = r.roomid
    JOIN buildingT b ON r.buildingid = b.buildingid
    JOIN siteT s ON b.siteid = s.siteid
    WHERE ed.emergencydevicetypeid = $1
    `
	rows, err := db.Query(query, emergencyDeviceTypeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []models.EmergencyDevice
	for rows.Next() {
		var device models.EmergencyDevice
		err := rows.Scan(
			&device.EmergencyDeviceID,
			&device.EmergencyDeviceTypeID,
			&device.EmergencyDeviceTypeName,
			&device.ExtinguisherTypeName,
			&device.ExtinguisherTypeID,
			&device.RoomID,
			&device.RoomCode,
			&device.BuildingID,
			&device.BuildingCode,
			&device.SiteID,
			&device.SiteName,
			&device.SerialNumber,
			&device.ManufactureDate,
			&device.LastInspectionDateTime,
			&device.Description,
			&device.Size,
			&device.Status,
		)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func (db *DB) GetExtinguisherTypeByID(extinguisherTypeID int) (*models.ExtinguisherType, error) {
	query := `
	SELECT extinguishertypeid, extinguishertypename
	FROM extinguisher_typeT
	WHERE extinguishertypeid = $1
	`

	var extinguisherType models.ExtinguisherType
	err := db.QueryRow(query, extinguisherTypeID).Scan(
		&extinguisherType.ExtinguisherTypeID,
		&extinguisherType.ExtinguisherTypeName,
	)

	if err != nil {
		return nil, err
	}

	return &extinguisherType, nil
}

func (db *DB) AddEmergencyDevice(device *models.EmergencyDevice) error {
	query := `
	INSERT INTO emergency_deviceT (emergencydevicetypeid, extinguishertypeid, roomid, serialnumber, manufacturedate, description, size, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	insertStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, err = insertStmt.Exec(
		device.EmergencyDeviceTypeID,
		device.ExtinguisherTypeID,
		device.RoomID,
		device.SerialNumber,
		device.ManufactureDate,
		device.Description,
		device.Size,
		device.Status,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateEmergencyDevice(device *models.EmergencyDevice) error {
	query := `
	UPDATE emergency_deviceT
	SET emergencydevicetypeid = $1, extinguishertypeid = $2, roomid = $3, serialnumber = $4, manufacturedate = $5, description = $6, size = $7, status = $8
	WHERE emergencydeviceid = $9
	`
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(
		device.EmergencyDeviceTypeID,
		device.ExtinguisherTypeID,
		device.RoomID,
		device.SerialNumber,
		device.ManufactureDate,
		device.Description,
		device.Size,
		device.Status,
		device.EmergencyDeviceID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteEmergencyDevice(deviceID int) error {
	query := "DELETE FROM Emergency_DeviceT WHERE EmergencyDeviceID = $1"
	deleteStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(deviceID)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetAllInspectionsByDeviceID(deviceID int) ([]models.Inspection, error) {
	query := `
	SELECT edi.emergencydeviceinspectionid, edi.emergencydeviceid, ed.serialnumber, edi.userid, u.username, edi.inspectiondatetime AT TIME ZONE 'Pacific/Auckland' AS lastinspectiondate_nzdt, edi.createdat AT TIME ZONE 'Pacific/Auckland' AS createdat_nzdt,
		   edi.IsConspicuous, edi.IsAccessible, edi.IsAssignedLocation, edi.IsSignVisible, edi.IsAntiTamperDeviceIntact,
		   edi.IsSupportBracketSecure, edi.AreOperatingInstructionsClear, edi.IsMaintenanceTagAttached,
		   edi.isNoExternalDamage, edi.IsChargeGaugeNormal, edi.IsReplaced, edi.AreMaintenanceRecordsComplete, edi.WorkOrderRequired,
		   edi.InspectionStatus, edi.Notes
	FROM emergency_device_inspectionT edi
	JOIN userT u ON edi.userid = u.userid
	JOIN emergency_deviceT ed ON edi.emergencydeviceid = ed.emergencydeviceid
	WHERE edi.emergencydeviceid = $1
	ORDER BY edi.inspectiondatetime DESC
	`

	rows, err := db.Query(query, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inspections []models.Inspection

	// Scan the results
	for rows.Next() {
		var inspection models.Inspection
		err := rows.Scan(
			&inspection.EmergencyDeviceInspectionID,
			&inspection.EmergencyDeviceID,
			&inspection.SerialNumber,
			&inspection.UserID,
			&inspection.InspectorName, // Scans the `username` field into `InspectorName`
			&inspection.InspectionDateTime,
			&inspection.CreatedAt,
			&inspection.IsConspicuous,
			&inspection.IsAccessible,
			&inspection.IsAssignedLocation,
			&inspection.IsSignVisible,
			&inspection.IsAntiTamperDeviceIntact,
			&inspection.IsSupportBracketSecure,
			&inspection.AreOperatingInstructionsClear,
			&inspection.IsMaintenanceTagAttached,
			&inspection.IsNoExternalDamage,
			&inspection.IsChargeGaugeNormal,
			&inspection.IsReplaced,
			&inspection.AreMaintenanceRecordsComplete,
			&inspection.WorkOrderRequired,
			&inspection.InspectionStatus,
			&inspection.Notes,
		)
		if err != nil {
			return nil, err
		}

		inspections = append(inspections, inspection)
	}

	return inspections, nil
}

func (db *DB) GetInspectionByID(inspectionID int) (*models.Inspection, error) {
	query := `
	SELECT edi.emergencydeviceinspectionid, edi.emergencydeviceid, ed.serialnumber, edi.userid, u.username, edi.inspectiondatetime AT TIME ZONE 'Pacific/Auckland' AS lastinspectiondate_nzdt, edi.createdat AT TIME ZONE 'Pacific/Auckland' AS createdat_nzdt,
		   edi.IsConspicuous, edi.IsAccessible, edi.IsAssignedLocation, edi.IsSignVisible, edi.IsAntiTamperDeviceIntact,
		   edi.IsSupportBracketSecure, edi.AreOperatingInstructionsClear, edi.IsMaintenanceTagAttached,
		   edi.isNoExternalDamage, edi.IsChargeGaugeNormal, edi.IsReplaced, edi.AreMaintenanceRecordsComplete, edi.WorkOrderRequired,
		   edi.InspectionStatus, edi.Notes
	FROM emergency_device_inspectionT edi
	JOIN userT u ON edi.userid = u.userid
	JOIN emergency_deviceT ed ON edi.emergencydeviceid = ed.emergencydeviceid
	WHERE edi.emergencydeviceinspectionid = $1
	`

	var inspection models.Inspection
	err := db.QueryRow(query, inspectionID).Scan(
		&inspection.EmergencyDeviceInspectionID,
		&inspection.EmergencyDeviceID,
		&inspection.SerialNumber,
		&inspection.UserID,
		&inspection.InspectorName, // Scans the `username` field into `InspectorName`
		&inspection.InspectionDateTime,
		&inspection.CreatedAt,
		&inspection.IsConspicuous,
		&inspection.IsAccessible,
		&inspection.IsAssignedLocation,
		&inspection.IsSignVisible,
		&inspection.IsAntiTamperDeviceIntact,
		&inspection.IsSupportBracketSecure,
		&inspection.AreOperatingInstructionsClear,
		&inspection.IsMaintenanceTagAttached,
		&inspection.IsNoExternalDamage,
		&inspection.IsChargeGaugeNormal,
		&inspection.IsReplaced,
		&inspection.AreMaintenanceRecordsComplete,
		&inspection.WorkOrderRequired,
		&inspection.InspectionStatus,
		&inspection.Notes,
	)

	if err != nil {
		return nil, err
	}

	return &inspection, nil
}

func (db *DB) AddInspection(inspection *models.Inspection) error {

	query := `
	INSERT INTO emergency_device_inspectionT (emergencydeviceid, userid, inspectiondatetime, IsConspicuous, IsAccessible, IsAssignedLocation, IsSignVisible, IsAntiTamperDeviceIntact, IsSupportBracketSecure, AreOperatingInstructionsClear, IsMaintenanceTagAttached, IsNoExternalDamage, IsChargeGaugeNormal, IsReplaced, AreMaintenanceRecordsComplete, WorkOrderRequired, InspectionStatus, Notes)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`
	insertStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, err = insertStmt.Exec(
		inspection.EmergencyDeviceID,
		inspection.UserID,
		inspection.InspectionDateTime,
		inspection.IsConspicuous.Bool,
		inspection.IsAccessible.Bool,
		inspection.IsAssignedLocation.Bool,
		inspection.IsSignVisible.Bool,
		inspection.IsAntiTamperDeviceIntact.Bool,
		inspection.IsSupportBracketSecure.Bool,
		inspection.AreOperatingInstructionsClear.Bool,
		inspection.IsMaintenanceTagAttached.Bool,
		inspection.IsNoExternalDamage.Bool,
		inspection.IsChargeGaugeNormal.Bool,
		inspection.IsReplaced.Bool,
		inspection.AreMaintenanceRecordsComplete.Bool,
		inspection.WorkOrderRequired.Bool,
		inspection.InspectionStatus,
		inspection.Notes.String,
	)

	if err != nil {
		return err
	}

	return nil
}
