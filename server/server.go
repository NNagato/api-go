package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	core "github.com/KyberNetwork/api-server/api-core"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r    *gin.Engine
	core *core.Core
}

func NewServer(core *core.Core) *Server {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.MaxAge = 5 * time.Minute

	r.Use(cors.New(corsConfig))
	return &Server{
		r:    r,
		core: core,
	}
}

func (self *Server) GetListToken(c *gin.Context) {
	data, err := self.core.GetListTokenAPI()
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

func (self *Server) GetAccountInfo(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{"error": false, "data": "comming soon"},
	)
}

func (self *Server) GetRateBuy(c *gin.Context) {
	id := c.QueryArray("id")
	qty := c.QueryArray("qty")
	if len(id) == 0 || len(id) != len(qty) {
		c.JSON(
			http.StatusOK,
			gin.H{"error": true, "reason": "params are invalid!"},
		)
		return
	}
	qtyFloat := [][]float64{}
	for _, qtyElem := range qty {
		qtyElemArr := strings.Split(qtyElem, "-")
		qFloatArr := []float64{}
		for _, q := range qtyElemArr {
			qFloat, err := strconv.ParseFloat(q, 64)
			if err != nil {
				c.JSON(
					http.StatusOK,
					gin.H{"error": true, "reason": "params are invalid!"},
				)
				return
			}
			qFloatArr = append(qFloatArr, qFloat)
		}
		qtyFloat = append(qtyFloat, qFloatArr)
	}
	isNewRate := self.core.GetIsNewRate()
	if isNewRate == false {
		c.JSON(
			http.StatusOK,
			gin.H{"error": true, "reason": "Can't get rate now"},
		)
		return
	}
	data, err := self.core.GetRateBuy(id, qtyFloat)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"error": true, "reason": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{"error": false, "data": data},
	)
}

func (self *Server) Run(port string) {
	currencies := self.r.Group("/currencies")
	currencies.GET("/getList", self.GetListToken)
	currencies.GET("/get_ethrate_buy", self.GetRateBuy)

	account := self.r.Group("/account")
	account.GET("/getInfo", self.GetAccountInfo)

	self.r.Run(fmt.Sprintf(":%s", port))
}
