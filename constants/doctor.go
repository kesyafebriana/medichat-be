package constants

const (
	DoctorSortByStartWorkDate = "start_work_date"
	DoctorSortByName          = "name"
	DoctorSortByPrice         = "price"
)

var (
	DoctorSortBys = map[string]bool{
		DoctorSortByStartWorkDate: true,
		DoctorSortByName:          true,
		DoctorSortByPrice:         true,
	}
)
