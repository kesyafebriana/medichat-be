package service

import (
	"context"
	"fmt"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/util"
	"mime/multipart"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type categoryService struct {
	dataRepository domain.DataRepository
	cloud          util.CloudinaryProvider
}

type CategoryServiceOpts struct {
	DataRepository domain.DataRepository
	Cloud          util.CloudinaryProvider
}

func NewCategoryService(opts CategoryServiceOpts) *categoryService {
	return &categoryService{
		dataRepository: opts.DataRepository,
		cloud:          opts.Cloud,
	}
}

func (s *categoryService) CreateCategoryLevelOne(ctx context.Context, category domain.Category, file *multipart.File) (domain.Category, error) {
	categoryRepo := s.dataRepository.CategoryRepository()
	category.Name = strings.TrimSpace(strings.ToLower(category.Name))
	category.Slug = util.GenerateSlug(category.Name)

	c, err := categoryRepo.GetByName(ctx, category.Name)
	if err != nil && !apperror.IsErrorCode(err, apperror.CodeNotFound) {
		return domain.Category{}, apperror.Wrap(err)
	}

	if c.Name == category.Name {
		return domain.Category{}, apperror.NewAlreadyExists("category")
	}

	if file != nil {
		res, err := s.cloud.UploadImage(ctx, *file, uploader.UploadParams{})
		if err == nil {
			category.PhotoUrl = &res.SecureURL
		}
	}

	savedCategory, err := categoryRepo.Add(ctx, category)
	if err != nil {
		return domain.Category{}, apperror.Wrap(err)
	}

	return savedCategory, nil
}

func (s *categoryService) CreateCategoryLevelTwo(ctx context.Context, category domain.Category, parentSlug string) (domain.CategoryWithParentName, error) {
	categoryRepo := s.dataRepository.CategoryRepository()
	category.Name = strings.TrimSpace(strings.ToLower(category.Name))
	category.Slug = util.GenerateSlug(category.Name)

	c, err := categoryRepo.GetByName(ctx, category.Name)
	if err != nil && !apperror.IsErrorCode(err, apperror.CodeNotFound) {
		return domain.CategoryWithParentName{}, apperror.Wrap(err)
	}

	if c.Name == category.Name {
		return domain.CategoryWithParentName{}, apperror.NewAlreadyExists("category")
	}

	c, err = categoryRepo.GetBySlug(ctx, parentSlug)
	if err != nil {
		if apperror.IsErrorCode(err, apperror.CodeNotFound) {
			return domain.CategoryWithParentName{}, apperror.NewEntityNotFound(fmt.Sprintf(`category with slug %s`, parentSlug))
		}
		return domain.CategoryWithParentName{}, apperror.Wrap(err)
	}

	if c.ParentID != nil {
		return domain.CategoryWithParentName{}, apperror.NewCreateCategoryParentRestrict()
	}

	category.ParentID = &c.ID
	savedCategory, err := categoryRepo.Add(ctx, category)
	if err != nil {
		return domain.CategoryWithParentName{}, apperror.Wrap(err)
	}

	return domain.CategoryWithParentName{
		Category:   savedCategory,
		ParentName: &c.Name,
	}, nil

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

	if pageInfo.ItemsPerPage == 0 {
		pageInfo.PageCount = 0
	} else {
		pageInfo.PageCount = (int(pageInfo.ItemCount) + pageInfo.ItemsPerPage - 1) / pageInfo.ItemsPerPage
	}

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

func (s *categoryService) GetCategoryBySlug(ctx context.Context, slug string) (domain.CategoryWithParentName, error) {
	categoryRepo := s.dataRepository.CategoryRepository()

	categories, err := categoryRepo.GetBySlugWithParentName(ctx, slug)
	if err != nil {
		return domain.CategoryWithParentName{}, apperror.Wrap(err)
	}

	return categories, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, slug string) error {
	categoryRepo := s.dataRepository.CategoryRepository()

	c, err := categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return apperror.Wrap(err)
	}

	if c.ParentID != nil {
		err = categoryRepo.SoftDeleteBySlug(ctx, slug)
		if err != nil {
			return apperror.Wrap(err)
		}
		return nil
	}

	q := domain.DefaultCategoriesQuery()
	q.ParentSlug = slug
	childs, err := categoryRepo.GetCategoriesWithParentName(ctx, q)
	if err != nil {
		return apperror.Wrap(err)
	}

	slugs := make([]string, len(childs)+1)
	slugs[0] = slug
	for i := 0; i < len(childs); i++ {
		slugs[i+1] = childs[i].Category.Slug
	}

	err = categoryRepo.BulkSoftDeleteBySlug(ctx, slugs)
	if err != nil {
		return apperror.Wrap(err)
	}
	return nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, category domain.Category, file *multipart.File) (domain.Category, error) {
	categoryRepo := s.dataRepository.CategoryRepository()

	c, err := categoryRepo.GetBySlug(ctx, category.Slug)
	if err != nil {
		if apperror.IsErrorCode(err, apperror.CodeNotFound) {
			return domain.Category{}, apperror.NewEntityNotFound(fmt.Sprintf("category with slug %s", category.Slug))
		}
		return domain.Category{}, apperror.Wrap(err)
	}

	if category.Name == "" {
		category.Name = c.Name
	}

	if c.ParentID == nil {
		category.ParentID = nil
	}

	if category.ParentID == nil && c.ParentID != nil {
		category.ParentID = c.ParentID
	}

	if category.ParentID != nil && c.ParentID != nil {
		cParent, err := categoryRepo.GetById(ctx, *category.ParentID)
		if err != nil {
			if apperror.IsErrorCode(err, apperror.CodeNotFound) {
				return domain.Category{}, apperror.NewEntityNotFound(fmt.Sprintf("category with id %d", *category.ParentID))
			}
			return domain.Category{}, apperror.Wrap(err)
		}

		if cParent.ParentID != nil {
			return domain.Category{}, apperror.NewUpdateCategoryParentRestrict()
		}
	}

	if file != nil && c.ParentID == nil {
		res, err := s.cloud.UploadImage(ctx, *file, uploader.UploadParams{})
		if err == nil {
			category.PhotoUrl = &res.SecureURL
		}
	}

	category.ID = c.ID
	category.Name = strings.TrimSpace(strings.ToLower(category.Name))
	category.Slug = util.GenerateSlug(category.Name)
	updatedCategory, err := categoryRepo.Update(ctx, category)
	if err != nil {
		return domain.Category{}, apperror.Wrap(err)
	}

	return updatedCategory, nil
}
