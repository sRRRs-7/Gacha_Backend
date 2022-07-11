package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/sRRRs-7/GachaPon/db/sqlc"
	"github.com/sRRRs-7/GachaPon/utils"
)

type CreateUserRequest struct {
	UserName     string `json:"user_name" binding:"required,alphanum"`
    HashPassword string `json:"hash_password" binding:"required,min=5"`
    FullName     string `json:"full_name" binding:"required"`
    Email        string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	UserName     string `json:"user_name"`
    FullName     string `json:"full_name"`
    Email        string `json:"email"`
}

func (server *Server) CreateUserApi(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	hashPassword, err := utils.HashedPassword(req.HashPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	arg := db.CreateUserParams{
		UserName: req.UserName,
		HashPassword: hashPassword,
		FullName: req.FullName,
		Email: req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
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

	res := CreateUserResponse{
		UserName: user.UserName,
		FullName: user.FullName,
		Email: user.Email,
	}

	ctx.JSON(http.StatusOK, res)
}

type GetUserRequest struct {
	UserName string `uri:"user_name" binding:"required"`
}

func (server *Server) GetUserApi(ctx *gin.Context) {
	var req GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type LoginUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required,min=5"`
}

type LoginUserResponse struct {
	UserName  			string    `json:"user_name"`
    UserAgent 			string    `json:"user_agent"`
    ClientIp  			string    `json:"client_ip"`
    IsBlocked 			bool      `json:"is_blocked"`
    ExpiredAt 			time.Time `json:"expired_at"`
	AccessToken         string    `json:"access_token"`
	AccessTokenExpired  time.Time `json:"access_token_expired_at"`
}

func (server *Server) LoginUserApi(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}

	err = utils.CheckPassword(req.Password, user.HashPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.UserName, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	arg := db.CreateSessionParams{
		UserName:     user.UserName,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiredAt:    accessPayload.ExpiredAt,
	}

	session, err := server.store.CreateSession(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}

	resp := LoginUserResponse{
		UserName: session.UserName,
		UserAgent: session.UserAgent,
		ClientIp: session.ClientIp,
		IsBlocked: session.IsBlocked,
		ExpiredAt: session.ExpiredAt,
		AccessToken: accessToken,
		AccessTokenExpired: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, resp)
}