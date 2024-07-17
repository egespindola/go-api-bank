package api

import (
	db "github.com/egespindola/go-api-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server servers HTTP requets for the banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	server.router = router

	return server
}
