package handler

import (
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/dto"
	"net/http"
	"strings"

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

func (h *CategoryHandler) CreateCategoryLevelOne(ctx *gin.Context) {
	var req dto.CreateCategoryRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	category, err := h.categorySrv.CreateCategory(ctx, domain.Category{
		Name: req.Name,
	})

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, dto.ResponseCreated(dto.NewCategoryResponse(category)))
}

func (h *CategoryHandler) CreateCategoryLevelTwo(ctx *gin.Context) {
	var req dto.CreateCategoryRequest
	var params dto.CategoryParams

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

	category, err := h.categorySrv.CreateCategory(ctx, domain.Category{
		Name:     req.Name,
		ParentID: &params.ID,
	})
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, dto.ResponseCreated(dto.NewCategoryResponse(category)))
}

func (h *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	var params dto.CategoryParams

	err := ctx.ShouldBindUri(&params)
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	err = h.categorySrv.DeleteCategory(ctx, params.ID)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	var req dto.UpdateCategoryRequest
	var params dto.CategoryParams

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

	category, err := h.categorySrv.UpdateCategory(ctx, domain.Category{
		ID:       params.ID,
		Name:     strings.ToLower(req.Name),
		ParentID: req.ParentId,
	})
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseOk(dto.NewCategoryResponse(category)))
}
