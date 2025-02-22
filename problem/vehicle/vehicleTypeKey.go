package vehicle

import (
	"fmt"
	"gsprit/problem"
)

type VehicleTypeKey struct {
	problem.BaseTypeKey
	t             string
	startLocation string
	endLocation   string
	earliestStart float64
	latestEnd     float64
	skills        *problem.Skills
	returnToDepot bool
}

func NewVehicleTypeKey(typeID, startLocation, endLocation string, earliestStart, latestEnd float64, skills *problem.Skills, returnToDepot bool) *VehicleTypeKey {
	return &VehicleTypeKey{
		t:             typeID,
		startLocation: startLocation,
		endLocation:   endLocation,
		earliestStart: earliestStart,
		latestEnd:     latestEnd,
		skills:        skills,
		returnToDepot: returnToDepot,
	}
}

func (v *VehicleTypeKey) Type() string          { return v.t }
func (v *VehicleTypeKey) SetType(typeID string) { v.t = typeID }

func (v *VehicleTypeKey) StartLocation() string     { return v.startLocation }
func (v *VehicleTypeKey) SetStartLocation(s string) { v.startLocation = s }

func (v *VehicleTypeKey) EndLocation() string     { return v.endLocation }
func (v *VehicleTypeKey) SetEndLocation(e string) { v.endLocation = e }

func (v *VehicleTypeKey) EarliestStart() float64     { return v.earliestStart }
func (v *VehicleTypeKey) SetEarliestStart(e float64) { v.earliestStart = e }

func (v *VehicleTypeKey) LatestEnd() float64     { return v.latestEnd }
func (v *VehicleTypeKey) SetLatestEnd(l float64) { v.latestEnd = l }

func (v *VehicleTypeKey) Skills() *problem.Skills     { return v.skills }
func (v *VehicleTypeKey) SetSkills(s *problem.Skills) { v.skills = s }

func (v *VehicleTypeKey) ReturnToDepot() bool     { return v.returnToDepot }
func (v *VehicleTypeKey) SetReturnToDepot(r bool) { v.returnToDepot = r }

func (v *VehicleTypeKey) Equals(other problem.VehicleTypeKey) bool {
	if other == nil {
		return false
	}
	return v.t == other.Type() &&
		v.startLocation == other.StartLocation() &&
		v.endLocation == other.EndLocation() &&
		v.earliestStart == other.EarliestStart() &&
		v.latestEnd == other.LatestEnd() &&
		v.skills == other.Skills() &&
		v.returnToDepot == other.ReturnToDepot()
}

func (v *VehicleTypeKey) String() string {
	return fmt.Sprintf("%s_%s_%s_%.2f_%.2f", v.t, v.startLocation, v.endLocation, v.earliestStart, v.latestEnd)
}
