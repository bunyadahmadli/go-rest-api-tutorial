package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

func mainHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Main endpointine get isteği yapıldı")
}

func getUser(c echo.Context) error {
	dataType := c.Param("data")
	username := c.QueryParam("username")
	name := c.QueryParam("name")
	surname := c.QueryParam("surname")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("usernme : %s, name : %s, surname: %s ", username, name, surname))
	}
	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name":     name,
			"username": username,
			"surname":  surname,
		})
	}

	return c.String(http.StatusBadRequest, "Parametre olarak Json ve ya String gire bilirsiniz")
}

func addUser(c echo.Context) error {
	user := User{}
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return err
	}
	fmt.Println(user)
	return c.String(http.StatusOK, "Başarılı")
}
func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "Admin endpointindesin")
}

func loginAdmin(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	if username == "admin" && password == "123" {
		cookie := &http.Cookie{
			Name:    "userId",
			Value:   "user_id",
			Expires: time.Now().Add(48 * time.Hour),
		}

		c.SetCookie(cookie)
		return c.String(http.StatusOK, "login olundu!")
	}

	return c.String(http.StatusUnauthorized, "Kullanıcı adı ve ya şifre hatalı")
}

func setHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		contentType := c.Request().Header.Get("Content-Type")
		if contentType != "application/json" {
			return c.JSON(http.StatusBadRequest, "Sadece Application/json tipinde istek atılabilir!")
		}

		return next(c)
	}
}
func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("userId")

		if err != nil {
			fmt.Println(err)
			if strings.Contains(err.Error(), "named cookie not present") {
				return c.JSON(http.StatusBadRequest, "Her hangi bir  cookie gönderilmedi!")
			}
			return err
		}

		if cookie.Value == "user_id" {
			return next(c)
		}

		return c.JSON(http.StatusBadRequest, "Doğru cookie gönderilmedi!")
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "statusCode: ${status}",
		}))
	e.GET("/home", mainHandler)
	//e.Use(setHeader)
	adminGroup := e.Group("/admin")

	adminGroup.GET("/main", mainAdmin, checkCookie)
	adminGroup.GET("/login", loginAdmin)
	// adminGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	if username == "admin" && password == "123" {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))
	e.GET("/user/:data", getUser)
	e.POST("/user", addUser)

	e.Start(":8099")
}
