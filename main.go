package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func mainHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Main endpointine get isteği yapıldı")
}

func userHandler(c echo.Context) error {
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

func main() {
	e := echo.New()

	e.GET("/home", mainHandler)
	e.GET("/user/:data", userHandler)

	e.Start(":8099")
}
