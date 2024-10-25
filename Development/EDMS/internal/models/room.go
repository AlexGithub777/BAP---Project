package models

// RoomT represents rooms in each building
type Room struct {
	RoomID       int    `json:"room_id"`
	BuildingID   int    `json:"building_id"`
	RoomCode     string `json:"room_code"`
	BuildingCode string `json:"building_code"`
	SiteName     string `json:"site_name"`
	SiteID       int    `json:"site_id"`
}

type RoomDto struct {
	RoomID       string `json:"room_id"`
	BuildingID   string `json:"building_id"`
	RoomCode     string `json:"room_code"`
	BuildingCode string `json:"building_code"`
	SiteName     string `json:"site_name"`
	SiteID       string `json:"site_id"`
}
