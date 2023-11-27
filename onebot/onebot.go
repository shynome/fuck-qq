package onebot

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Inject(e *echo.Echo) {

	group := e.Group("/onebot")
	group.GET("/get_login_info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"user_id":  0,
			"nickname": "",
		})
	})
	initAPI(group)

}
