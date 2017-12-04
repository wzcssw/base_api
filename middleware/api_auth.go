package mw

import (
	"fmt"
	"strings"
	"tx_base_api/models"
	"tx_base_api/util"

	"github.com/labstack/echo"
)

type (
	// APIAuth 认证中间件
	APIAuth struct {
		IsLogin    bool     `json:"is_login"`
		PassURLArr []string `json:"pass_url_arr"`
	}
)

// NewAPIAuth 初始化认证中间件
func NewAPIAuth() *APIAuth {
	return &APIAuth{
		IsLogin:    false,
		PassURLArr: []string{"/static"},
	}
}

// CanPass 判断是否已能通过访问
func (a *APIAuth) CanPass(next echo.HandlerFunc) echo.HandlerFunc {
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
			flag, result := a.BaseValidate(c)
			if !flag {
				return c.JSON(500, result)
			}
			flag, result = a.UserValidate(c)
			if !flag {
				return c.JSON(500, result)
			}
		}
		return next(c)
	}
}

// BaseValidate 基础验证
func (a *APIAuth) BaseValidate(c echo.Context) (flag bool, result map[string]interface{}) {
	fmt.Printf("===============>in BaseValidate \n")
	fmt.Printf("===============>params is %v\n", util.GetParam(c, "app_key"))
	result = map[string]interface{}{}
	timeStamp := util.GetParam(c, "timestamp")
	appKey := util.GetParam(c, "app_key")
	if timeStamp == "" || appKey == "" {
		flag = false
		result["error"] = "时间戳或者appkey不存在"
		return
	}
	fmt.Printf("===============>timeStamp is %s,and appKey is %s\n", timeStamp, appKey)
	sign := util.GetParam(c, "sign")
	params := util.GetParamNames(c)
	var appToken models.AppToken
	models.Mdb.First(&appToken, "app_key = ?", appKey)
	if appToken.ID == 0 {
		flag = false
		result["error"] = strings.Join([]string{"请求的appkey:", appKey, "不存在"}, "")
		return
	}
	params["app_token"] = appToken.Token
	tmpParams := map[string]string{}
	for k, v := range params {
		tmpParams[k] = v
	}
	delete(tmpParams, "sign")
	signWord := util.Sign(tmpParams)
	if sign == signWord {
		flag = true
	} else {
		flag = false
		result["error"] = strings.Join([]string{"签名不一致,raw串为:", sign, "gen串为:", signWord}, "")
	}
	return
}

// UserValidate 登录验证
func (a *APIAuth) UserValidate(c echo.Context) (flag bool, result map[string]interface{}) {
	result = map[string]interface{}{}
	accessToken := util.GetParam(c, "access_token")
	fmt.Printf("===============>access_token is %s\n", accessToken)
	if accessToken == "" {
		flag = false
		result["error"] = "用户未认证"
	}
	var apikey models.APIKey
	models.Mdb.First(&apikey, "login_type = 'user' and access_token = ?", accessToken)
	if apikey.ID == 0 {
		flag = false
		result["error"] = "用户未认证"
	} else if apikey.IsExpired() {
		flag = false
		result["error"] = "认证信息已过期"
	} else {
		flag = true
	}
	return
}

// // CanPass 判断是否已能通过访问
// func (a *APIAuth) CanPass(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		isPass := false
// 		// 先判断访问的是否pass_url
// 		curURL := c.Path()
// 		fmt.Printf("===============>current path is %s\n", curURL)
// 		for _, url := range a.PassURLArr {
// 			if strings.Index(curURL, url) == 0 {
// 				isPass = true
// 				break
// 			}
// 		}
// 		if isPass == false {
// 			flag, result := a.TransformParam(c)
// 			fmt.Printf("#############Flag is %v\n", flag)
// 			if !flag {
// 				fmt.Printf("Error is %s\n", result)
// 				return c.JSON(500, result)
// 			}
// 			flag, result = a.BaseValidate(c)
// 			if !flag {
// 				return c.JSON(500, result)
// 			}
// 			flag, result = a.UserValidate(c)
// 			if !flag {
// 				return c.JSON(500, result)
// 			}
// 		}
// 		return next(c)
// 	}
// }

// // TransformParam 转换参数
// func (a *APIAuth) TransformParam(c echo.Context) (flag bool, result map[string]interface{}) {
// 	fmt.Printf("===============>in TransformParam \n")
// 	jsonBody, err := util.BodyToJSON(c)
// 	result = map[string]interface{}{}
// 	if err != nil {
// 		fmt.Printf("=============>转换请求数据为JSON格式时出错,%s\n", err)
// 		flag = false
// 		result["error"] = "请求body转换请求数据为JSON格式时出错" + err.Error()
// 		return
// 	}
// 	fmt.Printf("=============>JSONBody is ,%v\n", jsonBody)
// 	c.Set("JSONBody", jsonBody)
// 	return
// }

// // BaseValidate 基础验证
// func (a *APIAuth) BaseValidate(c echo.Context) (flag bool, result map[string]interface{}) {
// 	fmt.Printf("===============>in BaseValidate \n")
// 	jsonBody := c.Get("JSONBody").(map[string]interface{})
// 	fmt.Printf("===============>params is %v\n", jsonBody["app_key"])
// 	timeStamp := jsonBody["timestamp"].(string)
// 	appKey := jsonBody["app_key"].(string)
// 	if timeStamp == "" || appKey == "" {
// 		flag = false
// 		result["error"] = "时间戳或者appkey不存在"
// 		return
// 	}
// 	fmt.Printf("===============>timeStamp is %s,and appKey is %s\n", timeStamp, appKey)
// 	sign := jsonBody["sign"].(string)
// 	var appToken models.AppToken
// 	models.Mdb.First(&appToken, "app_key = ?", appKey)
// 	if appToken.ID == 0 {
// 		flag = false
// 		result["error"] = strings.Join([]string{"请求的appkey:", appKey, "不能存在"}, "")
// 		return
// 	}
// 	jsonBody["app_token"] = appToken.Token
// 	delete(jsonBody, "sign")
// 	signWord := util.SignJSON(jsonBody)
// 	fmt.Printf("===============>Sign is %s\n", sign)
// 	fmt.Printf("===============>SignWord is %s\n", signWord)
// 	if sign == signWord {
// 		flag = true
// 	} else {
// 		flag = false
// 		result["error"] = strings.Join([]string{"签名不一致,raw串为:", sign, "gen串为:", signWord}, "")
// 	}
// 	fmt.Printf("===============>Finish！！！！！！！！\n")
// 	return
// }

// // UserValidate 登录验证
// func (a *APIAuth) UserValidate(c echo.Context) (flag bool, result map[string]interface{}) {
// 	jsonBody := c.Get("JSONBody").(map[string]interface{})
// 	accessToken := jsonBody["access_token"]
// 	fmt.Printf("===============>access_token is %s\n", accessToken)
// 	if accessToken == "" {
// 		flag = false
// 		result["error"] = "用户未认证"
// 	}
// 	var apikey models.APIKey
// 	models.Mdb.First(&apikey, "login_type = 'user' and access_token = ?", accessToken)
// 	if apikey.ID == 0 {
// 		flag = false
// 		result["error"] = "用户未认证"
// 	} else if apikey.IsExpired() {
// 		flag = false
// 		result["error"] = "认证信息已过期"
// 	} else {
// 		flag = true
// 	}
// 	return
// }
