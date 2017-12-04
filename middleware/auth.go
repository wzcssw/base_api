package mw

import (
	"fmt"
	"strings"

	session "github.com/ipfans/echo-session"
	"github.com/labstack/echo"
)

type (
	// Auth 认证中间件
	Auth struct {
		IsLogin    bool     `json:"is_login"`
		PassURLArr []string `json:"pass_url_arr"`
	}
)

// NewAuth 初始化认证中间件
func NewAuth() *Auth {
	return &Auth{
		IsLogin:    false,
		PassURLArr: []string{"/static", "/login", "/logon", "/logout", "/captcha"},
	}
}

// CanPass 判断是否已能通过访问
func (a *Auth) CanPass(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		isPass := false
		// 先判断访问的是否pass_url
		curURL := c.Path()
		fmt.Printf("===============>current path is %s\n", curURL)
		for _, url := range a.PassURLArr {
			if strings.Index(curURL, url) == 0 {
				isPass = true
				break
			}
		}
		if isPass == false {
			sess := session.Default(c)
			userID := sess.Get("user_id")
			fmt.Print("===============>in auth user_id is ", userID)
			if userID == nil {
				a.IsLogin = false
				return c.Redirect(302, "/login")
			}
		}
		return next(c)
	}
}

// IsLogin 判断是否登录
// func (s *Auth) IsLogin(c echo.Context) error {
// 	s.mutex.RLock()
// 	defer s.mutex.RUnlock()
// 	return c.JSON(http.StatusOK, s)
// }
