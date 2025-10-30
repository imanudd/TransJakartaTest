package domain

type GetVehicleHistoryRequest struct {
	Start int64 `json:"start" form:"start"`
	End   int64 `json:"end" form:"end"`
}

type Vehicle struct {
	VehicleId   string `json:"vehicle_id"`
	VehicleType string `json:"vehicle_type"`
	VehicleCode string `json:"vehicle_code"`
}

type VehicleLocation struct {
	VehicleId string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}
