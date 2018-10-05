package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r       *gin.Engine
	fetcher FetcherInteface
	storage StorageInterface
}

func NewServer(fetcher FetcherInteface, storage StorageInterface) *Server {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.MaxAge = 5 * time.Minute

	r.Use(cors.New(corsConfig))
	return &Server{
		r:       r,
		fetcher: fetcher,
		storage: storage,
	}
}

func (self *Server) GetListToken(c *gin.Context) {
	data, err := self.fetcher.GetListTokenAPI()
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"error": true, "reason": err.Error(), "data": data},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{"error": false, "data": data},
	)
}

func (self *Server) Run(port string) {
	self.r.GET("/currencies/getList", self.GetListToken)
	self.r.Run(fmt.Sprintf(":%s", port))
}
