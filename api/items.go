package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/sRRRs-7/GachaPon/db/sqlc"
)

type CreateItemRequest struct {
	ItemName   string `json:"item_name" binding:"required,alphanum"`
    Rating     int32  `json:"rating" binding:"required,min=1,max=7"`
    ItemUrl    string `json:"item_url" binding:"required"`
    CategoryID int32  `json:"category_id" binding:"required"`
}

func (server *Server) CreateItemApi(ctx *gin.Context) {
	var req CreateItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.CreateItemParams{
		ItemName: 	req.ItemName,
		Rating: 	req.Rating,
		ItemUrl: 	req.ItemUrl,
		CategoryID: req.CategoryID,
	}

	item, err := server.store.CreateItem(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name(){
				case "unique_violation":
					ctx.JSON(http.StatusForbidden, errRes(err))
					return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, item)
}

type GetItemRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) GetItemApi(ctx *gin.Context) {
	var req GetItemRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	item, err := server.store.GetItem(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, item)
}

type ListItemByCategoryIdRequest struct {
	CategoryID 	int32 `json:"category_id" binding:"required"`
    PageID      int32 `json:"page_id" binding:"required,min=1"`
    PageSize    int32 `json:"page_size" binding:"required,min=1"`
}

func (server *Server) ListItemByCategoryIdApi(ctx *gin.Context) {
	var req ListItemByCategoryIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.ListItemByCategoryIdParams{
		CategoryID: req.CategoryID,
		Limit: 		req.PageSize,
		Offset: 	(req.PageID - 1) * req.PageSize,
	}

	items, err := server.store.ListItemByCategoryId(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, items)
}

type ListItemsByCategoryIdRequest struct {
	PageID      int32 `json:"page_id" binding:"required,min=1"`
    PageSize    int32 `json:"page_size" binding:"required,min=1"`
}

func (server *Server) ListItemsByCategoryIdApi(ctx *gin.Context) {
	var req ListItemsByCategoryIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.ListItemsByCategoryIdParams{
		Limit: 	req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	items, err := server.store.ListItemsByCategoryId(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, items)
}

type ListItemsByIdRequest struct {
	PageID      int32 `json:"page_id" binding:"required,min=1"`
    PageSize    int32 `json:"page_size" binding:"required,min=1"`
}

func (server *Server) ListItemsByIdApi(ctx *gin.Context) {
	var req ListItemsByIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.ListItemsByIdParams{
		Limit: 	req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	items, err := server.store.ListItemsById(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, items)
}

type ListItemsByItemNameRequest struct {
	PageID      int32 `json:"page_id" binding:"required,min=1"`
    PageSize    int32 `json:"page_size" binding:"required,min=1"`
}

func (server *Server) ListItemsByItemNameApi(ctx *gin.Context) {
	var req ListItemsByItemNameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.ListItemsByItemNameParams{
		Limit: 	req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	items, err := server.store.ListItemsByItemName(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, items)
}

type ListItemsByRatingRequest struct {
	PageID      int32 `json:"page_id" binding:"required,min=1"`
    PageSize    int32 `json:"page_size" binding:"required,min=1"`
}

func (server *Server) ListItemsByRatingApi(ctx *gin.Context) {
	var req ListItemsByRatingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.ListItemsByRatingParams{
		Limit: 	req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	items, err := server.store.ListItemsByRating(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, items)
}

type UpdateItemRequest struct {
	ID         int64  `json:"id" binding:"required"`
    ItemName   string `json:"item_name" binding:"required,alphanum"`
    Rating     int32  `json:"rating" binding:"required,min=1,max=7"`
    ItemUrl    string `json:"item_url" binding:"required"`
    CategoryID int32  `json:"category_id" binding:"required"`
}

func (server *Server) UpdateItemApi(ctx *gin.Context) {
	var req UpdateItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.UpdateItemParams{
		ID: 		req.ID,
		ItemName: 	req.ItemName,
		Rating: 	req.Rating,
		ItemUrl: 	req.ItemUrl,
		CategoryID:	req.CategoryID,
	}

	item, err := server.store.UpdateItem(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, item)
}

type DeleteItemRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) DeleteItemApi(ctx *gin.Context) {
	var req DeleteItemRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusNotFound, errRes(err))
		return
	}

	err := server.store.DeleteItem(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}