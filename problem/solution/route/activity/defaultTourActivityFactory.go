package activity

import (
	"gsprit/problem"
)

type DefaultTourActivityFactory struct {
}

func (dt *DefaultTourActivityFactory) CreateActivity(service problem.Service) problem.AbstractActivity {

	if service.JobType() == problem.JobTypePickupService {
		return NewPickupService(service)
	}

	if service.JobType() == problem.JobTypeDeliveryService {
		return NewDeliverService(service)
	}

	return dt.createDefaultServiceActivity(service)
}

func (dt *DefaultTourActivityFactory) createDefaultServiceActivity(service problem.Service) problem.AbstractActivity {
	if service.Location() == nil {
		return NewActWithoutStaticLocation(service)
	}

	return NewPickupService(service)
}
