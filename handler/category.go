package handler

import (
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/domain"
	"medichat-be/dto"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categorySrv domain.CategoryService
	domain      string
}

type CategoryHandlerOpts struct {
	CategorySrv domain.CategoryService
	Domain      string
}

func NewCategoryHandler(opts CategoryHandlerOpts) *CategoryHandler {
	return &CategoryHandler{
		categorySrv: opts.CategorySrv,
		domain:      opts.Domain,
	}
}

func (h *CategoryHandler) GetCategoriesHierarchy(ctx *gin.Context) {
	query := dto.GetCategoriesQuery{}

	categories, err := h.categorySrv.GetCategoriesHierarchy(ctx, query.ToCategoriesQuery())
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseOk(dto.NewCategoriesHierarchyResponse(categories)))
}

func (h *CategoryHandler) GetCategories(ctx *gin.Context) {
	var query dto.GetCategoriesQuery

	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	categories, pageInfo, err := h.categorySrv.GetCategories(ctx, query.ToCategoriesQuery())
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseOk(dto.NewCategoriesWithParentNameResponse(categories, pageInfo)))
}

func (h *CategoryHandler) GetCategoryBySlug(ctx *gin.Context) {
	var params dto.CategorySlugParams

	err := ctx.ShouldBindUri(&params)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	category, err := h.categorySrv.GetCategoryBySlug(ctx, params.Slug)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseOk(dto.NewCategoryWithParentNameResponse(category)))
}

func (h *CategoryHandler) CreateCategoryLevelOne(ctx *gin.Context) {
	var form dto.CreateCategoryForm
	contentLength := ctx.Request.Header.Get("Content-Length")
	i, err := strconv.Atoi(contentLength)
	if err != nil {
		ctx.Error(apperror.NewInternal(err))
		ctx.Abort()
		return
	}

	if i > constants.MaxImageSize {
		ctx.Error(apperror.NewImageSizeExceeded("500kb"))
		ctx.Abort()
		return
	}

	if err := ctx.ShouldBind(&form); err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	var file *multipart.File
	if form.Image != nil {
		fileType := form.Image.Header.Get("Content-Type")
		if fileType != constants.ImageJpeg && fileType != constants.ImageJpg && fileType != constants.ImagePng {
			ctx.Error(apperror.NewRestrictredFileType(constants.ImageJpeg, constants.ImageJpg, constants.ImagePng))
			ctx.Abort()
			return
		}
		f, err := form.Image.Open()
		if err != nil {
			ctx.Error(apperror.NewInternal(err))
			ctx.Abort()
			return
		}
		file = &f
		defer f.Close()
	}

	category, err := h.categorySrv.CreateCategoryLevelOne(ctx, domain.Category{
		Name: form.Name,
	}, file)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, dto.ResponseCreated(dto.NewCategoryResponse(category)))
}

func (h *CategoryHandler) CreateCategoryLevelTwo(ctx *gin.Context) {
	var req dto.CreateCategoryRequest
	var params dto.CategorySlugParams

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = ctx.ShouldBindUri(&params)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	category, err := h.categorySrv.CreateCategoryLevelTwo(ctx, domain.Category{
		Name: req.Name,
	}, params.Slug)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, dto.ResponseCreated(dto.NewCategoryWithParentNameResponse(category)))
}

func (h *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	var params dto.CategorySlugParams

	err := ctx.ShouldBindUri(&params)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.categorySrv.DeleteCategory(ctx, params.Slug)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	var form dto.UpdateCategoryForm
	contentLength := ctx.Request.Header.Get("Content-Length")
	var params dto.CategorySlugParams

	err := ctx.ShouldBindUri(&params)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	i, err := strconv.Atoi(contentLength)
	if err != nil {
		ctx.Error(apperror.NewInternal(err))
		ctx.Abort()
		return
	}

	if i > constants.MaxImageSize {
		ctx.Error(apperror.NewImageSizeExceeded("500kb"))
		ctx.Abort()
		return
	}

	err = ctx.ShouldBind(&form)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	var file *multipart.File
	if form.Image != nil {
		fileType := form.Image.Header.Get("Content-Type")
		if fileType != constants.ImageJpeg && fileType != constants.ImageJpg && fileType != constants.ImagePng {
			ctx.Error(apperror.NewRestrictredFileType(constants.ImageJpeg, constants.ImageJpg, constants.ImagePng))
			ctx.Abort()
			return
		}
		f, err := form.Image.Open()
		if err != nil {
			ctx.Error(apperror.NewInternal(err))
			ctx.Abort()
			return
		}
		file = &f
		defer f.Close()
	}

	category, err := h.categorySrv.UpdateCategory(ctx, domain.Category{
		Name:     form.Name,
		ParentID: form.ParentId,
		Slug:     params.Slug,
	}, file)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseOk(dto.NewCategoryResponse(category)))
}
