package models

import "database/sql"

// SiteT represents the sites in the system
type Site struct {
	SiteID           int            `json:"site_id"`
	SiteName         string         `json:"site_name"`
	SiteAddress      string         `json:"site_address"`
	SiteMapImagePath sql.NullString `json:"site_map_image_path"`
}
