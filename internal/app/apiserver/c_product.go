package apiserver

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"warehouse/internal/app/models"

	"github.com/bingoohuang/xlsx"
	"github.com/gin-gonic/gin"
)

func (s *Server) LoadProductFile(c *gin.Context) {

	resp := models.Responce{}

	id := c.GetInt("id")
	model_id := c.GetInt("model_id")
	user_id := c.GetInt("user_id")
	file64 := c.GetString("file64")

	data, err := base64Decode(file64)
	if err != nil {
		s.SendError(c, err, "GsCodeCheckCount, base64Decode file", "")
		return
	}
	err = os.WriteFile("temp.xlsx", []byte(data), 0644)
	if err != nil {
		s.SendError(c, err, "GsCodeCheckCount, WriteFile file", "")
		return
	}

	var file []models.ProductInfo
	x, _ := xlsx.New(xlsx.WithInputFile("temp.xlsx"))
	defer x.Close()

	if err := x.Read(&file); err != nil {
		s.SendError(c, err, "GsCodeCheckCount, Read file", "")
		return
	}
	count := 0

	new_file := []models.ProductInfo{}
	for i := 0; i < len(file); i++ {
		t := strings.ReplaceAll(file[i].Serial, " ", "")
		if t == "" || t == " " {
		} else {
			new_file = append(new_file, file[i])
			count++
		}
	}

	err = s.Store.Repo().GsCodeCheckCount(model_id, count)
	if err != nil {
		errorString := fmt.Sprintf(`%s , faylgadi seriya soni:  %d`, err.Error(), count)
		s.SendError(c, errors.New(errorString), "LoadProductFile, GsCodeCheckCount", "")
		return
	}
	prodErr := []models.ProductInfo{}
	for i := 0; i < len(new_file); i++ {
		if err = s.Store.Repo().ProductAdd(id, model_id, user_id, &new_file[i]); err != nil {
			prodErr = append(prodErr, new_file[i])
		} else {
			if err = s.Store.Repo().GsCodeUpdate(new_file[i].ID, user_id, model_id); err != nil {
				s.SendError(c, err, "GsCodeCheckCount, GsCodeUpdate", "")
				return
			}
		}
	}

	if len(prodErr) > 0 {
		s.SendError(c, errors.New("seria nomer bazada yuklangan"), "GsCodeCheckCount, ProductAdd", prodErr)
		return
	}

	resp.Result = "ok"
	c.JSON(http.StatusOK, resp)
}

func (s *Server) ProductStatusListOutcome(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().ProductStatusListOutcome()

	if err != nil {
		s.SendError(c, err, "ProductStatusList", "")
		return
	}

	resp.Result = "ok"
	resp.Data = data
	c.JSON(http.StatusOK, resp)
}

func (s *Server) ProductStatusListIncome(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().ProductStatusListIncome()

	if err != nil {
		s.SendError(c, err, "ProductStatusListIncome", "")
		return
	}

	resp.Result = "ok"
	resp.Data = data
	c.JSON(http.StatusOK, resp)
}

func (s *Server) ProductStatusOutcome(c *gin.Context) {
	resp := models.Responce{}

	user_id := c.GetInt("user_id")
	file64 := c.GetString("file64")

	data, err := base64Decode(file64)
	if err != nil {
		s.SendError(c, err, "ProductStatusOutcome, base64Decode file", "")
		return
	}
	err = os.WriteFile("temp.xlsx", []byte(data), 0644)
	if err != nil {
		s.SendError(c, err, "ProductStatusOutcome, WriteFile file", "")
		return
	}

	var file []models.ProductInfo
	x, _ := xlsx.New(xlsx.WithInputFile("temp.xlsx"))
	defer x.Close()

	if err := x.Read(&file); err != nil {
		s.SendError(c, err, "ProductStatusOutcome, Read file", "")
		return
	}
	count := 0

	new_file := []models.ProductInfo{}
	for i := 0; i < len(file); i++ {
		t := strings.ReplaceAll(file[i].Serial, " ", "")
		if t == "" || t == " " {
		} else {
			new_file = append(new_file, file[i])
			count++
		}
	}

	serials, err := s.Store.Repo().ProductUpdateStatusToOutcome(user_id, new_file)
	if err != nil {
		s.SendError(c, err, "ProductStatusOutcome, ProductUpdateStatusToOutcome", serials)
		return
	}

	resp.Result = "ok"
	c.JSON(http.StatusOK, resp)
}

func (s *Server) ProductStatusIncome(c *gin.Context) {
	resp := models.Responce{}

	user_id := c.GetInt("user_id")
	file64 := c.GetString("file64")

	data, err := base64Decode(file64)
	if err != nil {
		s.SendError(c, err, "ProductStatusIncome, base64Decode file", "")
		return
	}
	err = os.WriteFile("temp.xlsx", []byte(data), 0644)
	if err != nil {
		s.SendError(c, err, "ProductStatusOutcome, WriteFile file", "")
		return
	}

	var file []models.ProductInfo
	x, _ := xlsx.New(xlsx.WithInputFile("temp.xlsx"))
	defer x.Close()

	if err := x.Read(&file); err != nil {
		s.SendError(c, err, "ProductStatusOutcome, Read file", "")
		return
	}
	count := 0

	new_file := []models.ProductInfo{}
	for i := 0; i < len(file); i++ {
		t := strings.ReplaceAll(file[i].Serial, " ", "")
		if t == "" || t == " " {
		} else {
			new_file = append(new_file, file[i])
			count++
		}
	}

	serials, err := s.Store.Repo().ProductUpdateStatusToIncome(user_id, new_file)
	if err != nil {
		s.SendError(c, err, "ProductStatusOutcome, ProductUpdateStatusToOutcome", serials)
		return
	}

	resp.Result = "ok"
	c.JSON(http.StatusOK, resp)
}

