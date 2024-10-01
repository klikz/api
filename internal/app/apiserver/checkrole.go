package apiserver

import (
	"warehouse/internal/app/models"

	"github.com/gin-gonic/gin"
)

func (s *Server) CheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := models.Request{}
		resp := models.Responce{}

		if err := c.ShouldBind(&req); err != nil {
			s.Logger.Error("Error Pasing body in CheckRole(): ", err)
			resp.Result = "error"
			resp.Error = err
			c.JSON(401, resp)
			c.Abort()
			return
		}
		parsedToken, err := ParseToken(req.Token)

		if err != nil {
			s.Logger.Error("Wrong Token: ", req.UserName, " error: ", err)
			resp.Result = "error"
			resp.Error = "Wrong Credentials"
			c.JSON(401, resp)
			c.Abort()
			return
		}

		user_id, _ := s.Store.Repo().GetUserID(parsedToken.UserName)

		res, err := s.Store.Repo().CheckRole(c.Request.URL.String(), user_id)
		if err != nil {
			s.Logger.Error("CheckRole: ", req.UserName, " error: ", err)
			resp.Result = "error"
			resp.Error = "Wrong Credentials"
			c.JSON(401, resp)
			c.Abort()
			return
		}

		if !res {
			s.Logger.Error("CheckRole: ", req.UserName, " error: ", err)
			resp.Result = "error"
			resp.Error = "Wrong Credentials"
			c.JSON(401, resp)
			c.Abort()
			return
		}

		c.Set("username", req.UserName)
		c.Set("password", req.Password)
		c.Set("name", req.Name)
		c.Set("model_id", req.ModelID)
		c.Set("product_id", req.ProductID)
		c.Set("file64", req.File64)
		c.Set("id", req.ID)
		c.Set("user_id", user_id)
		c.Set("date1", req.Date1)
		c.Set("date2", req.Date2)
		c.Set("c_time", req.C_Time)
		c.Set("status_id", req.Status_id)
		c.Set("serial", req.Serial)
		c.Set("printerid", req.PrinterID)
		c.Set("retry", req.Retry)

	}
}
