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

func (s *productService) CreateProduct(ctx context.Context, request domain.AddProductRequest, file *multipart.File) (domain.Product, error) {
	productRepo := s.dataRepository.ProductRepository()
	categoryRepo := s.dataRepository.CategoryRepository()
	detailRepo := s.dataRepository.ProductDetailsRepository()

	product := domain.Product{}

	product.IsActive = true
	product.Name = strings.TrimSpace(strings.ToLower(request.Name))
	product.Slug = util.GenerateSlug(request.Name)

	product.KeyWord = product.Name + " " + request.Manufacturer + " " + request.Composition

	p, err := productRepo.GetByName(ctx, product.Name)
	if err != nil && !apperror.IsErrorCode(err, apperror.CodeNotFound) {
		return domain.Product{}, apperror.Wrap(err)
	}

	if p.Name == request.Name {
		return domain.Product{}, apperror.NewAlreadyExists("product")
	}

	cat, err := categoryRepo.GetById(ctx, request.ProductCategoryId)
	if err != nil && !apperror.IsErrorCode(err, apperror.CodeNotFound) {
		return domain.Product{}, apperror.Wrap(err)
	}

	product.ProductCategoryId = cat.ID

	if file != nil {
		res, err := s.cloud.UploadImage(ctx, *file, uploader.UploadParams{})
		if err == nil {
			product.Picture = &res.SecureURL
		}
	}

	detail := domain.ProductDetails{
		GenericName:           request.GenericName,
		Content:               request.Content,
		Manufacturer:          request.Manufacturer,
		Composition:           request.Composition,
		Description:           request.Description,
		ProductClassification: request.ProductClassification,
		ProductForm:           request.ProductForm,
		UnitInPack:            request.UnitInPack,
		SellingUnit:           request.SellingUnit,
		Weight:                request.Weight,
		Height:                request.Height,
		Length:                request.Length,
		Width:                 request.Width,
	}

	det, err := detailRepo.Add(ctx, detail)
	if err != nil && !apperror.IsErrorCode(err, apperror.CodeNotFound) {
		return domain.Product{}, apperror.Wrap(err)
	}

	product.ProductDetailId = det.ID

	prod, err := productRepo.Add(ctx, product)
	if err != nil && !apperror.IsErrorCode(err, apperror.CodeNotFound) {
		return domain.Product{}, apperror.Wrap(err)
	}

	if err != nil {
		return domain.Product{}, apperror.Wrap(err)
	}

	return prod, nil
}

func (s *productService) GetProduct(ctx context.Context, slug string) (domain.Product, domain.ProductDetails, domain.CategoryWithParentName, error) {
	productRepo := s.dataRepository.ProductRepository()
	productDetailRepo := s.dataRepository.ProductDetailsRepository()
	categoryRepo := s.dataRepository.CategoryRepository()

	products, err := productRepo.GetBySlug(ctx, slug)
	if err != nil {
		return domain.Product{}, domain.ProductDetails{}, domain.CategoryWithParentName{}, apperror.Wrap(err)
	}

	productDetail, err := productDetailRepo.GetById(ctx, products.ProductDetailId)
	if err != nil {
		return domain.Product{}, domain.ProductDetails{}, domain.CategoryWithParentName{}, apperror.Wrap(err)
	}

	c, err := categoryRepo.GetById(ctx, products.ProductCategoryId)
	if err != nil {
		return domain.Product{}, domain.ProductDetails{}, domain.CategoryWithParentName{}, apperror.Wrap(err)
	}

	category, err := categoryRepo.GetBySlugWithParentName(ctx, c.Slug)

	return products, productDetail, category, nil
}

func (s *productService) GetProductLocation(ctx context.Context, query domain.ProductsQuery) ([]domain.Product, domain.PageInfo, error) {
	productRepo := s.dataRepository.ProductRepository()

	products, err := productRepo.GetProductsFromArea(ctx, query)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	pageInfo, err := productRepo.GetPageInfoFromArea(ctx, query)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	pageInfo.ItemsPerPage = int(query.Limit)
	if query.Limit == 0 {
		pageInfo.ItemsPerPage = len(products)
	}

	if pageInfo.ItemsPerPage == 0 {
		pageInfo.PageCount = 0
	} else {
		pageInfo.PageCount = (int(pageInfo.ItemCount) + pageInfo.ItemsPerPage - 1) / pageInfo.ItemsPerPage
	}

	return products, pageInfo, nil
}

