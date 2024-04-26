package constants

const (
	DoctorSortByYearExperience = "year_experience"
	DoctorSortByName           = "name"
	DoctorSortByPrice          = "price"
)

var (
	DoctorSortBys = map[string]bool{
		DoctorSortByYearExperience: true,
		DoctorSortByName:           true,
		DoctorSortByPrice:          true,
	}
)
