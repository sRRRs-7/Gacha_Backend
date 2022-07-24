package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/sRRRs-7/GachaPon/db/sqlc"
	"github.com/sRRRs-7/GachaPon/token"
	"github.com/sRRRs-7/GachaPon/utils"
)

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	userRouter := router.Group("/user")
	userRouter.POST("/create", server.CreateUserApi)
	userRouter.POST("/login", server.LoginUserApi)
	userRouter.GET("/get/:user_name", server.GetUserApi)

	itemRouter := router.Group("/item")
	itemRouter.POST("/create", server.CreateItemApi)
	itemRouter.GET("/get/:id", server.GetItemApi)
	itemRouter.GET("/listByCategoryId", server.ListItemByCategoryIdApi)
	itemRouter.GET("/listByCategoriesId", server.ListItemsByCategoryIdApi)
	itemRouter.GET("/listById", server.ListItemsByIdApi)
	itemRouter.GET("/listByItemName", server.ListItemsByItemNameApi)
	itemRouter.GET("/listByRating", server.ListItemsByRatingApi)
	itemRouter.PUT("/update", server.UpdateItemApi)
	itemRouter.DELETE("/delete/:id", server.DeleteItemApi)

	categoryRouter := router.Group("/category")
	categoryRouter.POST("/create", server.CreateCategoryApi)
	categoryRouter.GET("/get/:category", server.GetCategoryApi)
	categoryRouter.GET("/list", server.ListCategoryApi)

	// authenticated router
	accountRouter := router.Group("/account").Use(authMiddleware(server.tokenMaker))
	accountRouter.POST("/create", server.CreateAccountApi)
	accountRouter.GET("/get/:id", server.GetAccountApi)
	accountRouter.GET("/list", server.ListAccountsApi)
	accountRouter.PUT("/updateAccount", server.UpdateAccountApi)
	accountRouter.PUT("/updateBalance", server.UpdateBalanceApi)
	accountRouter.DELETE("/delete/:id", server.DeleteAccountApi)

	galleryRouter := router.Group("/gallery").Use(authMiddleware(server.tokenMaker))
	galleryRouter.GET("/get/:id", server.GetGalleryApi)
	galleryRouter.GET("/listById", server.ListGalleriesByIdApi)
	galleryRouter.GET("/listByItemId", server.ListGalleriesByItemIdApi)

	gachaRouter := router.Group("/gacha").Use(authMiddleware(server.tokenMaker))
	gachaRouter.POST("/create", server.CreateGachaApi)
	gachaRouter.GET("/get/:id", server.GetGachaApi)
	gachaRouter.GET("/list", server.ListGachaApi)

	exchangeRouter := router.Group("/exchange").Use(authMiddleware(server.tokenMaker))
	exchangeRouter.POST("/create", server.CreateExchangeApi)
	exchangeRouter.GET("/get/:id", server.GetExchangeApi)
	exchangeRouter.GET("/listFromExchange", server.ListExchangeFromAccountApi)
	exchangeRouter.GET("/listToExchange", server.ListExchangeToAccountApi)

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errRes(err error) gin.H {
	return gin.H{"error": err.Error()}
}