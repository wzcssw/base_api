package util

import (
	"strings"

	"github.com/labstack/echo"
)

// GetParam 获得请求参数
func GetParam(c echo.Context, paramName string) string {
	var result string
	req := c.Request()
	if req.Method == "GET" {
		result = c.QueryParam(paramName)
	} else {
		result = c.FormValue(paramName)
	}
	return result
}

// GetParamNames 获得请求参数
func GetParamNames(c echo.Context) map[string]string {
	var result map[string]string
	req := c.Request()
	if req.Method == "GET" {
		result = GetParamMapByQuery(c)
	} else {
		result = GetParamMapByForm(c)
	}
	return result
}

// GetParamMapByQuery 获得map
func GetParamMapByQuery(c echo.Context) map[string]string {
	params := map[string]string{}
	paramsNames := c.ParamNames()
	paramsValues := c.ParamValues()
	for i, p := range paramsNames {
		params[p] = paramsValues[i]
	}
	return params
}

// GetParamMapByForm 获得map
func GetParamMapByForm(c echo.Context) map[string]string {
	params := map[string]string{}
	fromParam, err := c.FormParams()
	if err != nil {
		return params
	}
	for key, value := range fromParam {
		if len(value) == 0 {
			params[key] = value[0]
		} else {
			params[key] = strings.Join(value, ",")
		}
	}
	return params
}
