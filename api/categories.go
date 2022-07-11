package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/sRRRs-7/GachaPon/db/sqlc"
)

type CreateCategoryRequest struct {
	Category string `json:"category" binding:"required"`
}

func (server *Server) CreateCategoryApi(ctx *gin.Context) {
	var req CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	category, err := server.store.CreateCategory(ctx, req.Category)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				ctx.JSON(http.StatusForbidden, errRes(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type GetCategoryRequest struct {
	Category string `uri:"category" binding:"required"`
}

func (server *Server) GetCategoryApi(ctx *gin.Context) {
	var req GetCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	category, err := server.store.GetCategory(ctx, req.Category)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
				return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type ListCategoryRequest struct {
	PageID 		int32 `form:"page_id" binding:"required,min=1"`
	PageSize 	int32 `form:"page_size" binding:"required,min=10"`
}

func (server *Server) ListCategoryApi(ctx *gin.Context) {
	var req ListCategoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.ListCategoriesParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	category, err := server.store.ListCategories(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
				return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}
