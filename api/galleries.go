package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/sRRRs-7/GachaPon/db/sqlc"
)

type GetGalleryRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) GetGalleryApi(ctx *gin.Context) {
	var req GetGalleryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	gallery, err := server.store.GetGallery(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusForbidden, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, gallery)
}

type ListGalleriesByIdRequest struct {
	OwnerID 	int64 `json:"owner_id" binding:"required,min=1"`
	PageID  	int32 `json:"page_id" binding:"required,min=1"`
    PageSize 	int32 `json:"page_size" binding:"required,min=10"`
}

func (server *Server) ListGalleriesByIdApi(ctx *gin.Context) {
	var req ListGalleriesByIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.ListGalleriesByIdParams{
		OwnerID: req.OwnerID,
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	gallery, err := server.store.ListGalleriesById(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, gallery)
}

type ListGalleriesByItemIdRequest struct {
    ItemID 		int64 `json:"item_id" binding:"required,min=1"`
	PageID  	int32 `json:"page_id" binding:"required,min=1"`
    PageSize 	int32 `json:"page_size" binding:"required,min=10"`
}

func (server *Server) ListGalleriesByItemIdApi(ctx *gin.Context) {
	var req ListGalleriesByItemIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.ListGalleriesByItemIdParams{
		ItemID: req.ItemID,
		Limit: 	req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	gallery, err := server.store.ListGalleriesByItemId(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, gallery)
}