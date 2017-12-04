// Package routers is good
// @APIVersion 1.0.0
// @Title mobile API
// @Description mobile has every tool to get any job done, so codename for the new mobile APIs.
// @Contact astaxie@gmail.com
package routers

import (
	"tx_base_api/api"

	"github.com/labstack/echo"
)

// InitRouter 加载路由配置
func InitRouter(e *echo.Echo) {
	e.POST("/logon", api.Logon)
	// user控制器
	e.GET("/user/get_all_user", api.GetAllUser)

}
