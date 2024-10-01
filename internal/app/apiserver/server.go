package apiserver

import (
	"io"
	"os"
	"warehouse/internal/app/store/sqlstore"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var Secret_key = []byte("Some123SecretKeyPremier1")

type Server struct {
	Router *gin.Engine
	Logger *log.Logger
	Store  sqlstore.Store
}

func newServer(store sqlstore.Store) *Server {
	s := &Server{
		Router: gin.New(),
		Logger: log.New(),
		Store:  store,
	}
	f, err := os.OpenFile("logger.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	wrt := io.MultiWriter(os.Stdout, f)
	// cors.AllowAll()
	// cors.Default()
	// s.Router.Use()

	s.Logger.SetOutput(wrt)
	s.Logger.SetFormatter(&log.JSONFormatter{})
	s.configureRouter()
	return s
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) configureRouter() {
	s.Router.SetTrustedProxies([]string{"localhost"})
	s.Router.Use(CORSMiddleware())

	s.Router.POST("user/login", s.Login)
	global := s.Router.Group("/api")
	global.Use(s.CheckRole())
	{
		global.POST("/user/create", s.Create)
		global.POST("/user/update", s.UpdatePassword)
		global.POST("/category/getall", s.CategoryGetAll)
		global.POST("/category/add", s.CategoryAdd)
		global.POST("/category/delete", s.CategoryDelete)
		global.POST("/models/getall", s.ModelsGetAll)
		global.POST("/models/add", s.ModelsAdd)
		global.POST("/models/delete", s.ModelsDelete)
		global.POST("/gscode/file", s.GsCodeFile)
		global.POST("/gscode/getlist", s.GsCodeGetList)
		global.POST("/load/file", s.LoadProductFile)
		global.POST("/products/getlist/outcome", s.ProductStatusListOutcome)
		global.POST("/products/status/outcome", s.ProductStatusOutcome)
		global.POST("/products/status/income", s.ProductStatusIncome)
		global.POST("/products/getlist/income", s.ProductStatusListIncome)
		global.POST("/products/serial/info", s.ProductSerialInfo)
		global.POST("/products/last", s.ProductGetLast)
		global.POST("/products/delete", s.ProductDelete)
		global.POST("/report/models", s.ReportGetByModels)
		global.POST("/report/serials", s.ReportGetBySerials)
		global.POST("/status/getall", s.StatusGetAll)
		global.POST("/sklad/serial/input", s.SkladSerialInput)
		global.POST("/printers/list", s.PrintersList)

	}

}
