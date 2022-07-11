package api

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/sRRRs-7/GachaPon/db/sqlc"
)

type CreateGachaRequest struct {
	AccountID int64 `json:"account_id" binding:"required,min=1"`
}

func (server *Server) CreateGachaApi(ctx *gin.Context) {
	var req CreateGachaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusAccepted, errRes(err))
		return
	}

	wg := &sync.WaitGroup{}

	arg1 := db.ListItemsByIdParams{
		Limit: 1000,
		Offset: 0,
	}
	wg.Add(1)
	items, err := server.store.ListItemsById(ctx, arg1)
	if err != nil {
		fmt.Println("ListItemById error")
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}
	wg.Done()

	wg.Add(1)
	randomNum := rand.Int63n(int64(len(items) + 1))
	item, err := server.store.GetItem(ctx, randomNum)
	if err != nil {
		fmt.Println("GetItem error")
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}
	wg.Done()

	arg2 := db.CreateGachaParams{
		AccountID: req.AccountID,
		ItemID: item.ID,
	}
	wg.Add(1)
	gacha, err := server.store.CreateGacha(ctx, arg2)
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
	wg.Done()

	wg.Add(1)
	arg3 := db.CreateGalleryParams {
		OwnerID: gacha.AccountID,
		ItemID: gacha.ItemID,
	}
	gallery, err := server.store.CreateGallery(ctx, arg3)
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
	wg.Done()

	ctx.JSON(http.StatusOK, gallery)
}

type GetGachaRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) GetGachaApi(ctx *gin.Context) {
	var req GetGachaRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	gacha, err := server.store.GetGacha(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, gacha)
}

type ListGachaRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=10"`
}

func (server *Server) ListGachaApi(ctx *gin.Context) {
	var req ListGachaRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.ListGachasParams{
		Limit: req.PageSize,
		Offset:(req.PageID - 1) * req.PageSize,
	}

	gachas, err := server.store.ListGachas(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, gachas)
}
