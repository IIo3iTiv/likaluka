package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
	rgm    *gin.RouterGroup
}

func Init(mode string) *Server {
	var (
		router *gin.Engine
	)

	gin.SetMode(gin.DebugMode)
	if mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = os.Stdout
	router = gin.New()
	router.Use(gin.Recovery())
	return &Server{Router: router}
}

func (s *Server) Run(port string) {
	err := s.Router.Run(port)
	if err != nil {
		log.Fatalf("Server: Run. Error:%s", err.Error())
	}
}

func MWTimeout(td time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(td),
		timeout.WithResponse(func(c *gin.Context) { c.JSON(http.StatusRequestTimeout, []byte("timneout")) }),
	)
}
