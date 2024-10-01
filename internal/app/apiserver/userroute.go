package apiserver

import (
	"net/http"
	"warehouse/internal/app/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (s *Server) Create(c *gin.Context) {

	user := models.User{}
	resp := models.Responce{}
	user.UserName = c.GetString("username")
	user.Password = c.GetString("password")

	enc, err := encryptString(user.Password)
	if err != nil {
		s.Logger.Error("Create: encryptString: ", err.Error())
		resp.Result = "error"
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	user.EncryptedPassword = enc

	err = s.Store.Repo().Create(&user)
	if err != nil {
		s.Logger.Error("Create: Create: ", err.Error())
		resp.Result = "error"
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Result = "ok"
	c.JSON(200, resp)

}

func (s *Server) UpdatePassword(c *gin.Context) {

	user := models.User{}
	resp := models.Responce{}
	user.UserName = c.GetString("username")
	user.Password = c.GetString("password")

	enc, err := encryptString(user.Password)
	if err != nil {
		s.Logger.Error("UpdatePassword: encryptString: ", err.Error())
		resp.Result = "error"
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	user.EncryptedPassword = enc
	err = s.Store.Repo().UpdatePassword(&user)
	if err != nil {
		s.Logger.Error("UpdatePassword: UpdatePassword: ", err.Error())
		resp.Result = "error"
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Result = "ok"
	c.JSON(200, resp)

}

func (s *Server) Login(c *gin.Context) {
	user := models.User{}
	resp := models.Responce{}

	if err := c.ShouldBind(&user); err != nil {
		logrus.Error("Login: Error Parsing body: ", err.Error())
	}

	if err := s.Store.Repo().FindByUserName(&user); err != nil {
		resp.Result = "error"
		resp.Error = "Wrong Credentials"
		s.Logger.Error("Login: FindByUserName: ", user.UserName, " error: ", err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if !ComparePassword(user.Password, user.EncryptedPassword) {
		resp.Result = "error"
		resp.Error = "Wrong Credentials"
		s.Logger.Error("Login: ComparePassword: ", user.UserName)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := GetToken(&user); err != nil {
		s.Logger.Error("Login: GetToken: ", user.UserName, " error: ", err.Error())
		resp.Result = "error"
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	user.Password = ""
	s.Logger.Info("User Logged: ", user.UserName, " client ip: ", c.ClientIP(), " remote ip: ", c.RemoteIP())
	c.JSON(200, user)

}
