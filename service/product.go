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

type productService struct {
	dataRepository domain.DataRepository
	cloud          util.CloudinaryProvider
}

type ProductServiceOpts struct {
	DataRepository domain.DataRepository
	Cloud          util.CloudinaryProvider
}

func NewProductService(opts ProductServiceOpts) *productService {
	return &productService{
		dataRepository: opts.DataRepository,
		cloud:          opts.Cloud,
	}
}

func (s *productService) CreateCategoryLevelOne(ctx context.Context, product domain.Product, file *multipart.File) (domain.Category, error) {
	productRepo := s.dataRepository.ProductRepository()
	product.Name = strings.TrimSpace(strings.ToLower(product.Name))
	product.Slug = util.GenerateSlug(product.Name)

	c, err := productRepo.GetByName(ctx, category.Name)
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

	savedCategory, err := productRepo.Add(ctx, category)
	if err != nil {
		return domain.Category{}, apperror.Wrap(err)
	}

	return savedCategory, nil
}

func (s *productService) GetCategories(ctx context.Context, query domain.CategoriesQuery) ([]domain.CategoryWithParentName, domain.PageInfo, error) {
	productRepo := s.dataRepository.CategoryRepository()

	categories, err := productRepo.GetCategoriesWithParentName(ctx, query)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	pageInfo, err := productRepo.GetPageInfo(ctx, query)
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

func (s *productService) GetCategoriesHierarchy(ctx context.Context, query domain.CategoriesQuery) ([]domain.Category, error) {
	productRepo := s.dataRepository.CategoryRepository()

	categories, err := productRepo.GetCategories(ctx, query)
	if err != nil {
		return nil, apperror.Wrap(err)
	}

	return categories, nil
}

func (s *productService) GetCategoryBySlug(ctx context.Context, slug string) (domain.CategoryWithParentName, error) {
	productRepo := s.dataRepository.CategoryRepository()

	categories, err := productRepo.GetBySlugWithParentName(ctx, slug)
	if err != nil {
		return domain.CategoryWithParentName{}, apperror.Wrap(err)
	}

	return categories, nil
}

func (s *productService) DeleteCategory(ctx context.Context, slug string) error {
	productRepo := s.dataRepository.CategoryRepository()

	c, err := productRepo.GetBySlug(ctx, slug)
	if err != nil {
		return apperror.Wrap(err)
	}

	if c.ParentID != nil {
		err = productRepo.SoftDeleteBySlug(ctx, slug)
		if err != nil {
			return apperror.Wrap(err)
		}
		return nil
	}

	q := domain.DefaultCategoriesQuery()
	q.ParentSlug = slug
	childs, err := productRepo.GetCategoriesWithParentName(ctx, q)
	if err != nil {
		return apperror.Wrap(err)
	}

	slugs := make([]string, len(childs)+1)
	slugs[0] = slug
	for i := 0; i < len(childs); i++ {
		slugs[i+1] = childs[i].Category.Slug
	}

	err = productRepo.BulkSoftDeleteBySlug(ctx, slugs)
	if err != nil {
		return apperror.Wrap(err)
	}
	return nil
}

func (s *productService) UpdateCategory(ctx context.Context, category domain.Category, file *multipart.File) (domain.Category, error) {
	productRepo := s.dataRepository.CategoryRepository()

	c, err := productRepo.GetBySlug(ctx, category.Slug)
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
		cParent, err := productRepo.GetById(ctx, *category.ParentID)
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
	updatedCategory, err := productRepo.Update(ctx, category)
	if err != nil {
		return domain.Category{}, apperror.Wrap(err)
	}

	return updatedCategory, nil
}
