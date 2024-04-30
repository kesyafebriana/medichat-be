package constants

const (
	SortAsc  = "asc"
	SortDesc = "desc"
)

var (
	SortOrders = map[string]bool{
		SortAsc:  true,
		SortDesc: true,
	}
)
