package models

// BuildingT represents the buildings in each site
type Building struct {
	BuildingID   int    `json:"building_id"`
	SiteID       int    `json:"site_id"`
	BuildingCode string `json:"building_code"`
	SiteName     string `json:"site_name"`
}
