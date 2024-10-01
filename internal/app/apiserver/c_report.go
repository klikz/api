package apiserver

import (
	"net/http"
	"warehouse/internal/app/models"

	"github.com/gin-gonic/gin"
)

func (s *Server) ReportGetByModels(c *gin.Context) {
	resp := models.Responce{}
	date1 := c.GetString("date1")
	date2 := c.GetString("date2")
	c_time := c.GetBool("c_time")
	model_id := c.GetInt("model_id")
	category_id := c.GetInt("id")
	status_id := c.GetInt("status_id")
	serial := c.GetString("serial")

	if c_time {
		data, err := s.Store.Repo().ReportGetByModelsC(date1, date2, serial, category_id, model_id, status_id)
		if err != nil {
			s.SendError(c, err, "ReportGetByModels, ReportGetByModelsC", "")
			return
		}
		resp.Data = data
	} else {
		data, err := s.Store.Repo().ReportGetByModelsU(date1, date2, serial, category_id, model_id, status_id)
		if err != nil {
			s.SendError(c, err, "ReportGetByModels, ReportGetByModelsU", "")
			return
		}
		resp.Data = data
	}

	resp.Result = "ok"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) ReportGetBySerials(c *gin.Context) {
	resp := models.Responce{}
	date1 := c.GetString("date1")
	date2 := c.GetString("date2")
	c_time := c.GetBool("c_time")
	model_id := c.GetInt("model_id")
	category_id := c.GetInt("id")
	status_id := c.GetInt("status_id")
	serial := c.GetString("serial")
	if c_time {
		data, err := s.Store.Repo().ReportGetBySerialsC(date1, date2, serial, category_id, model_id, status_id)
		if err != nil {
			s.SendError(c, err, "ReportGetBySerials, ReportGetBySerialsC", "")
			return
		}
		resp.Data = data
	} else {
		data, err := s.Store.Repo().ReportGetBySerialsU(date1, date2, serial, category_id, model_id, status_id)
		if err != nil {
			s.SendError(c, err, "ReportGetBySerials, ReportGetBySerialsU", "")
			return
		}
		resp.Data = data
	}

	resp.Result = "ok"

	c.JSON(http.StatusOK, resp)
}
