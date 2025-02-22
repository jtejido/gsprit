package driver

import "math"

type DriverImpl struct {
	id            string
	earliestStart float64
	latestEnd     float64
	home          string
}

func NewDriver(id string) *DriverImpl {
	return &DriverImpl{
		id:        id,
		latestEnd: math.MaxFloat64,
	}
}

func (d *DriverImpl) Id() string {
	return d.id
}

func (d *DriverImpl) EarliestStart() float64 {
	return d.earliestStart
}

func (d *DriverImpl) SetEarliestStart(earliestStart float64) {
	d.earliestStart = earliestStart
}

func (d *DriverImpl) LatestEnd() float64 {
	return d.latestEnd
}

func (d *DriverImpl) SetLatestEnd(latestEnd float64) {
	d.latestEnd = latestEnd
}

func (d *DriverImpl) SetHomeLocation(locationId string) {
	d.home = locationId
}

func (d *DriverImpl) HomeLocation() string {
	return d.home
}

type NoDriver struct {
	DriverImpl
}

func NewNoDriver() *NoDriver {
	return &NoDriver{
		DriverImpl: DriverImpl{
			id:        "noDriver",
			latestEnd: math.MaxFloat64,
		},
	}
}
