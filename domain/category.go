package domain

import (
	"context"
)

const (
	CategorySortById     = "id"
	CategorySortByName   = "name"
	CategorySortByLevel  = "level"
	CategorySortByParent = "parent"
)

type Category struct {
	ID       int64
	ParentID *int64
	Name     string
	Slug     string
}

type CategoryWithParentName struct {
	Category   Category
	ParentName *string
}

type CategoriesQuery struct {
	ParentId   *int64
	Page       int64
	Limit      int64
	Level      int64
	Term       string
	SortBy     string
	SortType   string
	ParentSlug string
}

type CategoryRepository interface {
	GetCategoriesWithParentName(ctx context.Context, query CategoriesQuery) ([]CategoryWithParentName, error)
	GetCategories(ctx context.Context, query CategoriesQuery) ([]Category, error)
	GetPageInfo(ctx context.Context, query CategoriesQuery) (PageInfo, error)
	GetByName(ctx context.Context, name string) (Category, error)
	GetById(ctx context.Context, id int64) (Category, error)

	Add(ctx context.Context, category Category) (Category, error)
	Update(ctx context.Context, category Category) (Category, error)
	SoftDeleteById(ctx context.Context, id int64) error
	BulkSoftDelete(ctx context.Context, ids []int64) error
}

type CategoryService interface {
	CreateCategory(ctx context.Context, category Category) (Category, error)
	GetCategories(ctx context.Context, query CategoriesQuery) ([]CategoryWithParentName, PageInfo, error)
	GetCategoriesHierarchy(ctx context.Context, query CategoriesQuery) ([]Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	UpdateCategory(ctx context.Context, category Category) (Category, error)
}
