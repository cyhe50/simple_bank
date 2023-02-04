package api

import (
	db "github.com/cyhe50/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

// server serves all HTTP request
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP request and setup routing
func NewServer(s *db.Store) (server *Server) {
	server = &Server{store: s}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	//add router later
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.createTransfer)
	server.router = router
	return
}

// run the HTTP servcer on a specific address to start listening for API request
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
