package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/sRRRs-7/GachaPon/db/sqlc"
	"github.com/sRRRs-7/GachaPon/token"
)

type CreateExchangeRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required"`
	ToAccountID   int64 `json:"to_account_id" binding:"required"`
	ItemID1       int64 `json:"item_id_1" binding:"required"`
	ItemID2       int64 `json:"item_id_2" binding:"required"`
}

func (server *Server) CreateExchangeApi(ctx *gin.Context) {
	var req CreateExchangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	account ,err := server.store.GetAccount(ctx, req.FromAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}

	arg := db.ExchangeTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		ItemID1: req.ItemID1,
		ItemID2: req.ItemID2,
	}

	result, err := server.store.ExchangeTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type GetExchangeRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) GetExchangeApi(ctx *gin.Context) {
	var req GetExchangeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	exchange, err := server.store.GetExchange(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, exchange)
}

type ListExchangeFromAccountRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required"`
	PageID         int32 `json:"page_id" binding:"required,min=1"`
	PageSize       int32 `json:"page_size" binding:"required,max=10"`
}

func (server *Server) ListExchangeFromAccountApi(ctx *gin.Context) {
	var req ListExchangeFromAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.ListExchangeFromAccountParams{
		FromAccountID: req.FromAccountID,
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	exchanges, err := server.store.ListExchangeFromAccount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, exchanges)
}

type ListExchangeToAccountRequest struct {
	ToAccountID    int64 `json:"to_account_id" binding:"required"`
	PageID         int32 `json:"page_id" binding:"required,min=1"`
	PageSize       int32 `json:"page_size" binding:"required,max=10"`
}

func (server *Server) ListExchangeToAccountApi(ctx *gin.Context) {
	var req ListExchangeToAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.ListExchangeToAccountParams{
		ToAccountID: req.ToAccountID,
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	exchanges, err := server.store.ListExchangeToAccount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, exchanges)
}
