package distance

import "math"

type Distance struct {
	Latitude1  float64
	Latitude2  float64
	Longitude1 float64
	Longitude2 float64
}

func (d *Distance) getDistance() float64 {

	d.Longitude1 = d.Longitude1 * math.Pi / 180
	d.Longitude2 = d.Longitude2 * math.Pi / 180

	d.Latitude1 = d.Latitude1 * math.Pi / 180
	d.Latitude2 = d.Latitude2 * math.Pi / 180

	dlon := d.Longitude2 - d.Longitude1
	dlat := d.Latitude2 - d.Latitude1
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(d.Latitude1)*math.Cos(d.Latitude2)*math.Pow(math.Sin(dlon/2), 2)

	return 2 * math.Asin(math.Sqrt(a))
}

func (d *Distance) GetDistanceOnMeter() float64 {
	if d.Latitude1 == 0 && d.Longitude1 == 0 {
		return 0
	}
	return d.getDistance() * 6371.0 * 1000
}
