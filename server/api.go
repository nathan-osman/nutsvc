package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) apiConfig_GET(c *gin.Context) {
	c.JSON(http.StatusOK, s.conf.GetAll())
}

func (s *Server) apiConfig_POST(c *gin.Context) {
	v := make(map[string]string)
	if err := c.ShouldBindJSON(&v); err != nil {
		panic(err)
	}
	if err := s.conf.SetMultiple(v); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (s *Server) apiMonitorReload_POST(c *gin.Context) {
	if err := s.tryReloadConf(); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{})
}
