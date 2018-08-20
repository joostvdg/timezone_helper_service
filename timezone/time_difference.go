package timezone

type TimeDifference struct {
	LocationFrom string		`json:"location_from"`
	LocationFromTime string `json:"location_from_time"`
	LocationToo string      `json:"location_too"`
	LocationTooTime string  `json:"location_too_time"`
	TimeDifference int      `json:"time_difference"`
}
