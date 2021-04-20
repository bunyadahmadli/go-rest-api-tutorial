package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "statusCode: ${status}",
		}))
	e.GET("/home", mainHandler)

	adminGroup := e.Group("/admin", middleware.Logger())
	adminGroup.GET("/main", mainAdmin)

	e.GET("/user/:data", getUser)
	e.POST("/user", addUser)

	e.Start(":8099")
}
