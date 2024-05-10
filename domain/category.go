package domain

import (
	"context"
	"mime/multipart"
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
	PhotoUrl *string
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

func DefaultCategoriesQuery() CategoriesQuery {
	return CategoriesQuery{
		Page:     1,
		SortBy:   CategorySortById,
		SortType: "ASC",
	}
}

type CategoryRepository interface {
	GetCategoriesWithParentName(ctx context.Context, query CategoriesQuery) ([]CategoryWithParentName, error)
	GetCategories(ctx context.Context, query CategoriesQuery) ([]Category, error)
	GetPageInfo(ctx context.Context, query CategoriesQuery) (PageInfo, error)
	GetByName(ctx context.Context, name string) (Category, error)
	GetById(ctx context.Context, id int64) (Category, error)
	GetBySlug(ctx context.Context, slug string) (Category, error)
	GetBySlugWithParentName(ctx context.Context, slug string) (CategoryWithParentName, error)

	Add(ctx context.Context, category Category) (Category, error)
	Update(ctx context.Context, category Category) (Category, error)
	SoftDeleteBySlug(ctx context.Context, slug string) error
	BulkSoftDeleteBySlug(ctx context.Context, slug []string) error
}

type CategoryService interface {
	CreateCategoryLevelOne(ctx context.Context, category Category, file *multipart.File) (Category, error)
	CreateCategoryLevelTwo(ctx context.Context, category Category, parentSlug string) (CategoryWithParentName, error)
	GetCategories(ctx context.Context, query CategoriesQuery) ([]CategoryWithParentName, PageInfo, error)
	GetCategoriesHierarchy(ctx context.Context, query CategoriesQuery) ([]Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (CategoryWithParentName, error)
	DeleteCategory(ctx context.Context, slug string) error
	UpdateCategory(ctx context.Context, category Category, file *multipart.File) (Category, error)
}
