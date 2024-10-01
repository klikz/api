package apiserver

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"warehouse/internal/app/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// func isPasswordValid(p string) bool {
// 	return len(p) >= 6
// }

func ComparePassword(password, encrypt string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encrypt), []byte(password)) == nil
}

func GetToken(u *models.User) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": u.UserName,
		"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	tokenString, err := token.SignedString(Secret_key)
	if err != nil {
		return err
	}
	u.Token = tokenString
	return nil
}

func ParseToken(tokenString string) (models.ParsedToken, error) {
	parsedToken := models.ParsedToken{}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return Secret_key, nil
	})
	if err != nil {
		return parsedToken, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

	} else {
		return parsedToken, errors.New("wrong token")
	}
	parsedToken.UserName = fmt.Sprint(claims["username"])

	return parsedToken, nil
}

func base64Decode(str string) (string, error) {

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (s *Server) SendError(c *gin.Context, err error, route string, data interface{}) {
	resp := models.Responce{}
	s.Logger.Error(route, ": ", err.Error())
	resp.Result = "error"
	resp.Error = err.Error()
	resp.Data = data
	c.JSON(http.StatusBadRequest, resp)
}

// func PrintLabel(jsonStr []byte, printerUrl string) {
// 	reprint := true
// 	count := 0

// 	for reprint {
// 		if count > 3 {
// 			return
// 		}
// 		req, err := http.NewRequest("POST", printerUrl, bytes.NewBuffer(jsonStr))
// 		if err != nil {
// 			fmt.Println("NewRequest: ", err.Error())
// 		}
// 		req.Header.Set("X-Custom-Header", "myvalue")
// 		req.Header.Set("Content-Type", "application/json")
// 		client := &http.Client{}
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			fmt.Println("resp: ", err.Error())
// 		}
// 		defer resp.Body.Close()

// 		body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Println("ReadAll: ", err.Error())
// 		}
// 		var jsonMap map[string]interface{}
// 		json.Unmarshal([]byte(string(body)), &jsonMap)

// 		if strings.Contains(string(body), "BarTender успешно отправил задание") {
// 			reprint = false
// 			return
// 		}
// 		count++
// 	}

// }

func PrintLabel(jsonStr []byte, channel chan string, wg *sync.WaitGroup, printerUrl string) {

	defer wg.Done()
	defer close(channel)
	reprint := true
	count := 0

	fmt.Println("print name full: ", printerUrl)

	for reprint {
		fmt.Println("try to print: ", count)
		if count > 3 {
			channel <- "qaytadan urinib ko'ring"
			return
		}
		req, err := http.NewRequest("POST", printerUrl, bytes.NewBuffer(jsonStr))
		if err != nil {
			channel <- err.Error()
			return
		}
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			channel <- err.Error()
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var jsonMap map[string]interface{}
		json.Unmarshal([]byte(string(body)), &jsonMap)
		fmt.Println("body from: ", string(body))
		if strings.Contains(string(body), "BarTender успешно отправил задание") {
			reprint = false
			channel <- "ok"
			return
		}
		count++
	}

	channel <- "error"
	// channel <- "ok"

}

func SendGSCODE(gscode string, gscodeUrl string) error {
	fmt.Println(gscodeUrl)
	response, err := http.PostForm(gscodeUrl, url.Values{
		"gscode": {gscode}})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Println("body: ", string(body))
	return nil
}
