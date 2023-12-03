package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/khafizullokh02/simplebank/db/sqlc"
	"github.com/khafizullokh02/simplebank/util"
)

type Server struct {
	config util.Config
	store  db.Store
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	server := &Server{
		config: config,
		store:  store,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}
func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)

	authRoutes := router.Group("/")
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.PUT("/accounts/:id", server.updateAccount)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// start runs http server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
