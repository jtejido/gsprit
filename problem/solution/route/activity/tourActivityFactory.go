package activity

import "gsprit/problem"

type TourActivityFactory interface {
	CreateActivity(service problem.Service) problem.AbstractActivity
}
