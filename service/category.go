package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
)

type categoryService struct {
	dataRepository domain.DataRepository
}

type CategoryServiceOpts struct {
	DataRepository domain.DataRepository
}

func NewCategoryService(opts CategoryServiceOpts) *categoryService {
	return &categoryService{
		dataRepository: opts.DataRepository,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, category domain.Category) (domain.Category, error) {
	categoryRepo := s.dataRepository.CategoryRepository()

	c, err := categoryRepo.GetByName(ctx, category.Name)
	if err != nil && !apperror.IsErrorCode(err, apperror.CodeNotFound) {
		return domain.Category{}, apperror.Wrap(err)
	}

	if c.Name == category.Name {
		return domain.Category{}, apperror.NewAlreadyExists("category")
	}

	if category.ParentID != nil {
		_, err := categoryRepo.GetById(ctx, *category.ParentID)
		if err != nil {
			return domain.Category{}, apperror.Wrap(err)
		}
	}

	savedCategory, err := categoryRepo.Add(ctx, category)
	if err != nil {
		return domain.Category{}, apperror.Wrap(err)
	}

	return savedCategory, nil
}

func (s *categoryService) GetCategories(ctx context.Context, query domain.CategoriesQuery) ([]domain.CategoryWithParentName, domain.PageInfo, error) {
	categoryRepo := s.dataRepository.CategoryRepository()

	categories, err := categoryRepo.GetCategoriesWithParentName(ctx, query)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	pageInfo, err := categoryRepo.GetPageInfo(ctx, query)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	pageInfo.ItemsPerPage = int(query.Limit)
	if query.Limit == 0 {
		pageInfo.ItemsPerPage = len(categories)
	}
	pageInfo.PageCount = (int(pageInfo.ItemCount) + pageInfo.ItemsPerPage - 1) / pageInfo.ItemsPerPage

	return categories, pageInfo, nil
}

func (s *categoryService) GetCategoriesHierarchy(ctx context.Context, query domain.CategoriesQuery) ([]domain.Category, error) {
	categoryRepo := s.dataRepository.CategoryRepository()

	categories, err := categoryRepo.GetCategories(ctx, query)
	if err != nil {
		return nil, apperror.Wrap(err)
	}

	return categories, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id int64) error {
	categoryRepo := s.dataRepository.CategoryRepository()

	c, err := categoryRepo.GetById(ctx, id)
	if err != nil {
		return apperror.Wrap(err)
	}

	if c.ParentID != nil {
		err = categoryRepo.SoftDeleteById(ctx, id)
		if err != nil {
			return apperror.Wrap(err)
		}
		return nil
	}

	childs, err := categoryRepo.GetCategories(ctx, domain.CategoriesQuery{
		ParentId: &id,
	})
	if err != nil {
		return apperror.Wrap(err)
	}

	ids := make([]int64, len(childs)+1)
	ids[0] = id
	for i := 0; i < len(childs); i++ {
		ids[i+1] = childs[i].ID
	}

	err = categoryRepo.BulkSoftDelete(ctx, ids)
	if err != nil {
		return apperror.Wrap(err)
	}
	return nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, category domain.Category) (domain.Category, error) {
	categoryRepo := s.dataRepository.CategoryRepository()

	_, err := categoryRepo.GetById(ctx, category.ID)
	if err != nil {
		return domain.Category{}, apperror.Wrap(err)
	}

	updatedCategory, err := categoryRepo.Update(ctx, category)
	if err != nil {
		return domain.Category{}, apperror.Wrap(err)
	}

	return updatedCategory, nil
}
