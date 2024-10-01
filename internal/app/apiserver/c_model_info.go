package apiserver

import (
	"errors"
	"net/http"
	"os"
	"warehouse/internal/app/models"

	"github.com/gin-gonic/gin"
)

func (s *Server) CategoryGetAll(c *gin.Context) {
	resp := models.Responce{}

	data, err := s.Store.Repo().CategoryGetAll()
	if err != nil {
		s.SendError(c, err, "CategoryGetAll", "")
		return
	}
	resp.Result = "ok"
	resp.Data = data
	c.JSON(http.StatusOK, resp)
}

func (s *Server) CategoryAdd(c *gin.Context) {
	resp := models.Responce{}
	name := c.GetString("name")
	user_id := c.GetInt("user_id")
	if name == "" {
		s.SendError(c, errors.New("kategoriya nomi kiritilmadi"), "CategoryAdd", "")
		return
	}
	err := s.Store.Repo().CategoryAdd(name, user_id)
	if err != nil {
		s.SendError(c, err, "CategoryAdd", "")
		return
	}

	resp.Result = "ok"
	c.JSON(http.StatusOK, resp)
}

func (s *Server) CategoryDelete(c *gin.Context) {
	resp := models.Responce{}
	id := c.GetInt("id")
	user_id := c.GetInt("user_id")

	err := s.Store.Repo().CategoryDelete(id, user_id)

	if err != nil {
		s.SendError(c, err, "CategoryDelete", "")
		return
	}

	resp.Result = "ok"
	c.JSON(http.StatusOK, resp)
}

func (s *Server) ModelsGetAll(c *gin.Context) {
	resp := models.Responce{}

	id := c.GetInt("id")
	data, err := s.Store.Repo().ModelsGetAll(id)
	if err != nil {
		s.SendError(c, err, "ModelsGetAll", "")
		return
	}
	resp.Result = "ok"
	resp.Data = data
	c.JSON(http.StatusOK, resp)
}

func (s *Server) ModelsAdd(c *gin.Context) {
	resp := models.Responce{}
	id := c.GetInt("id")
	name := c.GetString("name")
	user_id := c.GetInt("user_id")
	if name == "" || id == 0 {
		s.SendError(c, errors.New("ma'lumotlar to'liq emas"), "ModelsAdd", "")
		return
	}
	err := s.Store.Repo().ModelAdd(name, user_id, id)
	if err != nil {
		s.Logger.Error("ModelsAdd: ", err.Error())
		resp.Result = "error"
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Result = "ok"
	c.JSON(http.StatusOK, resp)
}

func (s *Server) ModelsDelete(c *gin.Context) {
	resp := models.Responce{}
	id := c.GetInt("id")
	user_id := c.GetInt("user_id")

	err := s.Store.Repo().ModelsDelete(id, user_id)

	if err != nil {
		s.SendError(c, err, "ModelsDelete", "")
		return
	}

	resp.Result = "ok"
	c.JSON(http.StatusOK, resp)
}

func (s *Server) GsCodeFile(c *gin.Context) {
	resp := models.Responce{}
	id := c.GetInt("id")
	user_id := c.GetInt("user_id")
	file64 := c.GetString("file64")

	data, err := base64Decode(file64)
	if err != nil {
		s.SendError(c, err, "GsCodeFile, base64Decode", "")
		return
	}
	err = os.WriteFile("gscode.csv", []byte(data), 0644)
	if err != nil {
		s.SendError(c, err, "GsCodeFile, WriteFile", "")
		return
	}

	stringArray, err := readLines("gscode.csv")
	if err != nil {
		s.SendError(c, err, "GsCodeFile, readLines", "")
		return
	}

	badCode, err := s.Store.Repo().InsertGsCode(stringArray, id, user_id)
	if err != nil {
		type errStruct struct {
			ID   int    `json:"id"`
			Data string `json:"data"`
		}

		errorCode := []errStruct{}
		code := errStruct{}
		for i := 0; i < len(badCode); i++ {
			code.ID = i + 1
			code.Data = badCode[i]
			errorCode = append(errorCode, code)
		}
		s.SendError(c, err, "GsCodeFile, InsertGsCode", errorCode)
		return
	}

	resp.Result = "ok"
	c.JSON(http.StatusOK, resp)
}

func (s *Server) GsCodeGetList(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().GsCodeGetList()

	if err != nil {
		s.SendError(c, err, "GsCodeGetList", "")
		return
	}

	resp.Result = "ok"
	resp.Data = data
	c.JSON(http.StatusOK, resp)
}

func (s *Server) StatusGetAll(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().StatusGetAll()

	if err != nil {
		s.SendError(c, err, "StatusGetAll", "")
		return
	}

	resp.Result = "ok"
	resp.Data = data
	c.JSON(http.StatusOK, resp)
}
