package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// Index 登录页面
func Index(c echo.Context) error {
	fmt.Printf("In Admin Index")
	return c.Render(http.StatusOK, "admin.html", map[string]interface{}{
		"test": "captchaID",
	})
}
