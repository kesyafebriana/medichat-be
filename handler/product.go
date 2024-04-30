package handler

import (
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/domain"
	"medichat-be/dto"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productsrv domain.ProductService
	domain      string
}

type ProductHandlerOpts struct {
	ProductSrv domain.ProductService
}

func NewProductHandler(opts ProductHandlerOpts) *ProductHandler {
	return &ProductHandler{
		productsrv: opts.ProductSrv,
	}
}

func (h *ProductHandler) GetProducts(ctx *gin.Context) {
	var query dto.GetProductsQuery

	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}
	products, pageInfo, err := h.productsrv.GetProducts(ctx,query.ToProductsQuery() )
	
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}
	
	ctx.JSON(http.StatusOK, dto.ResponseOk(dto.NewProductsResponse(products, pageInfo)))
}

func (h *ProductHandler) GetProductBySlug(ctx *gin.Context) {
	var params dto.CategorySlugParams

	err := ctx.ShouldBindUri(&params)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	product, err := h.productsrv.GetProduct(ctx, params.Slug)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseOk(dto.NewProductResponse(product)))
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var form dto.CreateProductForm

	if err := ctx.ShouldBind(&form); err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	var file *multipart.File
	if form.Picture != nil {
		fileType := form.Picture.Header.Get("Content-Type")
		if fileType != constants.ImageJpeg && fileType != constants.ImageJpg && fileType != constants.ImagePng {
			ctx.Error(apperror.NewRestrictredFileType(constants.ImageJpeg, constants.ImageJpg, constants.ImagePng))
			ctx.Abort()
			return
		}
		f, err := form.Picture.Open()
		if err != nil {
			ctx.Error(apperror.NewInternal(err))
			ctx.Abort()
			return
		}
		file = &f
		defer f.Close()
	}

	product, err := h.productsrv.CreateProduct(ctx, domain.AddProductRequest{
		Name: form.Name,
		ProductCategoryId: form.CategoryId,
		GenericName: form.GenericName,
		Content: form.Content,
		Manufacturer: form.Manufacturer,
		Description: form.Description,
		ProductClassification: form.ProductClassification,
		ProductForm: form.ProductForm,
		UnitInPack: form.UnitInPack,
		SellingUnit: form.SellingUnit,
		Weight: form.Weight,
		Height: form.Height,
		Length: form.Length,
		Width: form.Width,
	}, file)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, dto.ResponseCreated(dto.NewProductResponse(product)))
}

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	var params dto.ProductSlugParams

	err := ctx.ShouldBindUri(&params)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.productsrv.DeleteProducts(ctx, params.Slug)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	var form dto.UpdateProductForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	var file *multipart.File
	if form.Picture != nil {
		fileType := form.Picture.Header.Get("Content-Type")
		if fileType != constants.ImageJpeg && fileType != constants.ImageJpg && fileType != constants.ImagePng {
			ctx.Error(apperror.NewRestrictredFileType(constants.ImageJpeg, constants.ImageJpg, constants.ImagePng))
			ctx.Abort()
			return
		}
		f, err := form.Picture.Open()
		if err != nil {
			ctx.Error(apperror.NewInternal(err))
			ctx.Abort()
			return
		}
		file = &f
		defer f.Close()
	}

	product, err := h.productsrv.UpdateProduct(ctx, form.Slug,domain.UpdateProductRequest{
		Name: form.Name,
		ProductCategoryId: form.CategoryId,
		GenericName: form.GenericName,
		Content: form.Content,
		Manufacturer: form.Manufacturer,
		Description: form.Description,
		ProductClassification: form.ProductClassification,
		ProductForm: form.ProductForm,
		UnitInPack: form.UnitInPack,
		SellingUnit: form.SellingUnit,
		Weight: form.Weight,
		Height: form.Height,
		Length: form.Length,
		Width: form.Width,
	}, file)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, dto.ResponseOk(dto.NewProductResponse(product)))
}
