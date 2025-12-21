package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/javapub/mini-study/mini-study-backend/internal/dto"
	"github.com/javapub/mini-study/mini-study-backend/internal/middleware"
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/service"
	"github.com/javapub/mini-study/mini-study-backend/internal/utils"
)

// ContentHandler handles content related endpoints.
type ContentHandler struct {
	service *service.ContentService
}

// NewContentHandler creates a content handler.
func NewContentHandler(contentService *service.ContentService) *ContentHandler {
	return &ContentHandler{service: contentService}
}

// ListCategories godoc
// @Summary 查询可见内容分类
// @Description 返回当前登录用户基于角色可访问的内容分类列表
// @Tags 内容
// @Security Bearer
// @Produce json
// @Success 200 {object} utils.Response{data=[]dto.ContentCategoryResponse}
// @Failure 401 {object} utils.Response
// @Router /api/v1/contents/categories [get]
func (h *ContentHandler) ListCategories(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	categories, counts, err := h.service.ListCategoriesWithCount(userID)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	resp := make([]dto.ContentCategoryResponse, 0, len(categories))
	for i, item := range categories {
		count := int64(0)
		if i < len(counts) {
			count = counts[i]
		}
		resp = append(resp, dto.ContentCategoryResponse{
			ID:        item.ID,
			Name:      item.Name,
			RoleScope: item.RoleScope,
			SortOrder: item.SortOrder,
			CoverURL:  item.CoverURL,
			Count:     count,
		})
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

func (h *ContentHandler) AdminListCategories(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	categories, counts, err := h.service.AdminListCategories(adminID)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	resp := make([]dto.ContentCategoryResponse, 0, len(categories))
	for i, item := range categories {
		count := int64(0)
		if i < len(counts) {
			count = counts[i]
		}
		resp = append(resp, dto.ContentCategoryResponse{
			ID:        item.ID,
			Name:      item.Name,
			RoleScope: item.RoleScope,
			SortOrder: item.SortOrder,
			CoverURL:  item.CoverURL,
			Count:     count,
		})
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// ListPublishedContents godoc
// @Summary 查询已发布内容
// @Description 根据分类与类型筛选当前用户可访问的已发布内容
// @Tags 内容
// @Security Bearer
// @Produce json
// @Param category_id query int false "分类ID"
// @Param type query string false "内容类型(doc/video/article)"
// @Success 200 {object} utils.Response{data=[]dto.ContentResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/contents [get]
func (h *ContentHandler) ListPublishedContents(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var query dto.PublishedContentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	contents, err := h.service.ListPublished(userID, query.CategoryID, query.Type)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	utils.NewSuccessResponse(h.toContentResponses(contents)).JSON(c)
}

// GetContentDetail godoc
// @Summary 查询内容详情
// @Description 获取指定已发布内容的详细信息
// @Tags 内容
// @Security Bearer
// @Produce json
// @Param id path int true "内容ID"
// @Success 200 {object} utils.Response{data=dto.ContentResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/contents/{id} [get]
func (h *ContentHandler) GetContentDetail(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	idParam := c.Param("id")
	contentID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的内容ID").JSON(c)
		return
	}

	content, err := h.service.GetPublishedDetail(userID, uint(contentID))
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	utils.NewSuccessResponse(h.toContentResponse(content)).JSON(c)
}

// AdminListContents godoc
// @Summary 管理员查询内容列表
// @Description 管理员可按分类、类型与状态筛选内容
// @Tags 管理后台-内容
// @Security Bearer
// @Produce json
// @Param category_id query int false "分类ID"
// @Param type query string false "内容类型(doc/video/article)"
// @Param status query string false "内容状态(draft/published/offline)"
// @Success 200 {object} utils.Response{data=[]dto.ContentResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/contents [get]
func (h *ContentHandler) AdminListContents(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var query dto.AdminListContentRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	contents, err := h.service.AdminListContents(adminID, query)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(h.toContentResponses(contents)).JSON(c)
}

// AdminCreateContent godoc
// @Summary 管理员创建内容
// @Description 管理员上传文件并设置分类、可见角色后创建内容
// @Tags 管理后台-内容
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body dto.AdminCreateContentRequest true "内容信息"
// @Success 200 {object} utils.Response{data=dto.ContentResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/contents [post]
func (h *ContentHandler) AdminCreateContent(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	var req dto.AdminCreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	content, err := h.service.AdminCreateContent(adminID, req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(h.toContentResponse(content)).JSON(c)
}

func (h *ContentHandler) AdminUpdateCategory(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	idParam := c.Param("id")
	categoryID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的分类ID").JSON(c)
		return
	}

	var req dto.AdminUpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	category, err := h.service.AdminUpdateCategory(adminID, uint(categoryID), req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	resp := dto.ContentCategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		RoleScope: category.RoleScope,
		SortOrder: category.SortOrder,
		CoverURL:  category.CoverURL,
		Count:     0,
	}
	utils.NewSuccessResponse(resp).JSON(c)
}

// AdminUpdateContent godoc
// @Summary 管理员更新内容
// @Description 管理员可修改内容基础信息、状态与可见范围
// @Tags 管理后台-内容
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "内容ID"
// @Param body body dto.AdminUpdateContentRequest true "内容信息"
// @Success 200 {object} utils.Response{data=dto.ContentResponse}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/admin/contents/{id} [put]
func (h *ContentHandler) AdminUpdateContent(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		utils.NewErrorResponse(http.StatusUnauthorized, "未登录").JSON(c)
		return
	}

	idParam := c.Param("id")
	contentID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, "非法的内容ID").JSON(c)
		return
	}

	var req dto.AdminUpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}

	content, err := h.service.AdminUpdateContent(adminID, uint(contentID), req)
	if err != nil {
		utils.NewErrorResponse(http.StatusBadRequest, err.Error()).JSON(c)
		return
	}
	utils.NewSuccessResponse(h.toContentResponse(content)).JSON(c)
}

func (h *ContentHandler) toContentResponses(contents []model.Content) []dto.ContentResponse {
	resp := make([]dto.ContentResponse, 0, len(contents))
	for idx := range contents {
		resp = append(resp, h.toContentResponse(&contents[idx]))
	}
	return resp
}

func (h *ContentHandler) toContentResponse(content *model.Content) dto.ContentResponse {
	categoryName := ""
	if content.Category.ID != 0 {
		categoryName = content.Category.Name
	}
	var articleBlocks []dto.ArticleBlock
	if content.BodyBlocksJSON != "" {
		_ = json.Unmarshal([]byte(content.BodyBlocksJSON), &articleBlocks)
	}
	return dto.ContentResponse{
		ID:              content.ID,
		Title:           content.Title,
		Type:            content.Type,
		CategoryID:      content.CategoryID,
		CategoryName:    categoryName,
		FilePath:        content.FilePath,
		CoverURL:        content.CoverURL,
		Summary:         content.Summary,
		Status:          content.Status,
		VisibleRoles:    content.VisibleRoles,
		DurationSeconds: content.DurationSeconds,
		PublishAt:       content.PublishAt,
		ArticleBlocks:   articleBlocks,
	}
}