func (s *Server) ProductSerialInfo(c *gin.Context) {
	resp := models.Responce{}

	serial := c.GetString("serial")

	data, err := s.Store.Repo().ProductSerialInfo(serial)

	if err != nil {
		s.SendError(c, err, "ProductSerialInfo", nil)
		return
	}

	err = s.Store.Repo().ProductGeneratorBarcode(serial)
	if err != nil {
		s.SendError(c, err, "ProductSerialInfo, ProductGeneratorBarcode", nil)
		return
	}

	resp.Result = "ok"
	resp.Data = data
	c.JSON(200, resp)

}

func (s *Server) SkladSerialInput(c *gin.Context) {

	resp := models.Responce{}

	serial := c.GetString("serial")
	model_id := c.GetInt("model_id")
	category_id := c.GetInt("id")
	user_id := c.GetInt("user_id")
	printerId := c.GetInt("printerid")
	retry := c.GetBool("retry")

	product := models.ProductInfo{}
	product.Serial = serial

	err := s.Store.Repo().GetPrintInfo(&product)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			s.SendError(c, err, "SkladSerialInput, GetPrintInfo", "")
			return
		}
	}

	printerIp, printerName, err := s.Store.Repo().GetPrinterInfo(printerId)
	if err != nil {
		s.SendError(c, err, "SkladSerialInput, GetPrinterInfo", product.Serial)
		return
	}
	fmt.Println("retry: ", retry)
	fmt.Println("product.ID: ", product.ID)
	if product.ID == 0 {

		gscodeId, gscodeString, err := s.Store.Repo().GsCodeGetGetId(model_id)
		if err != nil {
			s.SendError(c, errors.New("GSCode yetarli emas"), "SkladSerialInput, GsCodeGetGetId", product.Serial)
			return
		}
		product.GsCode = gscodeString
		if err = s.Store.Repo().ProductAdd(category_id, model_id, user_id, &product); err != nil {
			s.SendError(c, err, "SkladSerialInput, ProductAdd", product.Serial)
			return
		}
		if err = s.Store.Repo().GsCodeUpdateByID(product.ID, user_id, gscodeId); err != nil {
			s.SendError(c, err, "SkladSerialInput, GsCodeUpdate", "")
			return
		}
	} else {

		if !retry {
			s.SendError(c, errors.New("serial kiritilgan"), "SkladSerialInput", product.Serial)
			return
		}
	}

	if err = SendGSCODE(product.GsCode, printerIp+":5555/gscode"); err != nil {
		s.SendError(c, errors.New("GS code yuborishda xatolik"), "SkladSerialInput, SendGSCODE", product.Serial)
		return
	}

	err = s.Store.Repo().GetPrintInfo(&product)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			s.SendError(c, err, "SkladSerialInput, GetPrintInfo", "")
			return
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	channel := make(chan string, 1)
	data := []byte(fmt.Sprintf(`
					{
						"libraryID": "986278f7-755f-4412-940f-a89e893947de",
						"absolutePath": "C:/inetpub/wwwroot/BarTender/wwwroot/Templates/premier/label.btw",
						"printRequestID": "fe80480e-1f94-4A2f-8947-e492800623aa",
						"printer": "%s",
						"startingPosition": 0,
						"copies": 0,
						"serialNumbers": 0,
						"dataEntryControls": {
								"modelInput": "%s",
								"serialInput": "%s",
								"categoryInput": "%s"
						}
					}`, printerName, product.Model, product.Serial, product.Category))

	// go PrintLabel(data, printerIp+"/BarTender/api/v1/print")
	go PrintLabel(data, channel, &wg, printerIp+":96/BarTender/api/v1/print")
	wg.Wait()
	errorText1 := <-channel

	if errorText1 != "ok" {
		fmt.Println("print error: ", errorText1)
		s.SendError(c, errors.New("printerda muammo. Qaytadan urinib ko'ring"), "SkladSerialInput, PrintLabel", serial)
		return
	}

	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) PrintersList(c *gin.Context) {
	resp := models.Responce{}

	data, err := s.Store.Repo().GetPrinterList()

	if err != nil {
		s.SendError(c, err, "PrintersList", nil)
		return
	}
	resp.Result = "ok"
	resp.Data = data
	c.JSON(200, resp)

}
func (s *Server) ProductGetLast(c *gin.Context) {
	resp := models.Responce{}
	id := c.GetInt("id")

	data, err := s.Store.Repo().ProductGetLast(id)

	if err != nil {
		s.SendError(c, err, "ProductGetLast", nil)
		return
	}
	resp.Result = "ok"
	resp.Data = data
	c.JSON(200, resp)

}
func (s *Server) ProductDelete(c *gin.Context) {
	resp := models.Responce{}
	id := c.GetInt("id")

	err := s.Store.Repo().ProductGSCodeClear(id)
	if err != nil {
		s.SendError(c, err, "ProductDelete, ProductGSCodeClear", nil)
		return
	}

	err = s.Store.Repo().ProductDelete(id)
	if err != nil {
		s.SendError(c, err, "ProductDelete", nil)
		return
	}
	resp.Result = "ok"
	// resp.Data = data
	c.JSON(200, resp)

}