func (s *productService) GetProducts(ctx context.Context, query domain.ProductsQuery) ([]domain.Product, domain.PageInfo, error) {
	productRepo := s.dataRepository.ProductRepository()
	categoryRepo := s.dataRepository.CategoryRepository()

	if query.CategorySlug != nil {
		category, err := categoryRepo.GetBySlug(ctx, *query.CategorySlug)
		if err != nil {
			return nil, domain.PageInfo{}, apperror.Wrap(err)
		}

		query.CategoryID = &category.ID
	}

	products, err := productRepo.GetProducts(ctx, query)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	pageInfo, err := productRepo.GetPageInfo(ctx, query)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	pageInfo.ItemsPerPage = int(query.Limit)
	if query.Limit == 0 {
		pageInfo.ItemsPerPage = len(products)
	}

	if pageInfo.ItemsPerPage == 0 {
		pageInfo.PageCount = 0
	} else {
		pageInfo.PageCount = (int(pageInfo.ItemCount) + pageInfo.ItemsPerPage - 1) / pageInfo.ItemsPerPage
	}

	return products, pageInfo, nil
}

func (s *productService) DeleteProducts(ctx context.Context, slug string) error {
	productRepo := s.dataRepository.ProductRepository()

	_, err := productRepo.GetBySlug(ctx, slug)
	if err != nil {
		return apperror.Wrap(err)
	}

	err = productRepo.SoftDeleteBySlug(ctx, slug)
	if err != nil {
		return apperror.Wrap(err)
	}
	return nil
}

func (s *productService) UpdateProduct(ctx context.Context, slug string, request domain.UpdateProductRequest, file *multipart.File) (domain.Product, error) {
	productRepo := s.dataRepository.ProductRepository()

	detailRepo := s.dataRepository.ProductDetailsRepository()

	prod, err := productRepo.GetBySlug(ctx, slug)
	if err != nil {
		if apperror.IsErrorCode(err, apperror.CodeNotFound) {
			return domain.Product{}, apperror.NewEntityNotFound(fmt.Sprintf("products with slug %s", slug))
		}
		return domain.Product{}, apperror.Wrap(err)
	}

	detailId := prod.ProductDetailId
	detail, err := detailRepo.GetById(ctx, detailId)
	if err != nil {
		if apperror.IsErrorCode(err, apperror.CodeNotFound) {
			return domain.Product{}, apperror.NewEntityNotFound(fmt.Sprintf("products with slug %s", slug))
		}
		return domain.Product{}, apperror.Wrap(err)
	}

	prod.Name = strings.TrimSpace(strings.ToLower(request.Name))

	if file != nil {
		res, err := s.cloud.UploadImage(ctx, *file, uploader.UploadParams{})
		if err == nil {
			prod.Picture = &res.SecureURL
		}
	}

	prod.KeyWord = prod.Name + " " + detail.GenericName
	prod.Slug = util.GenerateSlug(prod.Name)

	detail.GenericName = request.GenericName
	detail.Content = request.Content
	detail.Manufacturer = request.Manufacturer
	detail.Description = request.Description
	detail.ProductClassification = request.ProductClassification
	detail.ProductForm = request.ProductForm
	detail.UnitInPack = request.UnitInPack
	detail.SellingUnit = request.SellingUnit
	detail.Weight = request.Weight
	detail.Height = request.Height
	detail.Length = request.Length
	detail.Width = request.Width
	updatedDetail, err := detailRepo.Update(ctx, detail)
	if err != nil {
		return domain.Product{}, apperror.Wrap(err)
	}
	prod.ProductDetailId = updatedDetail.ID
	updatedProduct, err := productRepo.Update(ctx, prod)
	if err != nil {
		return domain.Product{}, apperror.Wrap(err)
	}
	return updatedProduct, nil
}
