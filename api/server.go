package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/yanshen1997/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	res := &Server{
		store: store,
	}
	router := gin.Default()

	router.POST("/accounts", res.createAccount)
	router.GET("/accounts/:id", res.getAccount)
	router.GET("/accounts", res.listAccounts)
	router.DELETE("/accounts", res.deleteAccount)

	res.router = router
	return res
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}
