package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/sRRRs-7/GachaPon/db/sqlc"
	"github.com/sRRRs-7/GachaPon/token"
)

type CreateAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
	Balance float64 `json:"balance" binding:"required"`
}

func (server *Server) CreateAccountApi(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner: authPayload.Username,
		Balance: int64(req.Balance),
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errRes(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) GetAccountApi(ctx *gin.Context) {
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountsRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=10"`
}

func (server *Server) ListAccountsApi(ctx *gin.Context) {
	var req ListAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner: authPayload.Username,
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type UpdateAccountRequest struct {
	ID 		int64 `json:"id" binding:"required"`
	Balance int64 `json:"balance" binding:"required"`
}

func (server *Server) UpdateAccountApi(ctx *gin.Context) {
	var req UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID: req.ID,
		Balance: req.Balance,
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type UpdateBalanceRequest struct {
	ID int64 `json:"id" binding:"required"`
	Balance int64 `json:"balance" binding:"required"`
}

func (server *Server) UpdateBalanceApi(ctx *gin.Context) {
	var req UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID: req.ID,
		Balance: req.Balance,
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type DeleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) DeleteAccountApi(ctx *gin.Context) {
	var req UpdateAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}