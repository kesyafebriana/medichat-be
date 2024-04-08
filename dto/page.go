package dto

import "medichat-be/domain"

type PageInfoResponse struct {
	CurrentPage  int   `json:"current_page"`
	ItemsPerPage int   `json:"items_per_page"`
	ItemCount    int64 `json:"item_count"`
	PageCount    int   `json:"page_count"`
}

func NewPageInfoResponse(i domain.PageInfo) PageInfoResponse {
	return PageInfoResponse{
		CurrentPage:  i.CurrentPage,
		ItemsPerPage: i.ItemsPerPage,
		ItemCount:    i.ItemCount,
		PageCount:    i.PageCount,
	}
}
