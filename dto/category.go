package dto

import (
	"medichat-be/constants"
	"medichat-be/domain"
	"mime/multipart"
	"sort"
)

type CreateCategoryForm struct {
	Name  string                `form:"name" binding:"required"`
	Image *multipart.FileHeader `form:"image"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryForm struct {
	Name     string                `form:"name"`
	ParentId *int64                `form:"parent_id" binding:"omitempty,numeric,min=1"`
	Image    *multipart.FileHeader `form:"image"`
}

type GetCategoriesQuery struct {
	Page       int64  `form:"page" binding:"numeric,omitempty,min=1"`
	Limit      int64  `form:"limit" binding:"numeric,omitempty,min=1"`
	Level      int64  `form:"level" binding:"numeric,omitempty,oneof=1 2"`
	SortBy     string `form:"sort_by" binding:"omitempty,oneof=name level parent"`
	SortType   string `form:"sort_type" binding:"omitempty,oneof=ASC DESC"`
	ParentSlug string `form:"parent_slug"`
	Term       string `form:"term"`
}

type CategorySlugParams struct {
	Slug string `uri:"slug" binding:"required"`
}

type CategoryResponse struct {
	ID       int64   `json:"id"`
	ParentID *int64  `json:"parent_id,omitempty"`
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	PhotoUrl *string `json:"photo_url,omitempty"`
}

type CategoryWithParentNameResponse struct {
	ID         int64   `json:"id"`
	ParentID   *int64  `json:"parent_id,omitempty"`
	Name       string  `json:"name"`
	Slug       string  `json:"slug"`
	PhotoUrl   *string `json:"photo_url,omitempty"`
	ParentName *string `json:"parent_name,omitempty"`
}

type CategoriesWithParentNameResponse struct {
	Categories []CategoryWithParentNameResponse `json:"categories"`
	PageInfo   PageInfoResponse                 `json:"page_info"`
}

type CategoriesResponse struct {
	Parent    CategoryResponse   `json:"parent"`
	Childrens []CategoryResponse `json:"childrens"`
}

func (q *GetCategoriesQuery) ToCategoriesQuery() domain.CategoriesQuery {
	var page int64 = q.Page
	var sortBy string = q.SortBy
	var sortType string = q.SortType
	if q.Page == 0 || q.Limit == 0 {
		page = 1
	}
	if q.SortBy == "" {
		sortBy = domain.CategorySortById
	}
	if q.SortType == "" {
		sortType = constants.SortAsc
	}
	if q.SortBy == domain.CategorySortByLevel {
		sortBy = "parent_id"
		if sortType == constants.SortAsc {
			sortType = constants.SortDesc
		} else {
			sortType = constants.SortAsc
		}
	}
	return domain.CategoriesQuery{
		Page:       page,
		Limit:      q.Limit,
		Level:      q.Level,
		Term:       q.Term,
		SortBy:     sortBy,
		SortType:   sortType,
		ParentSlug: q.ParentSlug,
	}
}

func NewCategoryResponse(c domain.Category) CategoryResponse {
	photoUrl := c.PhotoUrl
	if c.ParentID == nil && photoUrl == nil {
		t := constants.DefaultCategoryImageURL
		photoUrl = &t
	}
	return CategoryResponse{
		ID:       c.ID,
		ParentID: c.ParentID,
		Name:     c.Name,
		Slug:     c.Slug,
		PhotoUrl: photoUrl,
	}
}

func NewCategoryWithParentNameResponse(c domain.CategoryWithParentName) CategoryWithParentNameResponse {
	return CategoryWithParentNameResponse{
		ID:         c.Category.ID,
		ParentID:   c.Category.ParentID,
		Name:       c.Category.Name,
		ParentName: c.ParentName,
		Slug:       c.Category.Slug,
		PhotoUrl:   c.Category.PhotoUrl,
	}
}

func NewCategoriesHierarchyResponse(categories []domain.Category) []CategoriesResponse {
	res := []CategoriesResponse{}
	categoriesMap := map[int64][]CategoryResponse{}
	parentsMap := map[int64]*CategoryResponse{}
	childs := []CategoryResponse{}
	for i := 0; i < len(categories); i++ {
		cR := NewCategoryResponse(categories[i])
		if categories[i].ParentID != nil {
			childs = append(childs, cR)
			continue
		}
		if cR.PhotoUrl == nil {
			t := constants.DefaultCategoryImageURL
			cR.PhotoUrl = &t
		}
		parentsMap[cR.ID] = &cR
		categoriesMap[cR.ID] = []CategoryResponse{}
	}

	for i := 0; i < len(childs); i++ {
		categoriesMap[*childs[i].ParentID] = append(categoriesMap[*childs[i].ParentID], childs[i])
	}

	for k, v := range categoriesMap {
		res = append(res, CategoriesResponse{
			Parent:    *parentsMap[k],
			Childrens: v,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Parent.ID < res[j].Parent.ID
	})

	return res
}

func NewCategoriesWithParentNameResponse(categories []domain.CategoryWithParentName, pageInfo domain.PageInfo) CategoriesWithParentNameResponse {
	res := make([]CategoryWithParentNameResponse, len(categories))
	for i := 0; i < len(categories); i++ {
		if categories[i].Category.ParentID == nil && categories[i].Category.PhotoUrl == nil {
			t := constants.DefaultCategoryImageURL
			categories[i].Category.PhotoUrl = &t
		}
		res[i] = NewCategoryWithParentNameResponse(categories[i])
	}
	return CategoriesWithParentNameResponse{
		Categories: res,
		PageInfo:   NewPageInfoResponse(pageInfo),
	}
}
