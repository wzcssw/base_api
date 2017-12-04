package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/labstack/echo"
)

// BodyToJSON 获取JSON参数
func BodyToJSON(c echo.Context) (map[string]interface{}, error) {
	s, err := ioutil.ReadAll(c.Request().Body)
	fmt.Printf("============>body_str is %s\n", s)
	if err != nil {
		return nil, err
	}
	var body map[string]interface{}
	if err := json.Unmarshal(s, &body); err != nil {
		return nil, err
	}
	return body, nil
}
