package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/uditsaurabh/simple_bank/db/sqlc"
	token "github.com/uditsaurabh/simple_bank/token"
	util "github.com/uditsaurabh/simple_bank/util"
)

type Server struct {
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	config     *util.Config
}

func NewServer(store db.Store, config *util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create server %w", err)
	}
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}
	server.SetUpRoutes()
	return server, nil
}

func (server *Server) SetUpRoutes() {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/user", server.CreateUser)
	router.POST("/login", server.LoginWithPassword)

	//authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/account", server.listAccount)
	router.POST("/transfers", server.createTransfer)

	server.router = router
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	log.Println(err.Error())
	return gin.H{"error": err.Error()}
}
